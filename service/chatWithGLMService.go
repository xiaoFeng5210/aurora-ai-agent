package service

import (
	"aurora-agent/ai"
	"aurora-agent/ai/agent"
	"aurora-agent/ai/llm"
	"aurora-agent/database"
	"aurora-agent/database/model"
	"aurora-agent/handler/dto"
	"strings"

	"github.com/google/uuid"
)

type ChatStreamEvent struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

func ChatWithGLMStream(documentID int, req dto.ChatRequest, onSSEEvent func(ChatStreamEvent)) error {
	messages := buildChatMessages(req)
	chatAgent := agent.Agent{}
	chatAgent.NewAgentWithOptions(llm.ChatOptions{
		MaxTokens:    req.MaxTokens,
		Temperature:  req.Temperature,
		TopP:         req.TopP,
		ThinkingType: req.Thinking.Type,
	})

	agentResult, err := chatAgent.RunAgent(messages, func(event string, data any) {
		if onSSEEvent == nil {
			return
		}
		onSSEEvent(ChatStreamEvent{
			Event: event,
			Data:  data,
		})
	})

	err = saveChatHistory(documentID, &chatAgent, agentResult)
	if err != nil {
		return err
	}
	return err
}


// 将聊天记录保存到数据库
func saveChatHistory(documentID int, chatAgent *agent.Agent, agentResult agent.AgentResult) error {
	var userMessages []ai.Message
	var filterMessages []ai.Message
	var toolCalls []ai.ToolCall
	for _, toolCall := range chatAgent.ToolCalls {
		toolCalls = append(toolCalls, toolCall)
	}
	for _, message := range chatAgent.History {
		if message.Role == "assistant" {
			message.ToolCalls = toolCalls
		}
		if message.Role == "user" {
			userMessages = append(userMessages, message)
		}
		if message.Role == "assistant" || message.Role == "user" {
			filterMessages = append(filterMessages, message)
		}
	}
	

	var submitDBMessages []model.Message
	for _, message := range filterMessages {
		currentMessageId := strings.Join(strings.Split(uuid.New().String(), "-"), "")
		submitDBMessages = append(submitDBMessages, model.Message{
			MessageId: currentMessageId,
			DocumentId: documentID,
			Role: message.Role,
			Content: message.Content,
			ToolCalls: message.ToolCalls,
		})
	}

	var userMessagesDB []model.Message
	for _, message := range userMessages {
		userMessagesDB = append(userMessagesDB, model.Message{
			MessageId: strings.Join(strings.Split(uuid.New().String(), "-"), ""),
			DocumentId: documentID,
			Role: message.Role,
			Content: message.Content,
		})
	}

	switch agentResult.Result {
		case agent.AgentResultTypeSuccess:
			err := database.BatchCreateMessages(submitDBMessages)
			if err != nil {
		    return err
			}
			return nil
		case agent.AgentResultTypeTerminate:
			err := database.BatchCreateMessages(userMessagesDB)
			if err != nil {
				return err
			}
			return nil
		default:
			err := database.BatchCreateMessages(submitDBMessages)
			if err != nil {
				return err
			}
			return nil
	}
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
