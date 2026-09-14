package model

import "time"

type PointsRecord struct {
	Id           int       `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	UserId       int       `gorm:"column:user_id" json:"user_id"`
	Delta        int       `gorm:"column:delta" json:"delta"`
	TriggerMode  string    `gorm:"column:trigger_mode" json:"trigger_mode"`
	BalanceAfter int       `gorm:"column:balance_after" json:"balance_after"`
	Remark       string    `gorm:"column:remark" json:"remark"`
	CreatedAt    time.Time `gorm:"column:create_time" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:update_time" json:"updated_at"`
}

func (PointsRecord) TableName() string {
	return "points_record"
}
