package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	Id        int            `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	UserId    int            `gorm:"column:user_id" json:"user_id"`
	Name      string         `gorm:"column:name" json:"name"`
	CreatedAt time.Time      `gorm:"column:create_time" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Tag) TableName() string {
	return "tag"
}
