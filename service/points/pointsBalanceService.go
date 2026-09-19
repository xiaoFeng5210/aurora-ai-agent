package points

import (
	"aurora-agent/database"
	"aurora-agent/database/model"
	"aurora-agent/handler/vo"
	"aurora-agent/utils/enum"
	"errors"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"
)

func ensureUser(userID int) error {
	_, err := database.GetUserById(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return vo.ErrUserNotFound
	}
	return err
}

func validateChange(userID, amount int, remark string) error {
	if userID <= 0 || amount <= 0 || utf8.RuneCountInString(remark) > 255 {
		return gorm.ErrInvalidData
	}
	return ensureUser(userID)
}

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

// QueryByUserID returns a user's balance after confirming the user exists.
func QueryByUserID(userID int) (int, error) {
	if userID <= 0 {
		return 0, gorm.ErrInvalidData
	}
	if err := ensureUser(userID); err != nil {
		return 0, err
	}
	return Query(userID)
}

// Add is the business entry point for crediting an existing user.
func Add(userID, amount int, remark string, triggerMode enum.TriggerModeEnum) (model.PointsBalance, error) {
	remark = strings.TrimSpace(remark)
	if err := validateChange(userID, amount, remark); err != nil {
		return model.PointsBalance{}, err
	}

	return database.IncrementPointsBalance(userID, amount, remark, triggerMode)
}

// Deduct is the business entry point for reducing an existing user's points.
func Deduct(userID, amount int, remark string, triggerMode enum.TriggerModeEnum) (model.PointsBalance, error) {
	remark = strings.TrimSpace(remark)
	if err := validateChange(userID, amount, remark); err != nil {
		return model.PointsBalance{}, err
	}

	balance, err := database.DecrementPointsBalance(userID, amount, remark, triggerMode)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.PointsBalance{}, vo.ErrInsufficientPoints
	}
	return balance, err
}
