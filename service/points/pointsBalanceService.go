package points

import (
	"aurora-agent/database"
	"aurora-agent/database/model"
	"errors"
	"strings"
	"unicode/utf8"

	"aurora-agent/utils/enum"

	"gorm.io/gorm"
)

// Query leaves legacy users without a balance row unchanged.
func Query(userID int) (int, error) {
	if userID <= 0 {
		return 0, gorm.ErrInvalidData
	}
	balance, err := database.GetPointsBalanceByUserID(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return balance.BalanceAfter, err
}

// Add is the business entry point for crediting an existing user.
func Add(userID, amount int, remark string, triggerMode enum.TriggerModeEnum) (model.PointsBalance, error) {
	remark = strings.TrimSpace(remark)
	if userID <= 0 || amount <= 0 || utf8.RuneCountInString(remark) > 255 {
		return model.PointsBalance{}, gorm.ErrInvalidData
	}

	return database.IncrementPointsBalance(userID, amount, remark, triggerMode)
}
