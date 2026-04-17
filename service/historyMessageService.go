package service

import (
	"aurora-agent/ai"
	"aurora-agent/database"
	"aurora-agent/database/model"
	"aurora-agent/handler/dto"
	"aurora-agent/handler/vo"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateHistoryMessage(uid int, documentID int, req dto.CreateHistoryMessageRequest) (dto.HistoryMessageResponse, error) {
	if err := ensureHistoryDocumentAccessible(uid, documentID); err != nil {
		return dto.HistoryMessageResponse{}, err
	}

	role, err := normalizeHistoryMessageRole(req.Role)
	if err != nil {
		return dto.HistoryMessageResponse{}, err
	}

	messageID := normalizeHistoryMessageID(req.MessageId)
	message, err := database.CreateMessage(model.Message{
		MessageId:  messageID,
		DocumentId: documentID,
		Role:       role,
		Content:    req.Content,
		ToolCalls:  toMessageToolCalls(req.ToolCalls),
	})
	if err != nil {
		return dto.HistoryMessageResponse{}, err
	}

	return toHistoryMessageResponse(message), nil
}

func GetHistoryMessageByMessageID(uid int, documentID int, messageID string) (dto.HistoryMessageResponse, error) {
	if err := ensureHistoryDocumentAccessible(uid, documentID); err != nil {
		return dto.HistoryMessageResponse{}, err
	}

	messageID = strings.TrimSpace(messageID)
	if messageID == "" {
		return dto.HistoryMessageResponse{}, gorm.ErrInvalidData
	}

	message, err := database.GetMessageByMessageIDAndDocumentID(messageID, documentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.HistoryMessageResponse{}, vo.ErrMessageNotFound
		}
		return dto.HistoryMessageResponse{}, err
	}

	return toHistoryMessageResponse(message), nil
}

func QueryHistoryMessages(uid int, documentID int, req dto.QueryHistoryMessagesRequest) ([]dto.HistoryMessageResponse, error) {
	if err := ensureHistoryDocumentAccessible(uid, documentID); err != nil {
		return nil, err
	}

	order, err := normalizeHistoryMessageOrder(req.Order)
	if err != nil {
		return nil, err
	}
	page, pageSize := normalizePagination(req.Page, req.PageSize)

	messages, err := database.QueryMessagesByDocumentID(database.MessageQueryFilter{
		DocumentID: documentID,
		Role:       strings.TrimSpace(req.Role),
		Keyword:    strings.TrimSpace(req.Keyword),
		Page:       page,
		PageSize:   pageSize,
		Order:      order,
	})
	if err != nil {
		return nil, err
	}

	return toHistoryMessageResponses(messages), nil
}

func ProxyQueryHistoryMessages(uid int, documentID int, req dto.ProxyQueryHistoryMessagesRequest) ([]dto.HistoryMessageResponse, error) {
	if err := ensureHistoryDocumentAccessible(uid, documentID); err != nil {
		return nil, err
	}

	startTime, err := parseConversationTime(req.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := parseConversationTime(req.EndTime)
	if err != nil {
		return nil, err
	}

	if startTime != nil && endTime != nil && startTime.After(*endTime) {
		return nil, vo.ErrMessageTimeRange
	}

	page, pageSize := normalizePagination(req.Page, req.PageSize)

	messages, err := database.GetConversationMessages(database.ConversationMessageFilter{
		DocumentID: documentID,
		StartTime:  startTime,
		EndTime:    endTime,
		Keyword:    strings.TrimSpace(req.Keyword),
		Page:       page,
		PageSize:   pageSize,
		Order:      req.Order,
	})
	if err != nil {
		return nil, err
	}

	return toHistoryMessageResponses(messages), nil
}

func parseConversationTime(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, vo.ErrMessageTimeFormat
	}

	return &t, nil
}

