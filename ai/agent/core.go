package agent

import (
	"aurora-agent/ai"
	"aurora-agent/ai/llm"
	functioncall "aurora-agent/ai/llm/function-call"
	"fmt"
	"strings"

	utils "aurora-agent/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var SYSTEM_BASE_PROMPT = "你是一个知识库助手，你可以根据用户的问题，从知识库中查询相关信息，并返回给用户。"

var logger *zap.Logger

func init() {
	logger = utils.Logger
}

type Agent struct {
	Id          string // agent的唯一标识
	Llm         llm.Model
	Options     llm.ChatOptions
	History     []ai.Message
	MaxLoop     int
	CurrentLoop int
	ToolCalls   []ai.ToolCall
}

type AgentResultType string

const (
	AgentResultTypeSuccess   AgentResultType = "success"
	AgentResultTypeError     AgentResultType = "failure"
	AgentResultTypeTerminate AgentResultType = "terminate" // 中断
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
	a.Id = strings.Join(strings.Split(uuid.New().String(), "-"), "")
	a.CurrentLoop = 0
	a.MaxLoop = 6
	a.Llm = llm.InitModel()
	opts.Tools = functioncall.RAGTools
	a.Options = opts

	a.History = []ai.Message{
		{
			Role:    "system",
			Content: SYSTEM_BASE_PROMPT,
		},
	}
}

func (a *Agent) RunAgent(ctx *gin.Context, messages []ai.Message, onEvent llm.StreamEventHandler) (AgentResult, error) {
	if a.Llm == nil {
		a.NewAgent()
	}
	if a.MaxLoop <= 0 {
		a.MaxLoop = 6
	}

	a.CurrentLoop = 0
	var usage llm.Usage
	defer func() {
		logger.Info("agent round usage",
			zap.String("agent_id", a.Id),
			zap.Int("rounds", a.CurrentLoop),
			zap.Int("prompt_tokens", usage.PromptTokens),
			zap.Int("completion_tokens", usage.CompletionTokens),
			zap.Int("total_tokens", usage.TotalTokens),
			zap.Int("prompt_cache_hit_tokens", usage.PromptCacheHitTokens),
			zap.Int("prompt_cache_miss_tokens", usage.PromptCacheMissTokens),
			zap.Int("reasoning_tokens", usage.ReasoningTokens),
		)
	}()
	conversation := append([]ai.Message{}, messages...)
	emitAgentEvent(onEvent, "start", map[string]any{
		"model": a.Llm.Name(),
	})
	// The caller supplies the full conversation; avoid duplicating it on reuse.
	a.History = append([]ai.Message{}, messages...)
	a.ToolCalls = []ai.ToolCall{}

	for {
		a.CurrentLoop++
		response, err := a.Llm.Chat(ctx.Request.Context(), conversation, a.Options, onEvent)
		if err != nil {
			logger.Error("model chat failed", zap.Error(err))
			emitAgentEvent(onEvent, "error", map[string]any{
				"message": err.Error(),
			})
			return AgentResult{Result: AgentResultTypeError, message: err.Error()}, err
		}

		usage = usage.Add(response.Usage)
		content, toolCalls := response.Message.Content, response.Message.ToolCalls
		conversation = append(conversation, response.Message)

		if len(toolCalls) == 0 {
			a.History = append(a.History, ai.Message{
				Role:    "assistant",
				Content: content, // 最后的回答内容
			})
			emitAgentEvent(onEvent, "done", map[string]any{
				"content":       content,
				"finish_reason": response.FinishReason,
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

		for _, toolCall := range toolCalls {
			switch toolCall.Type {
			case "function":
				result, runErr := functioncall.RunToolFunction(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
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
				a.ToolCalls = append(a.ToolCalls, toolCall)
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

// RunAgentWithPormpt keeps the legacy entry point on the same tool-loop implementation.
func (a *Agent) RunAgentWithPormpt(ctx *gin.Context, userPrompt string) (AgentResult, error) {
	if a.Llm == nil {
		a.NewAgent()
	}
	messages := append([]ai.Message{}, a.History...)
	messages = append(messages, ai.Message{Role: "user", Content: userPrompt})
	return a.RunAgent(ctx, messages, nil)
}
