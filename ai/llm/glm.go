package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"aurora-agent/ai"
	functioncall "aurora-agent/ai/llm/function-call"
	utils "aurora-agent/utils"
)

const (
	GLM_MODEL_BASE_URL = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
)

type Delta struct {
	Role      string         `json:"role"`
	Content   string         `json:"content"`
	ToolCalls *[]ai.ToolCall `json:"tool_calls"`
}

type Choice struct {
	Delta Delta `json:"delta"`
}

type StreamResponse struct {
	Id      string   `json:"id"`
	Choices []Choice `json:"choices"`
}

type ChatOptions struct {
	MaxTokens    int
	Temperature  float64
	TopP         float64
	ThinkingType string
}

type StreamEventHandler func(event string, data any)

type GLM struct {
	APIKey   string
	Model    string
	MaxToken int
	Stream   bool

	Temperature  float64
	TopP         float64
	ThinkingType string
}

var (
	logger = utils.Logger
)

func InitModel() *GLM {
	apiKey := os.Getenv("GLM_API_KEY")
	glm := &GLM{
		APIKey:       apiKey,
		Model:        "glm-4.7",
		MaxToken:     65536,
		Stream:       true,
		Temperature:  1.0,
		TopP:         0.0,
		ThinkingType: "disabled",
	}
	return glm
}

func (glm *GLM) SetOptions(opts ChatOptions) {
	if opts.MaxTokens > 0 {
		glm.MaxToken = opts.MaxTokens
	}
	if opts.Temperature > 0 {
		glm.Temperature = opts.Temperature
	}
	if opts.TopP > 0 {
		glm.TopP = opts.TopP
	}
	if opts.ThinkingType != "" {
		glm.ThinkingType = opts.ThinkingType
	}
}

func (glm *GLM) ChatWithGLMInStream(messages []ai.Message) (bool, []ai.ToolCall, string, error) {
	return glm.ChatWithGLMInStreamWithEvents(messages, ChatOptions{}, nil)
}

func (glm *GLM) ChatWithGLMInStreamWithEvents(messages []ai.Message, opts ChatOptions, onEvent StreamEventHandler) (bool, []ai.ToolCall, string, error) {
	glm.SetOptions(opts)
	if glm.APIKey == "" {
		return false, nil, "", fmt.Errorf("GLM_API_KEY is not set")
	}

	requestBody := map[string]interface{}{
		"model":      glm.Model,
		"messages":   messages,
		"max_tokens": glm.MaxToken,
		"stream":     true,
		"thinking": map[string]interface{}{
			"type": glm.ThinkingType,
		},
		"tools": functioncall.WeatherTools,
	}

	if glm.Temperature > 0 {
		requestBody["temperature"] = glm.Temperature
	}
	if glm.TopP > 0 {
		requestBody["top_p"] = glm.TopP
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return false, nil, "", err
	}

	request, err := http.NewRequest("POST", GLM_MODEL_BASE_URL, bytes.NewBuffer(body))
	if err != nil {
		return false, nil, "", err
	}
	request.Header.Set("Authorization", "Bearer "+glm.APIKey)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return false, nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return false, nil, "", fmt.Errorf("glm request failed with status %d", resp.StatusCode)
		}
		return false, nil, "", fmt.Errorf("glm request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return glm.AIStreamResponseHandlerWithEvents(resp.Body, onEvent)
}

func (glm *GLM) AIStreamResponseHandler(body io.Reader) (bool, []ai.ToolCall, string, error) {
	return glm.AIStreamResponseHandlerWithEvents(body, nil)
}

// AI 流式回答处理
func (glm *GLM) AIStreamResponseHandlerWithEvents(body io.Reader, onEvent StreamEventHandler) (bool, []ai.ToolCall, string, error) {
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	content := ""
	toolCallMap := map[int]ai.ToolCall{}
	toolCallOrder := []int{}

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 || !bytes.HasPrefix(line, []byte("data:")) {
			continue
		}

		segment := bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data:")))
		if string(segment) == "[DONE]" {
			break
		}

		var streamResponse StreamResponse
		if err := json.Unmarshal(segment, &streamResponse); err != nil {
			return false, nil, "", err
		}
		if len(streamResponse.Choices) == 0 {
			continue
		}

		delta := streamResponse.Choices[0].Delta
		if delta.Content != "" {
			content += delta.Content
			emitStreamEvent(onEvent, "delta", map[string]any{
				"content": delta.Content,
			})
		}

		if delta.ToolCalls != nil {
			mergeToolCalls(toolCallMap, &toolCallOrder, *delta.ToolCalls)
		}
	}

	if err := scanner.Err(); err != nil {
		return false, nil, "", err
	}

	toolCalls := orderedToolCalls(toolCallMap, toolCallOrder)
	if len(toolCalls) > 0 {
		emitStreamEvent(onEvent, "tool_call", map[string]any{
			"tool_calls": toolCalls,
		})
		return true, toolCalls, content, nil
	}

	return false, nil, content, nil
}

func mergeToolCalls(toolCallMap map[int]ai.ToolCall, toolCallOrder *[]int, incoming []ai.ToolCall) {
	for _, toolCall := range incoming {
		current, ok := toolCallMap[toolCall.Index]
		if !ok {
			current = ai.ToolCall{
				Index: toolCall.Index,
				Type:  toolCall.Type,
				Id:    toolCall.Id,
				Function: ai.FunctionCall{
					Name:      toolCall.Function.Name,
					Arguments: toolCall.Function.Arguments,
				},
			}
			*toolCallOrder = append(*toolCallOrder, toolCall.Index)
		}

		if toolCall.Id != "" {
			current.Id = toolCall.Id
		}
		if toolCall.Type != "" {
			current.Type = toolCall.Type
		}
		if toolCall.Function.Name != "" {
			current.Function.Name = toolCall.Function.Name
		}
		if toolCall.Function.Arguments != "" {
			switch {
			case current.Function.Arguments == "":
				current.Function.Arguments = toolCall.Function.Arguments
			case strings.HasPrefix(toolCall.Function.Arguments, current.Function.Arguments):
				current.Function.Arguments = toolCall.Function.Arguments
			default:
				current.Function.Arguments += toolCall.Function.Arguments
			}
		}

		toolCallMap[toolCall.Index] = current
	}
}

func orderedToolCalls(toolCallMap map[int]ai.ToolCall, toolCallOrder []int) []ai.ToolCall {
	toolCalls := make([]ai.ToolCall, 0, len(toolCallOrder))
	for _, idx := range toolCallOrder {
		toolCalls = append(toolCalls, toolCallMap[idx])
	}
	return toolCalls
}

func emitStreamEvent(onEvent StreamEventHandler, event string, data any) {
	if onEvent == nil {
		return
	}
	onEvent(event, data)
}
