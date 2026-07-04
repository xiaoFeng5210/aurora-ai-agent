package model

import (
	"time"

	"gorm.io/gorm"
)

type CardTag struct {
	Id        int            `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	UserId    int            `gorm:"column:user_id" json:"user_id"`
	CardId    int            `gorm:"column:card_id" json:"card_id"`
	TagId     int            `gorm:"column:tag_id" json:"tag_id"`
	CreatedAt time.Time      `gorm:"column:create_time" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (CardTag) TableName() string {
	return "card_tag"
}
