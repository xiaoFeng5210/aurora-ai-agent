// 当前窗口上文 历史对话
package agent

import (
	"aurora-agent/ai"
	"aurora-agent/database"
)

type PriorMessages struct {
	DocumentID      int
	MaxLength       int
	HistoryMessages []ai.Message
}

func (p *PriorMessages) pushMessage(message ai.Message) {
	p.HistoryMessages = append(p.HistoryMessages, message)
}

func (p *PriorMessages) limitMessages(limit int) []ai.Message {
	if len(p.HistoryMessages) >= limit {
		p.HistoryMessages = p.HistoryMessages[len(p.HistoryMessages)-limit:]
	}
	return p.HistoryMessages
}

func (p *PriorMessages) Concact(messages []ai.Message) []ai.Message {
	p.HistoryMessages = append(p.HistoryMessages, messages...)
	return p.HistoryMessages
}

func (p *PriorMessages) getHistoryMessagesFromDB() ([]ai.Message, error) {
	messages := make([]ai.Message, 0)
	documentID := p.DocumentID
	filter := database.MessageQueryFilter{
		DocumentID: documentID,
	}
	entries, err := database.QueryMessagesByDocumentID(filter)
	if err != nil {
		return messages, err
	}

	actrualEntries := entries[len(entries)-p.MaxLength:]
	for _, entry := range actrualEntries {
		msg := ai.Message{
			Role:    entry.Role,
			Content: entry.Content,
			// 可根据需要延展其它字段
		}
		p.pushMessage(msg)
	}

	return messages, nil
}

func NewPriorMessages(documentID int) (*PriorMessages, error) {
	p := &PriorMessages{
		DocumentID: documentID,
	}
	_, err := p.getHistoryMessagesFromDB()
	return p, err
}
