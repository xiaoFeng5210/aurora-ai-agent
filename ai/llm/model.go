package llm

import (
	"aurora-agent/ai"
	"context"
)

// Model is the provider-independent boundary consumed by agents.
// Provider wire formats and stream parsing stay inside the implementation.
type Model interface {
	Name() string
	Chat(context.Context, []ai.Message, ChatOptions, StreamEventHandler) (ChatResult, error)
}

type ChatOptions struct {
	MaxTokens    int
	Temperature  float64
	TopP         float64
	ThinkingType string
	Tools        any
}

type ChatResult struct {
	Message      ai.Message
	FinishReason string
	Usage        Usage
}

// Usage is one model request's token accounting. Agent rounds sum these.
type Usage struct {
	PromptTokens          int
	CompletionTokens      int
	TotalTokens           int
	PromptCacheHitTokens  int
	PromptCacheMissTokens int
	ReasoningTokens       int
}

func (u Usage) Add(other Usage) Usage {
	u.PromptTokens += other.PromptTokens
	u.CompletionTokens += other.CompletionTokens
	u.TotalTokens += other.TotalTokens
	u.PromptCacheHitTokens += other.PromptCacheHitTokens
	u.PromptCacheMissTokens += other.PromptCacheMissTokens
	u.ReasoningTokens += other.ReasoningTokens
	return u
}

type StreamEventHandler func(event string, data any)

const CUR_AI_MODEL_TYPE = "deepseek"

func InitModel() Model {
	switch CUR_AI_MODEL_TYPE {
	case "deepseek":
		return NewDeepSeek()
	default:
		return NewDeepSeek()
	}
}
