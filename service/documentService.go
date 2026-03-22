package service

import (
	"aurora-agent/database"
	"aurora-agent/database/model"
	"aurora-agent/handler/dto"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrDocumentNotFound            = errors.New("document not found")
	ErrDocumentDisplayNameRequired = errors.New("display_name is required")
)

func CreateDocument(uid int, req dto.CreateDocumentRequest) (dto.DocumentResponse, error) {
	displayName, err := normalizeDocumentDisplayName(req.DisplayName)
	if err != nil {
		return dto.DocumentResponse{}, err
	}

	document, err := database.CreateDocument(model.Document{
		UserId:      uid,
		DisplayName: displayName,
		FileName:    normalizeOptionalString(req.FileName),
	})
	if err != nil {
		return dto.DocumentResponse{}, err
	}

	return toDocumentResponse(document), nil
}

func GetDocumentByID(uid int, id int) (dto.DocumentResponse, error) {
	document, err := database.GetDocumentByIDAndUserID(id, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.DocumentResponse{}, ErrDocumentNotFound
		}
		return dto.DocumentResponse{}, err
	}

	return toDocumentResponse(document), nil
}

func QueryDocuments(uid int, filter dto.QueryDocumentDTO) ([]dto.DocumentResponse, error) {
	page, pageSize := normalizePagination(filter.Page, filter.PageSize)

	documents, err := database.QueryDocumentsByUserID(database.DocumentQueryFilter{
		UserID:      uid,
		DisplayName: strings.TrimSpace(filter.DisplayName),
		FileName:    strings.TrimSpace(filter.FileName),
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		return nil, err
	}

	return toDocumentResponses(documents), nil
}

func UpdateDocument(uid int, id int, req dto.UpdateDocumentRequest) (dto.DocumentResponse, error) {
	updates, err := buildDocumentUpdates(req)
	if err != nil {
		return dto.DocumentResponse{}, err
	}

	if err := database.UpdateDocumentByIDAndUserID(id, uid, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.DocumentResponse{}, ErrDocumentNotFound
		}
		return dto.DocumentResponse{}, err
	}

	return GetDocumentByID(uid, id)
}

func DeleteDocument(uid int, id int) error {
	if err := database.SoftDeleteDocumentByIDAndUserID(id, uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDocumentNotFound
		}
		return err
	}
	return nil
}

func normalizeDocumentDisplayName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrDocumentDisplayNameRequired
	}
	return value, nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func normalizePagination(page int, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}

	switch {
	case pageSize <= 0:
		pageSize = defaultPageSize
	case pageSize > maxPageSize:
		pageSize = maxPageSize
	}

	return page, pageSize
}

func buildDocumentUpdates(req dto.UpdateDocumentRequest) (map[string]any, error) {
	updates := make(map[string]any)

	if req.DisplayName != nil {
		displayName, err := normalizeDocumentDisplayName(*req.DisplayName)
		if err != nil {
			return nil, err
		}
		updates["display_name"] = displayName
	}

	if req.FileName != nil {
		updates["file_name"] = normalizeOptionalString(req.FileName)
	}

	if len(updates) == 0 {
		return nil, ErrNoFieldsToUpdate
	}

	return updates, nil
}

func toDocumentResponses(documents []model.Document) []dto.DocumentResponse {
	resp := make([]dto.DocumentResponse, 0, len(documents))
	for _, document := range documents {
		resp = append(resp, toDocumentResponse(document))
	}
	return resp
}

func toDocumentResponse(document model.Document) dto.DocumentResponse {
	return dto.DocumentResponse{
		Id:          document.Id,
		UserId:      document.UserId,
		DisplayName: document.DisplayName,
		FileName:    document.FileName,
		CreatedAt:   document.CreatedAt,
		UpdatedAt:   document.UpdatedAt,
	}
}
