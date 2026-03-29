package service

import (
	"aurora-agent/ai"
	"aurora-agent/ai/agent"
	"aurora-agent/ai/llm"
	"aurora-agent/handler/dto"
)

type ChatStreamEvent struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

func ChatWithGLMStream(req dto.ChatRequest, onSSEEvent func(ChatStreamEvent)) error {
	messages := buildChatMessages(req)
	chatAgent := agent.Agent{}
	chatAgent.NewAgentWithOptions(llm.ChatOptions{
		MaxTokens:    req.MaxTokens,
		Temperature:  req.Temperature,
		TopP:         req.TopP,
		ThinkingType: req.Thinking.Type,
	})

	_, err := chatAgent.RunAgentWithMessages(messages, func(event string, data any) {
		if onSSEEvent == nil {
			return
		}
		onSSEEvent(ChatStreamEvent{
			Event: event,
			Data:  data,
		})
	})

	return err
}

func buildChatMessages(req dto.ChatRequest) []ai.Message {
	messages := []ai.Message{
		{
			Role:    "system",
			Content: agent.SYSTEM_BASE_PROMPT,
		},
	}

	for _, prompt := range req.Prompt {
		messages = append(messages, ai.Message{
			Role:    prompt.Role,
			Content: prompt.Content,
		})
	}

	return messages
}
