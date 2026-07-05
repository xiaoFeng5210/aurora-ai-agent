package model

import (
	"time"

	"gorm.io/gorm"
)

type Card struct {
	Id            int            `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	UserId        int            `gorm:"column:user_id" json:"user_id"`
	Title         string         `gorm:"column:title" json:"title"`
	Content       string         `gorm:"type:text;column:content" json:"content"`
	Tags          StringArray    `gorm:"type:text[];column:tags" json:"tags"`
	ExternalLinks StringArray    `gorm:"type:text[];column:external_links" json:"external_links"`
	InternalLinks StringArray    `gorm:"type:text[];column:internal_links" json:"internal_links"`
	CreatedAt     time.Time      `gorm:"column:create_time" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:update_time" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Card) TableName() string {
	return "card"
}
