package database

import (
	"aurora-agent/database/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpsertRagVectorizationStatus(userID int, fileName string, filePath string, status string, errorMessage *string) error {
	record := model.RagVectorization{
		UserId:       userID,
		FileName:     fileName,
		FilePath:     filePath,
		Status:       status,
		ErrorMessage: errorMessage,
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "file_path"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"file_name":     fileName,
			"status":        status,
			"error_message": errorMessage,
			"deleted_at":    nil,
			"update_time":   gorm.Expr("CURRENT_TIMESTAMP"),
		}),
	}).Create(&record).Error
}

func QueryRagVectorizationsByPaths(userID int, filePaths []string) (map[string]model.RagVectorization, error) {
	statuses := make(map[string]model.RagVectorization, len(filePaths))
	if len(filePaths) == 0 {
		return statuses, nil
	}

	var records []model.RagVectorization
	err := db.Model(&model.RagVectorization{}).
		Where("user_id = ? AND file_path IN ?", userID, filePaths).
		Find(&records).Error
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		statuses[record.FilePath] = record
	}
	return statuses, nil
}

func SoftDeleteRagVectorizationsByPaths(userID int, filePaths []string) error {
	if len(filePaths) == 0 {
		return nil
	}

	result := db.Where("user_id = ? AND file_path IN ?", userID, filePaths).
		Delete(&model.RagVectorization{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
