package database

import (
	"aurora-agent/database/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IncrementPointsBalance initializes missing rows and increments in one transaction.
// The unique user_id constraint serializes concurrent first-time credits.
func IncrementPointsBalance(userID, amount int, remark string) (model.PointsBalance, error) {
	var balance model.PointsBalance
	err := db.Transaction(func(tx *gorm.DB) error {
		initial := model.PointsBalance{UserId: userID, BalanceAfter: 0, Remark: "积分初始化"}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}}, DoNothing: true,
		}).Create(&initial).Error; err != nil {
			return err
		}
		return tx.Model(&balance).Clauses(clause.Returning{}).
			Where("user_id = ?", userID).
			Updates(map[string]any{
				"balance_after": gorm.Expr("balance_after + ?", amount),
				"remark":        remark,
			}).Error
	})
	return balance, err
}

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
