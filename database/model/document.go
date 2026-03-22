package model

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	Id          int            `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	UserId      int            `gorm:"column:user_id" json:"user_id"`
	DisplayName string         `gorm:"column:display_name" json:"display_name"`
	FileName    *string        `gorm:"column:file_name" json:"file_name"`
	CreatedAt   time.Time      `gorm:"column:create_time" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:update_time" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}