func UpdateHistoryMessage(uid int, documentID int, messageID string, req dto.UpdateHistoryMessageRequest) (dto.HistoryMessageResponse, error) {
	if err := ensureHistoryDocumentAccessible(uid, documentID); err != nil {
		return dto.HistoryMessageResponse{}, err
	}

	messageID = strings.TrimSpace(messageID)
	if messageID == "" {
		return dto.HistoryMessageResponse{}, gorm.ErrInvalidData
	}

	updates, err := buildHistoryMessageUpdates(req)
	if err != nil {
		return dto.HistoryMessageResponse{}, err
	}

	if err := database.UpdateMessageByMessageIDAndDocumentID(messageID, documentID, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.HistoryMessageResponse{}, vo.ErrMessageNotFound
		}
		return dto.HistoryMessageResponse{}, err
	}

	return GetHistoryMessageByMessageID(uid, documentID, messageID)
}

func DeleteHistoryMessage(uid int, documentID int, messageID string) error {
	if err := ensureHistoryDocumentAccessible(uid, documentID); err != nil {
		return err
	}

	messageID = strings.TrimSpace(messageID)
	if messageID == "" {
		return gorm.ErrInvalidData
	}

	if err := database.SoftDeleteMessageByMessageIDAndDocumentID(messageID, documentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vo.ErrMessageNotFound
		}
		return err
	}

	return nil
}

func ensureHistoryDocumentAccessible(uid int, documentID int) error {
	if _, err := database.GetDocumentByIDAndUserID(documentID, uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vo.ErrDocumentNotFound
		}
		return err
	}
	return nil
}

func normalizeHistoryMessageRole(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", vo.ErrMessageRoleRequired
	}
	return value, nil
}

func normalizeHistoryMessageOrder(order string) (string, error) {
	order = strings.TrimSpace(order)
	if order == "" {
		return "asc", nil
	}

	if strings.EqualFold(order, "asc") {
		return "asc", nil
	}
	if strings.EqualFold(order, "desc") {
		return "desc", nil
	}

	return "", vo.ErrMessageOrderInvalid
}

func normalizeHistoryMessageID(messageID *string) string {
	if messageID != nil {
		trimmed := strings.TrimSpace(*messageID)
		if trimmed != "" {
			return trimmed
		}
	}

	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func buildHistoryMessageUpdates(req dto.UpdateHistoryMessageRequest) (map[string]any, error) {
	updates := make(map[string]any)

	if req.Role != nil {
		role, err := normalizeHistoryMessageRole(*req.Role)
		if err != nil {
			return nil, err
		}
		updates["role"] = role
	}

	if req.Content != nil {
		updates["content"] = *req.Content
	}

	if req.ToolCalls != nil {
		updates["tool_calls"] = toMessageToolCalls(*req.ToolCalls)
	}

	if len(updates) == 0 {
		return nil, vo.ErrNoFieldsToUpdate
	}

	return updates, nil
}

func toHistoryMessageResponses(messages []model.Message) []dto.HistoryMessageResponse {
	resp := make([]dto.HistoryMessageResponse, 0, len(messages))
	for _, message := range messages {
		resp = append(resp, toHistoryMessageResponse(message))
	}
	return resp
}

func toHistoryMessageResponse(message model.Message) dto.HistoryMessageResponse {
	return dto.HistoryMessageResponse{
		Id:         message.Id,
		MessageId:  message.MessageId,
		DocumentId: message.DocumentId,
		Role:       message.Role,
		Content:    message.Content,
		ToolCalls:  toAIToolCalls(message.ToolCalls),
		CreatedAt:  message.CreatedAt,
		UpdatedAt:  message.UpdatedAt,
	}
}

func toMessageToolCalls(toolCalls []ai.ToolCall) model.MessageToolCalls {
	if len(toolCalls) == 0 {
		return model.MessageToolCalls{}
	}

	resp := make(model.MessageToolCalls, len(toolCalls))
	copy(resp, toolCalls)
	return resp
}

func toAIToolCalls(toolCalls model.MessageToolCalls) []ai.ToolCall {
	if len(toolCalls) == 0 {
		return []ai.ToolCall{}
	}

	resp := make([]ai.ToolCall, len(toolCalls))
	copy(resp, toolCalls)
	return resp
}
