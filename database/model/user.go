package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id         int            `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	Username   string         `json:"username"`
	Password   string         `json:"password"`
	Email      string         `json:"email"`
	Phone      string         `json:"phone"`
	Birthday   *time.Time     `gorm:"type:date;column:birthday" json:"birthday"`
	UserPrompt string         `gorm:"type:text;column:user_prompt" json:"user_prompt"`
	CreatedAt  time.Time      `gorm:"column:create_time" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:update_time" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}
