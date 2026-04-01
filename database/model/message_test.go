package model

import (
	"aurora-agent/ai"
	"reflect"
	"testing"
)

func TestMessageToolCallsValue(t *testing.T) {
	var empty MessageToolCalls
	value, err := empty.Value()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "[]" {
		t.Fatalf("expected empty json array, got %#v", value)
	}
}

func TestMessageToolCallsScan(t *testing.T) {
	var toolCalls MessageToolCalls
	err := toolCalls.Scan([]byte(`[{"id":"call_1","type":"function","function":{"name":"weather","arguments":"{}"}}]`))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(toolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(toolCalls))
	}
	if toolCalls[0].Id != "call_1" {
		t.Fatalf("unexpected tool call id: %q", toolCalls[0].Id)
	}
}

func TestNewMessageFromAIMessage(t *testing.T) {
	msg := ai.Message{
		Role:       "assistant",
		Content:    "done",
		ToolCallId: stringPtr("ignored"),
		ToolCalls: []ai.ToolCall{
			{
				Id:   "call_1",
				Type: "function",
				Function: ai.FunctionCall{
					Name:      "weather",
					Arguments: "{}",
				},
			},
		},
	}

	modelMessage := NewMessageFromAIMessage(12, "msg_1", msg)

	if modelMessage.DocumentId != 12 {
		t.Fatalf("expected document id 12, got %d", modelMessage.DocumentId)
	}
	if modelMessage.MessageId != "msg_1" {
		t.Fatalf("expected message id msg_1, got %q", modelMessage.MessageId)
	}
	if modelMessage.Role != "assistant" {
		t.Fatalf("unexpected role: %q", modelMessage.Role)
	}
	if modelMessage.Content != "done" {
		t.Fatalf("unexpected content: %q", modelMessage.Content)
	}
	if len(modelMessage.ToolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(modelMessage.ToolCalls))
	}
}

func TestMessageToAIMessage(t *testing.T) {
	message := Message{
		Role:    "tool",
		Content: "{\"ok\":true}",
		ToolCalls: MessageToolCalls{
			{
				Id:   "call_2",
				Type: "function",
				Function: ai.FunctionCall{
					Name:      "search",
					Arguments: "{\"q\":\"aurora\"}",
				},
			},
		},
	}

	got := message.ToAIMessage()
	want := ai.Message{
		Role:    "tool",
		Content: "{\"ok\":true}",
		ToolCalls: []ai.ToolCall{
			{
				Id:   "call_2",
				Type: "function",
				Function: ai.FunctionCall{
					Name:      "search",
					Arguments: "{\"q\":\"aurora\"}",
				},
			},
		},
	}

	if got.ToolCallId != nil {
		t.Fatalf("expected nil tool call id, got %v", got.ToolCallId)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected ai message: %#v", got)
	}
}

func stringPtr(value string) *string {
	return &value
}
