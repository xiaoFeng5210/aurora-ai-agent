package ai

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	Id       string       `json:"id"`
	Index    int          `json:"index,omitempty"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type Message struct {
	// Internal model continuation; excluded from API responses and stored history.
	ReasoningContent string     `json:"-"`
	Role             string     `json:"role"`
	Content          string     `json:"content"`
	ToolCallId       *string    `json:"tool_call_id,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
}
