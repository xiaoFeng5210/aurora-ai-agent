package database

import (
	"aurora-agent/database/model"

	"gorm.io/gorm"
)

type TagQueryFilter struct {
	UserID   int
	Name     string
	Page     int
	PageSize int
}

func CreateTag(tag model.Tag) (model.Tag, error) {
	err := db.Model(&model.Tag{}).Create(&tag).Error
	return tag, err
}

func GetTagByIDAndUserID(id int, userID int) (model.Tag, error) {
	var tag model.Tag
	err := db.Model(&model.Tag{}).
		Where("id = ? AND user_id = ?", id, userID).
		First(&tag).Error
	return tag, err
}

func QueryTagsByUserID(filter TagQueryFilter) ([]model.Tag, error) {
	queryDB := db.Model(&model.Tag{}).Where("user_id = ?", filter.UserID)

	if filter.Name != "" {
		queryDB = queryDB.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		queryDB = queryDB.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	} else if filter.PageSize > 0 {
		queryDB = queryDB.Limit(filter.PageSize)
	}

	var tags []model.Tag
	err := queryDB.Order("create_time DESC, id DESC").Find(&tags).Error
	return tags, err
}

func TagNameExists(userID int, name string, excludeID int) (bool, error) {
	var count int64
	queryDB := db.Model(&model.Tag{}).Where("user_id = ? AND name = ?", userID, name)
	if excludeID > 0 {
		queryDB = queryDB.Where("id <> ?", excludeID)
	}

	if err := queryDB.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateTagByIDAndUserID(id int, userID int, updates map[string]any) error {
	result := db.Model(&model.Tag{}).
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

func SoftDeleteTagByIDAndUserID(id int, userID int) error {
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Tag{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
