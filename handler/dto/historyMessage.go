package dto

import (
	"aurora-agent/ai"
	"time"
)

type CreateHistoryMessageRequest struct {
	MessageId *string       `json:"message_id"`
	Role      string        `json:"role" binding:"required"`
	Content   string        `json:"content"`
	ToolCalls []ai.ToolCall `json:"tool_calls"`
}

type UpdateHistoryMessageRequest struct {
	Role      *string        `json:"role"`
	Content   *string        `json:"content"`
	ToolCalls *[]ai.ToolCall `json:"tool_calls"`
}

type QueryHistoryMessagesRequest struct {
	Role     string `json:"role"`
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Order    string `json:"order"`
}

type ProxyQueryHistoryMessagesRequest struct {
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Keyword   string `form:"keyword"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

type HistoryMessageResponse struct {
	Id         int           `json:"id"`
	MessageId  string        `json:"message_id"`
	DocumentId int           `json:"document_id"`
	Role       string        `json:"role"`
	Content    string        `json:"content"`
	ToolCalls  []ai.ToolCall `json:"tool_calls"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}
