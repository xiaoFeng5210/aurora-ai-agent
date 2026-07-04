package database

import (
	"aurora-agent/database/model"

	"gorm.io/gorm"
)

type CardQueryFilter struct {
	UserID   int
	Content  string
	Tags     []string
	TagIDs   []int
	Page     int
	PageSize int
}

func CreateCard(card model.Card) (model.Card, error) {
	err := db.Model(&model.Card{}).Create(&card).Error
	return card, err
}

func GetCardByIDAndUserID(id int, userID int) (model.Card, error) {
	var card model.Card
	err := db.Model(&model.Card{}).
		Where("id = ? AND user_id = ?", id, userID).
		First(&card).Error
	return card, err
}

func QueryCardsByUserID(filter CardQueryFilter) ([]model.Card, error) {
	queryDB := db.Model(&model.Card{}).Where("user_id = ?", filter.UserID)

	if filter.Content != "" {
		queryDB = queryDB.Where("content ILIKE ?", "%"+filter.Content+"%")
	}

	if len(filter.Tags) > 0 {
		queryDB = queryDB.Where("tags && ?::text[]", model.StringArray(filter.Tags))
	}

	if len(filter.TagIDs) > 0 {
		queryDB = queryDB.Where(
			"EXISTS (SELECT 1 FROM card_tag ct WHERE ct.card_id = card.id AND ct.user_id = ? AND ct.tag_id IN ? AND ct.deleted_at IS NULL)",
			filter.UserID,
			filter.TagIDs,
		)
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		queryDB = queryDB.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	} else if filter.PageSize > 0 {
		queryDB = queryDB.Limit(filter.PageSize)
	}

	var cards []model.Card
	err := queryDB.Order("create_time DESC, id DESC").Find(&cards).Error
	return cards, err
}

func UpdateCardByIDAndUserID(id int, userID int, updates map[string]any) error {
	result := db.Model(&model.Card{}).
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

func SoftDeleteCardByIDAndUserID(id int, userID int) error {
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Card{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
