package database

import (
	"aurora-agent/database/model"

	"gorm.io/gorm"
)

type DocumentQueryFilter struct {
	UserID      int
	DisplayName string
	FileName    string
	Page        int
	PageSize    int
}

func CreateDocument(document model.Document) (model.Document, error) {
	err := db.Model(&model.Document{}).Create(&document).Error
	return document, err
}

func GetDocumentByIDAndUserID(id int, userID int) (model.Document, error) {
	var document model.Document
	err := db.Model(&model.Document{}).
		Where("id = ? AND user_id = ?", id, userID).
		First(&document).Error
	return document, err
}

func QueryDocumentsByUserID(filter DocumentQueryFilter) ([]model.Document, error) {
	queryDB := db.Model(&model.Document{}).Where("user_id = ?", filter.UserID)

	if filter.DisplayName != "" {
		queryDB = queryDB.Where("display_name = ?", filter.DisplayName)
	}

	if filter.FileName != "" {
		queryDB = queryDB.Where("file_name = ?", filter.FileName)
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		queryDB = queryDB.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	} else if filter.PageSize > 0 {
		queryDB = queryDB.Limit(filter.PageSize)
	}

	var documents []model.Document
	err := queryDB.Order("create_time DESC, id DESC").Find(&documents).Error
	return documents, err
}

func UpdateDocumentByIDAndUserID(id int, userID int, updates map[string]any) error {
	result := db.Model(&model.Document{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func SoftDeleteDocumentByIDAndUserID(id int, userID int) error {
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Document{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
