package model

import "time"

type User struct {
	Id         int `gorm:"primaryKey;column:id;autoIncrement"`
	Username   string
	Password   string
	Email      string
	Phone      string
	Birthday   time.Time `gorm:"type:date;column:birthday"`
	UserPrompt string    `gorm:"type:text;column:user_prompt"`
	CreatedAt  time.Time `gorm:"column:create_time"`
	UpdatedAt  time.Time `gorm:"column:update_time"`
}
