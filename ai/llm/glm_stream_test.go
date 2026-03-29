package llm

import (
	"strings"
	"testing"
)

func TestAIStreamResponseHandlerWithEvents(t *testing.T) {
	stream := strings.NewReader(`
data: {"id":"1","choices":[{"delta":{"role":"assistant","content":"你好"}}]}

data: {"id":"1","choices":[{"delta":{"content":"，上海"}}]}

data: {"id":"1","choices":[{"delta":{"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"get_weather","arguments":"{\"city\":\"上"}}]}}]}

data: {"id":"1","choices":[{"delta":{"tool_calls":[{"index":0,"function":{"arguments":"海\"}"}}]}}]}

data: [DONE]
`)

	glm := InitModel()
	events := make([]string, 0)
	payloads := make([]map[string]any, 0)

	needToolCall, toolCalls, content, err := glm.AIStreamResponseHandlerWithEvents(stream, func(event string, data any) {
		events = append(events, event)
		if payload, ok := data.(map[string]any); ok {
			payloads = append(payloads, payload)
		}
	})
	if err != nil {
		t.Fatalf("AIStreamResponseHandlerWithEvents returned error: %v", err)
	}

	if !needToolCall {
		t.Fatalf("expected needToolCall to be true")
	}
	if content != "你好，上海" {
		t.Fatalf("unexpected content: %s", content)
	}
	if len(toolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(toolCalls))
	}
	if toolCalls[0].Function.Name != "get_weather" {
		t.Fatalf("unexpected tool name: %s", toolCalls[0].Function.Name)
	}
	if toolCalls[0].Function.Arguments != "{\"city\":\"上海\"}" {
		t.Fatalf("unexpected tool arguments: %s", toolCalls[0].Function.Arguments)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	if events[0] != "delta" || events[1] != "delta" || events[2] != "tool_call" {
		t.Fatalf("unexpected events: %#v", events)
	}
	if payloads[0]["content"] != "你好" || payloads[1]["content"] != "，上海" {
		t.Fatalf("unexpected delta payloads: %#v", payloads)
	}
}
