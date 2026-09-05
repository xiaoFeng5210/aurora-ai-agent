package service

import (
	"aurora-agent/ai/agent"
	"aurora-agent/handler/dto"
	"testing"
)

func TestBuildChatMessages(t *testing.T) {
	req := dto.ChatRequest{
		Thinking: dto.ChatThinkingRequest{
			Type: "disabled",
		},
		MaxTokens: 1024,
		Prompt: []dto.ChatPromptRequest{
			{Role: "system", Content: "回答要简洁。"},
			{Role: "user", Content: "我在上海。"},
			{Role: "assistant", Content: "知道了。"},
			{Role: "user", Content: "今天天气怎么样？"},
		},
	}

	messages := buildChatMessages(req)
	if len(messages) != 5 {
		t.Fatalf("expected 5 messages, got %d", len(messages))
	}
	if messages[0].Role != "system" || messages[0].Content != agent.SYSTEM_BASE_PROMPT {
		t.Fatalf("unexpected base system prompt: %#v", messages[0])
	}
	if messages[1].Role != "system" || messages[1].Content != "回答要简洁。" {
		t.Fatalf("unexpected prompt order at index 1: %#v", messages[1])
	}
	if messages[4].Role != "user" || messages[4].Content != "今天天气怎么样？" {
		t.Fatalf("unexpected prompt order at index 4: %#v", messages[4])
	}
}
