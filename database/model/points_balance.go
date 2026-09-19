package model

import "time"

type PointsBalance struct {
	Id           int       `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	UserId       int       `gorm:"column:user_id;uniqueIndex:idx_points_balance_user_id" json:"user_id"`
	BalanceAfter int       `gorm:"column:balance_after" json:"balance_after"`
	Remark       string    `gorm:"column:remark" json:"remark"`
	CreatedAt    time.Time `gorm:"column:create_time" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:update_time" json:"updated_at"`
}

func (PointsBalance) TableName() string {
	return "points_balance"
}
