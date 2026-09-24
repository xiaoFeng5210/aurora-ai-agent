package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"aurora-agent/ai"
	"github.com/joho/godotenv"
)

// In-memory HTTP transport exercises the request/response boundary without binding ports.
type testTransport func(*http.Request) (*http.Response, error)

func (f testTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func testClient(handler http.HandlerFunc) *http.Client {
	return &http.Client{Transport: testTransport(func(r *http.Request) (*http.Response, error) {
		recorder := httptest.NewRecorder()
		handler(recorder, r)
		response := recorder.Result()
		return response, nil
	})}
}

func chunk(delta any, finish string) string {
	b, _ := json.Marshal(map[string]any{"choices": []any{map[string]any{"delta": delta, "finish_reason": finish}}})
	return "data: " + string(b) + "\n\n"
}

func TestDeepSeekStream(t *testing.T) {
	stream := ": keep-alive\n\n" + chunk(map[string]any{"reasoning_content": "internal reasoning"}, "") + chunk(map[string]any{"content": "你好"}, "")
	// Repeated fragments must not be mistaken for cumulative snapshots.
	for _, fragment := range []string{`{"text":"`, "a", "a", `"}`} {
		stream += chunk(map[string]any{"tool_calls": []any{map[string]any{"index": 0, "id": "call_1", "type": "function", "function": map[string]any{"arguments": fragment}}}}, "")
	}
	stream += chunk(map[string]any{"tool_calls": []any{map[string]any{"index": 0, "function": map[string]any{"name": "lookup"}}}}, "tool_calls") + "data: [DONE]\n\n"
	var events []string
	result, err := readDeepSeekStream(strings.NewReader(stream), func(event string, data any) { events = append(events, event) })
	if err != nil {
		t.Fatal(err)
	}
	if result.Message.Content != "你好" || result.Message.ReasoningContent != "internal reasoning" || result.Message.ToolCalls[0].Function.Arguments != `{"text":"aa"}` {
		t.Fatalf("unexpected result: %+v", result)
	}
	if !reflect.DeepEqual(events, []string{"delta", "tool_call"}) {
		t.Fatalf("events: %v", events)
	}
	data, _ := json.Marshal(result.Message)
	if strings.Contains(string(data), "internal reasoning") {
		t.Fatal("reasoning leaked into public JSON")
	}
}

func TestDeepSeekStreamFailures(t *testing.T) {
	for name, stream := range map[string]string{
		"empty":          "",
		"disconnected":   chunk(map[string]string{"content": "partial"}, ""),
		"missing done":   chunk(map[string]string{}, "stop"),
		"missing finish": "data: [DONE]\n",
		"invalid JSON":   "data: {\n",
		"api error":      "data: {\"error\":{\"message\":\"private\"}}\n",
		"length":         chunk(map[string]string{}, "length") + "data: [DONE]\n",
		"filter":         chunk(map[string]string{}, "content_filter") + "data: [DONE]\n",
		"resource":       chunk(map[string]string{}, "insufficient_system_resource") + "data: [DONE]\n",
		"missing tool":   chunk(map[string]string{}, "tool_calls") + "data: [DONE]\n",
		"bad arguments":  chunk(map[string]any{"tool_calls": []ai.ToolCall{{Id: "id", Type: "function", Function: ai.FunctionCall{Name: "lookup", Arguments: "{"}}}}, "tool_calls") + "data: [DONE]\n",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := readDeepSeekStream(strings.NewReader(stream), nil)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestDeepSeekToolRoundTrip(t *testing.T) {
	requests := 0
	client := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != "/chat/completions" || r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("wrong endpoint/auth")
		}
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if string(body["model"]) != `"deepseek-v4-flash"` || string(body["stream"]) != "true" || string(body["thinking"]) != `{"type":"enabled"}` {
			t.Error("wrong model/options")
		}
		if string(body["stream_options"]) != `{"include_usage":true}` {
			t.Error("missing usage stream option")
		}
		if body["temperature"] != nil || body["top_p"] != nil {
			t.Error("thinking mode must omit sampling")
		}
		if body["tools"] == nil {
			t.Error("missing tools")
		}
		if requests == 1 {
			fmt.Fprint(w, chunk(map[string]any{"reasoning_content": "use lookup", "tool_calls": []ai.ToolCall{{Id: "call", Index: 2, Type: "function", Function: ai.FunctionCall{Name: "lookup", Arguments: "{}"}}}}, "tool_calls")+"data: [DONE]\n")
		} else {
			var messages []map[string]json.RawMessage
			json.Unmarshal(body["messages"], &messages)
			if string(messages[1]["reasoning_content"]) != `"use lookup"` || string(messages[2]["tool_call_id"]) != `"call"` {
				t.Error("lost continuation")
			}
			if strings.Contains(string(messages[1]["tool_calls"]), "index") {
				t.Error("stream index leaked into request")
			}
			fmt.Fprint(w, chunk(map[string]string{"content": "answer"}, "stop")+"data: [DONE]\n")
		}
	}))
	d := &DeepSeek{apiKey: "test-key", endpoint: "https://test.invalid/chat/completions", client: client}
	messages := []ai.Message{{Role: "user", Content: "question"}}
	opts := ChatOptions{ThinkingType: "enabled", Temperature: 0.5, Tools: []any{map[string]any{"type": "function"}}}
	first, err := d.Chat(context.Background(), messages, opts, nil)
	if err != nil {
		t.Fatal(err)
	}
	id := first.Message.ToolCalls[0].Id
	messages = append(messages, first.Message, ai.Message{Role: "tool", Content: "result", ToolCallId: &id})
	second, err := d.Chat(context.Background(), messages, opts, nil)
	if err != nil {
		t.Fatal(err)
	}
	if second.Message.Content != "answer" || requests != 2 {
		t.Fatal("round trip failed")
	}
}

