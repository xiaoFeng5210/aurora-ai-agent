package ai

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	ToolCallId string `json:"tool_call_id"`
}
