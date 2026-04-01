package model

import (
	"aurora-agent/ai"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MessageToolCalls []ai.ToolCall

func (mtc MessageToolCalls) Value() (driver.Value, error) {
	if mtc == nil {
		return "[]", nil
	}
	data, err := json.Marshal(mtc)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (mtc *MessageToolCalls) Scan(value any) error {
	if mtc == nil {
		return fmt.Errorf("MessageToolCalls: Scan on nil pointer")
	}

	switch v := value.(type) {
	case nil:
		*mtc = MessageToolCalls{}
		return nil
	case []byte:
		if len(v) == 0 {
			*mtc = MessageToolCalls{}
			return nil
		}
		return json.Unmarshal(v, mtc)
	case string:
		if v == "" {
			*mtc = MessageToolCalls{}
			return nil
		}
		return json.Unmarshal([]byte(v), mtc)
	default:
		return fmt.Errorf("MessageToolCalls: unsupported Scan type %T", value)
	}
}

type Message struct {
	Id         int              `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	MessageId  string           `gorm:"column:message_id" json:"message_id"`
	DocumentId int              `gorm:"column:document_id" json:"document_id"`
	Role       string           `gorm:"column:role" json:"role"`
	Content    string           `gorm:"type:text;column:content" json:"content"`
	ToolCalls  MessageToolCalls `gorm:"type:jsonb;column:tool_calls" json:"tool_calls"`
	CreatedAt  time.Time        `gorm:"column:create_time" json:"created_at"`
	UpdatedAt  time.Time        `gorm:"column:update_time" json:"updated_at"`
	DeletedAt  gorm.DeletedAt   `gorm:"column:deleted_at;index" json:"-"`
}

func (Message) TableName() string {
	return "messages"
}

func NewMessageFromAIMessage(documentID int, messageID string, msg ai.Message) Message {
	toolCalls := make(MessageToolCalls, len(msg.ToolCalls))
	copy(toolCalls, msg.ToolCalls)

	return Message{
		MessageId:  messageID,
		DocumentId: documentID,
		Role:       msg.Role,
		Content:    msg.Content,
		ToolCalls:  toolCalls,
	}
}

func (m Message) ToAIMessage() ai.Message {
	toolCalls := make([]ai.ToolCall, len(m.ToolCalls))
	copy(toolCalls, m.ToolCalls)

	return ai.Message{
		Role:      m.Role,
		Content:   m.Content,
		ToolCalls: toolCalls,
	}
}