func TestDeepSeekTransportErrors(t *testing.T) {
	for _, status := range []int{401, 402, 429, 500, 503} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			client := testClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(status); fmt.Fprint(w, "secret-key") }))
			d := &DeepSeek{apiKey: "secret-key", endpoint: "https://test.invalid", client: client}
			_, err := d.Chat(context.Background(), []ai.Message{{Role: "user", Content: "hi"}}, ChatOptions{}, nil)
			if err == nil || !strings.Contains(err.Error(), fmt.Sprint(status)) || strings.Contains(err.Error(), "secret-key") {
				t.Fatalf("unsafe/unexpected error: %v", err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	d := &DeepSeek{apiKey: "test-key", endpoint: "https://api.deepseek.com/chat/completions", client: http.DefaultClient}
	_, err := d.Chat(ctx, []ai.Message{{Role: "user", Content: "hi"}}, ChatOptions{}, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
}

func TestDeepSeekConfig(t *testing.T) {
	t.Setenv("DEEPSEEK_API_KEY", "")
	t.Setenv("GLM_API_KEY", "embedding-only")
	_, err := NewDeepSeek().Chat(context.Background(), nil, ChatOptions{}, nil)
	if err == nil || !strings.Contains(err.Error(), "DEEPSEEK_API_KEY") {
		t.Fatal("must not fall back to embedding key")
	}
	d := &DeepSeek{apiKey: "test"}
	for _, opts := range []ChatOptions{{MaxTokens: -1}, {MaxTokens: 393217}, {ThinkingType: "invalid"}, {Temperature: 3}, {TopP: 2}} {
		if _, err := d.Chat(context.Background(), nil, opts, nil); err == nil {
			t.Fatal("expected options error")
		}
	}
}

// Explicitly opt in: sends only synthetic prompts to the real API; never runs RAG.
func TestDeepSeekLive(t *testing.T) {
	if os.Getenv("DEEPSEEK_LIVE_TEST") != "1" {
		t.Skip("set DEEPSEEK_LIVE_TEST=1 to call the real API")
	}
	if err := godotenv.Load("../../.env"); err != nil && os.Getenv("DEEPSEEK_API_KEY") == "" {
		t.Fatal("missing DeepSeek environment")
	}
	d := NewDeepSeek()
	for _, thinking := range []string{"disabled", "enabled"} {
		t.Run(thinking, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			tools := []any{map[string]any{"type": "function", "function": map[string]any{"name": "lookup_test_value", "description": "Get the test value.", "parameters": map[string]any{"type": "object", "properties": map[string]any{}}}}}
			opts := ChatOptions{MaxTokens: 2048, ThinkingType: thinking, Tools: tools}
			messages := []ai.Message{{Role: "user", Content: "Call lookup_test_value exactly once to retrieve the value. Do not guess it. After the tool returns, reply with just its value."}}
			first, err := d.Chat(ctx, messages, opts, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(first.Message.ToolCalls) != 1 || first.Message.ToolCalls[0].Function.Name != "lookup_test_value" {
				t.Fatal("expected test tool call")
			}
			id := first.Message.ToolCalls[0].Id
			messages = append(messages, first.Message, ai.Message{Role: "tool", Content: "ARIADNE_OK_42", ToolCallId: &id})
			final, err := d.Chat(ctx, messages, opts, nil)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(final.Message.Content, "ARIADNE_OK_42") || len(final.Message.ToolCalls) != 0 {
				t.Fatal("tool answer round trip failed")
			}
			t.Log("DeepSeek V4 Flash streaming + tool continuation passed")
		})
	}
}

func TestDeepSeekMultipleToolsAndUsage(t *testing.T) {
	stream := chunk(map[string]any{"tool_calls": []ai.ToolCall{
		{Index: 1, Id: "second", Type: "function", Function: ai.FunctionCall{Name: "second", Arguments: `{"q":`}},
		{Index: 0, Id: "first", Type: "function", Function: ai.FunctionCall{Name: "first", Arguments: `{}`}},
	}}, "")
	stream += chunk(map[string]any{"tool_calls": []ai.ToolCall{{Index: 1, Function: ai.FunctionCall{Arguments: `"value"}`}}}}, "")
	stream += "data: {\"choices\":[{\"delta\":{},\"finish_reason\":\"tool_calls\"}],\"usage\":{\"total_tokens\":1}}\n\n"
	stream += "data: {\"choices\":[],\"usage\":{\"prompt_tokens\":10,\"completion_tokens\":32,\"total_tokens\":42,\"prompt_cache_hit_tokens\":4,\"prompt_cache_miss_tokens\":6,\"completion_tokens_details\":{\"reasoning_tokens\":7}}}\n\ndata: [DONE]\n"
	result, err := readDeepSeekStream(strings.NewReader(stream), nil)
	if err != nil {
		t.Fatal(err)
	}
	calls := result.Message.ToolCalls
	if len(calls) != 2 || calls[0].Id != "first" || calls[1].Id != "second" || calls[1].Function.Arguments != `{"q":"value"}` {
		t.Fatalf("tools out of order or corrupted: %+v", calls)
	}
	if result.Usage != (Usage{PromptTokens: 10, CompletionTokens: 32, TotalTokens: 42, PromptCacheHitTokens: 4, PromptCacheMissTokens: 6, ReasoningTokens: 7}) {
		t.Fatalf("usage: %+v", result.Usage)
	}
}

func TestDeepSeekNonThinkingRequest(t *testing.T) {
	opts := ChatOptions{Temperature: 0.7}
	d := &DeepSeek{apiKey: "test", endpoint: "https://test.invalid/chat/completions", client: testClient(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if string(body["thinking"]) != `{"type":"disabled"}` || string(body["max_tokens"]) != "65536" || string(body["temperature"]) != "0.7" || body["tools"] != nil {
			t.Fatalf("unexpected default request: %s", body)
		}
		if strings.Contains(string(body["messages"]), "reasoning_content") {
			t.Fatal("non-thinking request leaked reasoning")
		}
		fmt.Fprint(w, chunk(map[string]string{"content": "answer"}, "stop")+"data: [DONE]\n")
	})}
	_, err := d.Chat(context.Background(), []ai.Message{{Role: "assistant", Content: "prior", ReasoningContent: "internal"}, {Role: "user", Content: "hi"}}, opts, nil)
	if err != nil {
		t.Fatal(err)
	}
	if opts.MaxTokens != 0 || opts.ThinkingType != "" {
		t.Fatal("request mutated shared options")
	}
}
