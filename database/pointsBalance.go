package database

import (
	"aurora-agent/database/model"

	"gorm.io/gorm"
)

func CreatePointsBalance(balance model.PointsBalance) (model.PointsBalance, error) {
	err := db.Model(&model.PointsBalance{}).Create(&balance).Error
	return balance, err
}

func GetPointsBalanceByID(id int) (model.PointsBalance, error) {
	var balance model.PointsBalance
	err := db.Model(&model.PointsBalance{}).
		Where("id = ?", id).
		First(&balance).Error
	return balance, err
}

func GetPointsBalanceByUserID(userID int) (model.PointsBalance, error) {
	var balance model.PointsBalance
	err := db.Model(&model.PointsBalance{}).
		Where("user_id = ?", userID).
		First(&balance).Error
	return balance, err
}

func UpdatePointsBalanceByUserID(userID int, updates map[string]any) error {
	result := db.Model(&model.PointsBalance{}).
		Where("user_id = ?", userID).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeletePointsBalanceByUserID(userID int) error {
	result := db.Where("user_id = ?", userID).Delete(&model.PointsBalance{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
