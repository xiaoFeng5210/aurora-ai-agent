package dto

import "time"

type CreateDocumentRequest struct {
	DisplayName string  `json:"display_name" binding:"required"`
	FileName    *string `json:"file_name"`
}

type UpdateDocumentRequest struct {
	DisplayName *string `json:"display_name"`
	FileName    *string `json:"file_name"`
}

type QueryDocumentDTO struct {
	DisplayName string `json:"display_name"`
	FileName    string `json:"file_name"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
}

type DocumentResponse struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	DisplayName string    `json:"display_name"`
	FileName    *string   `json:"file_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
