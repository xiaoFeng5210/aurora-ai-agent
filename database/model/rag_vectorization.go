package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	RagVectorStatusNotVectorized = "not_vectorized"
	RagVectorStatusVectorizing   = "vectorizing"
	RagVectorStatusCompleted     = "completed"
	RagVectorStatusFailed        = "failed"
)

type RagVectorization struct {
	Id           int            `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	UserId       int            `gorm:"column:user_id" json:"user_id"`
	FileName     string         `gorm:"column:file_name" json:"file_name"`
	FilePath     string         `gorm:"column:file_path" json:"file_path"`
	Status       string         `gorm:"column:status" json:"status"`
	ErrorMessage *string        `gorm:"column:error_message" json:"error_message"`
	CreatedAt    time.Time      `gorm:"column:create_time" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:update_time" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}
