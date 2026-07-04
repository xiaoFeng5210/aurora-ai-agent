package database

import (
	"aurora-agent/database/model"

	"gorm.io/gorm"
)

type CardTagQueryFilter struct {
	UserID   int
	CardID   int
	TagID    int
	Page     int
	PageSize int
}

func CreateCardTag(cardTag model.CardTag) (model.CardTag, error) {
	var existing model.CardTag
	err := db.Unscoped().
		Where("user_id = ? AND card_id = ? AND tag_id = ?", cardTag.UserId, cardTag.CardId, cardTag.TagId).
		First(&existing).Error
	if err == nil {
		if existing.DeletedAt.Valid {
			err = db.Unscoped().Model(&model.CardTag{}).
				Where("id = ?", existing.Id).
				Updates(map[string]any{"deleted_at": nil}).Error
			if err != nil {
				return model.CardTag{}, err
			}
			return GetCardTagByIDAndUserID(existing.Id, existing.UserId)
		}
		return existing, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return model.CardTag{}, err
	}

	err = db.Model(&model.CardTag{}).Create(&cardTag).Error
	return cardTag, err
}

func GetCardTagByIDAndUserID(id int, userID int) (model.CardTag, error) {
	var cardTag model.CardTag
	err := db.Model(&model.CardTag{}).
		Where("id = ? AND user_id = ?", id, userID).
		First(&cardTag).Error
	return cardTag, err
}

func QueryCardTagsByUserID(filter CardTagQueryFilter) ([]model.CardTag, error) {
	queryDB := db.Model(&model.CardTag{}).Where("user_id = ?", filter.UserID)

	if filter.CardID > 0 {
		queryDB = queryDB.Where("card_id = ?", filter.CardID)
	}

	if filter.TagID > 0 {
		queryDB = queryDB.Where("tag_id = ?", filter.TagID)
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		queryDB = queryDB.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	} else if filter.PageSize > 0 {
		queryDB = queryDB.Limit(filter.PageSize)
	}

	var cardTags []model.CardTag
	err := queryDB.Order("create_time DESC, id DESC").Find(&cardTags).Error
	return cardTags, err
}

func ReplaceCardTagsByCardIDAndUserID(userID int, cardID int, tagIDs []int) ([]model.CardTag, error) {
	var cardTags []model.CardTag

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND card_id = ?", userID, cardID).Delete(&model.CardTag{}).Error; err != nil {
			return err
		}

		cardTags = make([]model.CardTag, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			cardTag := model.CardTag{
				UserId: userID,
				CardId: cardID,
				TagId:  tagID,
			}

			var existing model.CardTag
			err := tx.Unscoped().
				Where("user_id = ? AND card_id = ? AND tag_id = ?", userID, cardID, tagID).
				First(&existing).Error
			switch {
			case err == nil:
				if err := tx.Unscoped().Model(&model.CardTag{}).
					Where("id = ?", existing.Id).
					Updates(map[string]any{"deleted_at": nil}).Error; err != nil {
					return err
				}
				if err := tx.Model(&model.CardTag{}).Where("id = ?", existing.Id).First(&existing).Error; err != nil {
					return err
				}
				cardTags = append(cardTags, existing)
			case err == gorm.ErrRecordNotFound:
				if err := tx.Model(&model.CardTag{}).Create(&cardTag).Error; err != nil {
					return err
				}
				cardTags = append(cardTags, cardTag)
			default:
				return err
			}
		}
		return nil
	})

	return cardTags, err
}

func SoftDeleteCardTagByCardIDAndTagID(userID int, cardID int, tagID int) error {
	result := db.Where("user_id = ? AND card_id = ? AND tag_id = ?", userID, cardID, tagID).
		Delete(&model.CardTag{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func SoftDeleteCardTagsByCardIDAndUserID(userID int, cardID int) error {
	return db.Where("user_id = ? AND card_id = ?", userID, cardID).
		Delete(&model.CardTag{}).Error
}

func SoftDeleteCardTagsByTagIDAndUserID(userID int, tagID int) error {
	return db.Where("user_id = ? AND tag_id = ?", userID, tagID).
		Delete(&model.CardTag{}).Error
}
