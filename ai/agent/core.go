package agent

import (
	"aurora-agent/ai"
	"aurora-agent/ai/llm"
	functioncall "aurora-agent/ai/llm/function-call"
	"fmt"

	utils "aurora-agent/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var SYSTEM_BASE_PROMPT = "你是一个知识库助手，你可以根据用户的问题，从知识库中查询相关信息，并返回给用户。"

var logger *zap.Logger

func init() {
	logger = utils.Logger
}

type Agent struct {
	Id         string     // agent的唯一标识
	Llm         *llm.GLM
	History     []ai.Message
	MaxLoop     int
	CurrentLoop int
}

type AgentResultType string

const (
	AgentResultTypeSuccess   AgentResultType = "success"
	AgentResultTypeError     AgentResultType = "failure"
	AgentResultTypeTerminate AgentResultType = "terminate"
)

type AgentResult struct {
	Result  AgentResultType
	message string
	content string
}

func (a *Agent) NewAgent() {
	a.NewAgentWithOptions(llm.ChatOptions{})
}

func (a *Agent) NewAgentWithOptions(opts llm.ChatOptions) {
	// 生成一个随机字符串作为agent的唯一标识
	a.Id = uuid.New().String()
	fmt.Println("agent id: ", a.Id)
	a.CurrentLoop = 0
	a.MaxLoop = 6
	a.Llm = llm.InitModel()
	a.Llm.SetOptions(opts)

	a.History = []ai.Message{
		{
			Role:    "system",
			Content: SYSTEM_BASE_PROMPT,
		},
	}
}



func (a *Agent) RunAgent(messages []ai.Message, onEvent llm.StreamEventHandler) (AgentResult, error) {
	if a.Llm == nil {
		a.NewAgent()
	}
	if a.MaxLoop <= 0 {
		a.MaxLoop = 6
	}

	a.CurrentLoop = 0
	conversation := append([]ai.Message{}, messages...)

	for {
		a.CurrentLoop++
		emitAgentEvent(onEvent, "start", map[string]any{
			"model": a.Llm.Model,
		})

		needToolCall, toolCalls, content, err := a.Llm.ChatWithGLMInStreamWithEvents(conversation, llm.ChatOptions{}, onEvent)
		if err != nil {
			logger.Error("ChatWithGLMInStreamWithEvents failed", zap.Error(err))
			emitAgentEvent(onEvent, "error", map[string]any{
				"message": err.Error(),
			})
			return AgentResult{Result: AgentResultTypeError, message: err.Error()}, err
		}

		if !needToolCall {
			emitAgentEvent(onEvent, "done", map[string]any{
				"content":       content,
				"finish_reason": "stop",
			})
			return AgentResult{Result: AgentResultTypeSuccess, message: "success", content: content}, nil
		}

		if a.CurrentLoop >= a.MaxLoop {
			err = fmt.Errorf("max loop reached")
			emitAgentEvent(onEvent, "error", map[string]any{
				"message": err.Error(),
			})
			return AgentResult{Result: AgentResultTypeTerminate, message: err.Error()}, err
		}

		conversation = append(conversation, ai.Message{
			Role:      "assistant",
			Content:   content,
			ToolCalls: toolCalls,
		})

		for _, toolCall := range toolCalls {
			switch toolCall.Type {
			case "function":
				result, runErr := functioncall.RunToolFunction(toolCall.Function.Name, toolCall.Function.Arguments)
				if runErr != nil {
					logger.Error("RunToolFunction failed", zap.Error(runErr))
					emitAgentEvent(onEvent, "error", map[string]any{
						"message": runErr.Error(),
					})
					return AgentResult{Result: AgentResultTypeError, message: runErr.Error()}, runErr
				}

				resultContent := string(result)
				conversation = append(conversation, ai.Message{
					Role:       "tool",
					Content:    resultContent,
					ToolCallId: &toolCall.Id,
				})
				emitAgentEvent(onEvent, "tool_result", map[string]any{
					"tool_call_id": toolCall.Id,
					"name":         toolCall.Function.Name,
					"content":      resultContent,
				})
			default:
				err = fmt.Errorf("unknown tool call type: %s", toolCall.Type)
				emitAgentEvent(onEvent, "error", map[string]any{
					"message": err.Error(),
				})
				return AgentResult{Result: AgentResultTypeError, message: err.Error()}, err
			}
		}
	}
}

func emitAgentEvent(onEvent llm.StreamEventHandler, event string, data any) {
	if onEvent == nil {
		return
	}
	onEvent(event, data)
}





func (a *Agent) RunAgentWithPormpt(userPrompt string) (AgentResult, error) {
	for {
		a.CurrentLoop++

		sendMessages := append([]ai.Message{}, a.History...)
		sendMessages = append(sendMessages, ai.Message{Role: "user", Content: userPrompt})

		needToolCall, toolCalls, content, err := a.Llm.ChatWithGLMInStream(sendMessages)
		if err != nil {
			logger.Error("ChatWithGLMInStream failed", zap.Error(err))
			return AgentResult{Result: AgentResultTypeError, message: err.Error()}, err
		}

		if !needToolCall {
			return AgentResult{Result: AgentResultTypeSuccess, message: "success", content: content}, nil
		}

		if a.CurrentLoop >= a.MaxLoop {
			return AgentResult{Result: AgentResultTypeTerminate, content: "", message: "max loop reached"}, fmt.Errorf("max loop reached")
		}

		for _, toolCall := range toolCalls {
			switch toolCall.Type {
			case "function":
				functionName := toolCall.Function.Name
				functionArguments := toolCall.Function.Arguments
				result, err := functioncall.RunToolFunction(functionName, functionArguments)
				if err != nil {
					logger.Error("RunToolFunction failed", zap.Error(err))
					return AgentResult{Result: AgentResultTypeError, message: err.Error()}, err
				}
				a.History = append(a.History, ai.Message{
					Role:       "tool",
					Content:    string(result),
					ToolCallId: &toolCall.Id,
				})
			default:
				logger.Error("unknown tool call type", zap.String("type", toolCall.Type))
			}
		}
	}
}
