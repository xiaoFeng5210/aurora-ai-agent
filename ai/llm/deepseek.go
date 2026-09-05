package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"aurora-agent/ai"
)

const DeepSeekModel = "deepseek-v4-flash"

// DeepSeek owns transport, request conversion and response normalization.
// It is stateless between calls; reasoning travels with the conversation.
type DeepSeek struct {
	apiKey   string
	endpoint string
	client   *http.Client
}

func NewDeepSeek() *DeepSeek {
	return &DeepSeek{
		apiKey:   strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY")),
		endpoint: "https://api.deepseek.com/chat/completions",
		client:   &http.Client{Timeout: 10 * time.Minute},
	}
}

func (d *DeepSeek) Name() string { return DeepSeekModel }

// Provider-only fields never leak into the application's SSE/history JSON.
type deepSeekMessage struct {
	Role       string             `json:"role"`
	Content    string             `json:"content"`
	Reasoning  *string            `json:"reasoning_content,omitempty"`
	ToolCallID *string            `json:"tool_call_id,omitempty"`
	ToolCalls  []deepSeekToolCall `json:"tool_calls,omitempty"`
}

type deepSeekToolCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Function ai.FunctionCall `json:"function"`
}

func (d *DeepSeek) Chat(ctx context.Context, messages []ai.Message, opts ChatOptions, onEvent StreamEventHandler) (ChatResult, error) {
	if d.apiKey == "" {
		return ChatResult{}, fmt.Errorf("DEEPSEEK_API_KEY is not set")
	}
	if opts.MaxTokens == 0 {
		opts.MaxTokens = 65536
	}
	if opts.ThinkingType == "" {
		opts.ThinkingType = "disabled"
	}
	if opts.MaxTokens < 1 || opts.MaxTokens > 393216 {
		return ChatResult{}, fmt.Errorf("max_tokens must be between 1 and 393216")
	}
	if opts.ThinkingType != "enabled" && opts.ThinkingType != "disabled" {
		return ChatResult{}, fmt.Errorf("thinking.type must be enabled or disabled")
	}
	if opts.Temperature < 0 || opts.Temperature > 2 || opts.TopP < 0 || opts.TopP > 1 {
		return ChatResult{}, fmt.Errorf("invalid sampling options")
	}
	wireMessages := make([]deepSeekMessage, 0, len(messages))
	for _, m := range messages {
		wire := deepSeekMessage{Role: m.Role, Content: m.Content, ToolCallID: m.ToolCallId}
		if m.Role == "assistant" {
			if opts.ThinkingType == "enabled" {
				reasoning := m.ReasoningContent
				wire.Reasoning = &reasoning
			}
			for _, tc := range m.ToolCalls {
				wire.ToolCalls = append(wire.ToolCalls, deepSeekToolCall{ID: tc.Id, Type: tc.Type, Function: tc.Function})
			}
		}
		wireMessages = append(wireMessages, wire)
	}
	body := map[string]any{"model": d.Name(), "messages": wireMessages, "max_tokens": opts.MaxTokens, "stream": true, "thinking": map[string]string{"type": opts.ThinkingType}}
	if opts.Tools != nil {
		body["tools"] = opts.Tools
	}
	// Sampling parameters have no effect in thinking mode.
	if opts.ThinkingType == "disabled" {
		if opts.TopP > 0 {
			body["top_p"] = opts.TopP
		} else if opts.Temperature > 0 {
			body["temperature"] = opts.Temperature
		}
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return ChatResult{}, fmt.Errorf("encode model request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.endpoint, bytes.NewReader(payload))
	if err != nil {
		return ChatResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+d.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	resp, err := d.client.Do(req)
	if err != nil {
		return ChatResult{}, fmt.Errorf("deepseek request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Do not relay upstream bodies: they can contain request data or credentials.
		return ChatResult{}, fmt.Errorf("deepseek request failed (HTTP %d)", resp.StatusCode)
	}
	return readDeepSeekStream(resp.Body, onEvent)
}

type deepSeekChunk struct {
	Error   json.RawMessage `json:"error"`
	Choices []struct {
		Delta struct {
			Content   string        `json:"content"`
			Reasoning string        `json:"reasoning_content"`
			ToolCalls []ai.ToolCall `json:"tool_calls"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func readDeepSeekStream(body io.Reader, onEvent StreamEventHandler) (ChatResult, error) {
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	result := ChatResult{Message: ai.Message{Role: "assistant"}}
	var content, reasoning strings.Builder
	calls := map[int]ai.ToolCall{}
	done := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		} // Includes keep-alive comments.
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" {
			continue
		}
		if data == "[DONE]" {
			done = true
			break
		}
		var chunk deepSeekChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return ChatResult{}, fmt.Errorf("invalid deepseek stream JSON")
		}
		if len(chunk.Error) > 0 && string(chunk.Error) != "null" {
			return ChatResult{}, fmt.Errorf("deepseek stream returned an error")
		}
		if len(chunk.Choices) == 0 {
			continue
		} // Tolerate usage-only chunks.
		choice := chunk.Choices[0]
		if choice.Delta.Content != "" {
			content.WriteString(choice.Delta.Content)
			emitStreamEvent(onEvent, "delta", map[string]any{"content": choice.Delta.Content})
		}
		reasoning.WriteString(choice.Delta.Reasoning)
		for _, delta := range choice.Delta.ToolCalls {
			current := calls[delta.Index]
			current.Index = delta.Index
			if delta.Id != "" {
				current.Id = delta.Id
			}
			if delta.Type != "" {
				current.Type = delta.Type
			}
			current.Function.Name += delta.Function.Name
			// DeepSeek sends incremental fragments, never cumulative snapshots.
			current.Function.Arguments += delta.Function.Arguments
			calls[delta.Index] = current
		}
		if choice.FinishReason != "" {
			result.FinishReason = choice.FinishReason
		}
	}
	if err := scanner.Err(); err != nil {
		return ChatResult{}, fmt.Errorf("read deepseek stream: %w", err)
	}
	if !done || result.FinishReason == "" {
		return ChatResult{}, fmt.Errorf("deepseek stream ended prematurely")
	}
	if result.FinishReason != "stop" && result.FinishReason != "tool_calls" {
		return ChatResult{}, fmt.Errorf("deepseek generation ended with %s", result.FinishReason)
	}
	result.Message.Content = content.String()
	result.Message.ReasoningContent = reasoning.String()
	indices := make([]int, 0, len(calls))
	for index := range calls {
		indices = append(indices, index)
	}
	sort.Ints(indices)
	for _, index := range indices {
		call := calls[index]
		if call.Id == "" || call.Type != "function" || call.Function.Name == "" || !json.Valid([]byte(call.Function.Arguments)) {
			return ChatResult{}, fmt.Errorf("deepseek returned an incomplete tool call")
		}
		result.Message.ToolCalls = append(result.Message.ToolCalls, call)
	}
	if (result.FinishReason == "tool_calls") != (len(calls) > 0) {
		return ChatResult{}, fmt.Errorf("deepseek finish reason does not match tool calls")
	}
	if len(calls) > 0 {
		emitStreamEvent(onEvent, "tool_call", map[string]any{"tool_calls": result.Message.ToolCalls})
	}
	return result, nil
}

func emitStreamEvent(handler StreamEventHandler, event string, data any) {
	if handler != nil {
		handler(event, data)
	}
}
