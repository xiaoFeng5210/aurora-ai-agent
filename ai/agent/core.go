package agent

import "aurora-agent/ai/llm"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	ToolCallId string `json:"tool_call_id"`
}

type Agent struct {
	llm *llm.GLM
	history []Message
}
