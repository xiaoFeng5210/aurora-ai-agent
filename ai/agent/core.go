package agent

import (
	"aurora-agent/ai"
	"aurora-agent/ai/llm"
	functioncall "aurora-agent/ai/llm/function-call"
	"fmt"

	utils "aurora-agent/utils"

	"go.uber.org/zap"
)

var SYSTEM_BASE_PROMPT = "你是一个知识库助手，你可以根据用户的问题，从知识库中查询相关信息，并返回给用户。"

var logger *zap.Logger

func init() {
	logger = utils.Logger
}

type Agent struct {
	Llm *llm.GLM
	History []ai.Message
	MaxLoop int
	CurrentLoop int
}

type AgentResultType string
const (
	AgentResultTypeSuccess AgentResultType = "success"
	AgentResultTypeError AgentResultType = "failure"
	// 终止
	AgentResultTypeTerminate AgentResultType = "terminate"
)

type AgentResult struct {
	Result AgentResultType
	message string
	content string
}

func (a *Agent) NewAgent() {
	a.CurrentLoop = 0
	a.MaxLoop = 6
	a.Llm = llm.InitModel()

	a.History = []ai.Message{}
	a.History = append(a.History, ai.Message{
		Role:    "system",
		Content: SYSTEM_BASE_PROMPT,
	})
}


func (a *Agent) RunAgent(userPrompt string) (AgentResult, error) {
	for {
		a.CurrentLoop++

		sendMessages := []ai.Message{}
		sendMessages = append(sendMessages, a.History...)
		sendMessages = append(sendMessages, ai.Message{ Role:    "user", Content: userPrompt })

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
				// 添加工具调用结果到历史记录
				a.History = append(a.History, ai.Message{
					Role:    "tool",
					Content: string(result),
					ToolCallId: &toolCall.Id,
				})

			default:
				logger.Error("unknown tool call type", zap.String("type", string(toolCall.Type)))
			}
		}
	}
}
