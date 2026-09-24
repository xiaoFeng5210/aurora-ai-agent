package agent

import (
	"context"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"

	"aurora-agent/ai"
	"aurora-agent/ai/llm"
	"github.com/gin-gonic/gin"
)

type fakeModel struct {
	calls int
	fail  bool
	t     *testing.T
}

func (f *fakeModel) Name() string { return "test-provider" }
func (f *fakeModel) Chat(ctx context.Context, messages []ai.Message, opts llm.ChatOptions, emit llm.StreamEventHandler) (llm.ChatResult, error) {
	f.calls++
	if f.fail {
		return llm.ChatResult{}, errors.New("model unavailable")
	}
	if f.calls == 1 {
		calls := []ai.ToolCall{{Id: "weather", Type: "function", Function: ai.FunctionCall{Name: "get_weather", Arguments: `{"city":"Shanghai"}`}}}
		emit("tool_call", map[string]any{"tool_calls": calls})
		return llm.ChatResult{Message: ai.Message{Role: "assistant", ReasoningContent: "continue", ToolCalls: calls}, FinishReason: "tool_calls", Usage: llm.Usage{PromptTokens: 10, CompletionTokens: 4, TotalTokens: 14, ReasoningTokens: 2}}, nil
	}
	if len(messages) != 3 || messages[1].ReasoningContent != "continue" || messages[2].Role != "tool" || *messages[2].ToolCallId != "weather" {
		f.t.Fatal("lost tool conversation")
	}
	emit("delta", map[string]any{"content": "answer"})
	return llm.ChatResult{Message: ai.Message{Role: "assistant", Content: "answer"}, FinishReason: "stop", Usage: llm.Usage{PromptTokens: 20, CompletionTokens: 6, TotalTokens: 26, PromptCacheHitTokens: 8, PromptCacheMissTokens: 12}}, nil
}

func TestAgentModelBoundary(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/", nil)
	fake := &fakeModel{t: t}
	a := Agent{Llm: fake, MaxLoop: 6}
	var events []string
	result, err := a.RunAgent(ctx, []ai.Message{{Role: "user", Content: "weather"}}, func(event string, data any) { events = append(events, event) })
	if err != nil || result.Result != AgentResultTypeSuccess || result.content != "answer" {
		t.Fatalf("agent failed: %v", err)
	}
	if !reflect.DeepEqual(events, []string{"start", "tool_call", "tool_result", "delta", "done"}) {
		t.Fatalf("events: %v", events)
	}
	fake.fail = true
	events = nil
	result, err = a.RunAgent(ctx, []ai.Message{{Role: "user", Content: "weather"}}, func(event string, data any) { events = append(events, event) })
	if err == nil || result.Result != AgentResultTypeError || !reflect.DeepEqual(events, []string{"start", "error"}) {
		t.Fatal("model failure masked")
	}
}
