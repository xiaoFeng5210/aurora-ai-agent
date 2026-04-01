package database

import (
	"aurora-agent/ai"
	"aurora-agent/database/model"
	"strings"

	"gorm.io/gorm"
)

type MessageQueryFilter struct {
	DocumentID int
	Role       string
	Keyword    string
	Page       int
	PageSize   int
	Order      string
}

func CreateMessage(message model.Message) (model.Message, error) {
	err := db.Model(&model.Message{}).Create(&message).Error
	return message, err
}

func BatchCreateMessages(messages []model.Message) error {
	if len(messages) == 0 {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Message{}).Create(&messages).Error
	})
}

func GetMessageByIDAndDocumentID(id int, documentID int) (model.Message, error) {
	var message model.Message
	err := db.Model(&model.Message{}).
		Where("id = ? AND document_id = ?", id, documentID).
		First(&message).Error
	return message, err
}

func GetMessageByMessageIDAndDocumentID(messageID string, documentID int) (model.Message, error) {
	var message model.Message
	err := db.Model(&model.Message{}).
		Where("message_id = ? AND document_id = ?", messageID, documentID).
		First(&message).Error
	return message, err
}

func QueryMessagesByDocumentID(filter MessageQueryFilter) ([]model.Message, error) {
	queryDB := db.Model(&model.Message{}).Where("document_id = ?", filter.DocumentID)

	if filter.Role != "" {
		queryDB = queryDB.Where("role = ?", filter.Role)
	}

	if filter.Keyword != "" {
		queryDB = queryDB.Where("content ILIKE ?", "%"+filter.Keyword+"%")
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		queryDB = queryDB.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	} else if filter.PageSize > 0 {
		queryDB = queryDB.Limit(filter.PageSize)
	}

	var messages []model.Message
	err := queryDB.Order(messageOrderClause(filter.Order)).Find(&messages).Error
	return messages, err
}

func GetConversationMessages(documentID int) ([]model.Message, error) {
	var messages []model.Message
	err := db.Model(&model.Message{}).
		Where("document_id = ?", documentID).
		Order("create_time ASC, id ASC").
		Find(&messages).Error
	return messages, err
}

func UpdateMessageByMessageIDAndDocumentID(messageID string, documentID int, updates map[string]any) error {
	safeUpdates := sanitizeMessageUpdates(updates)
	if len(safeUpdates) == 0 {
		return gorm.ErrInvalidData
	}

	result := db.Model(&model.Message{}).
		Where("message_id = ? AND document_id = ?", messageID, documentID).
		Updates(safeUpdates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func SoftDeleteMessageByMessageIDAndDocumentID(messageID string, documentID int) error {
	result := db.Where("message_id = ? AND document_id = ?", messageID, documentID).
		Delete(&model.Message{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func messageOrderClause(order string) string {
	if strings.EqualFold(order, "desc") {
		return "create_time DESC, id DESC"
	}
	return "create_time ASC, id ASC"
}

func sanitizeMessageUpdates(updates map[string]any) map[string]any {
	if len(updates) == 0 {
		return nil
	}

	safeUpdates := make(map[string]any)
	for key, value := range updates {
		switch key {
		case "role", "content":
			safeUpdates[key] = value
		case "tool_calls":
			if toolCalls, ok := normalizeToolCallsUpdate(value); ok {
				safeUpdates[key] = toolCalls
			}
		}
	}

	return safeUpdates
}

func normalizeToolCallsUpdate(value any) (model.MessageToolCalls, bool) {
	switch v := value.(type) {
	case nil:
		return model.MessageToolCalls{}, true
	case model.MessageToolCalls:
		if v == nil {
			return model.MessageToolCalls{}, true
		}
		return v, true
	case []ai.ToolCall:
		if v == nil {
			return model.MessageToolCalls{}, true
		}
		toolCalls := make(model.MessageToolCalls, len(v))
		copy(toolCalls, v)
		return toolCalls, true
	default:
		return nil, false
	}
}
