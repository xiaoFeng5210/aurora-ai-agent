package points

import (
	"aurora-agent/database"
	"aurora-agent/database/model"
	"errors"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"
)

// GetBalance leaves legacy users without a balance row unchanged.
func GetBalance(userID int) (int, error) {
	if userID <= 0 {
		return 0, gorm.ErrInvalidData
	}
	balance, err := database.GetPointsBalanceByUserID(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return balance.BalanceAfter, err
}

// AddPoints is the business entry point for crediting an existing user.
func AddPoints(userID, amount int, remark string) (model.PointsBalance, error) {
	remark = strings.TrimSpace(remark)
	if userID <= 0 || amount <= 0 || utf8.RuneCountInString(remark) > 255 {
		return model.PointsBalance{}, gorm.ErrInvalidData
	}
	if _, err := database.GetUserById(userID); err != nil {
		return model.PointsBalance{}, err
	}
	return database.IncrementPointsBalance(userID, amount, remark)
}
