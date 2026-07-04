package dto

import "time"

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTagRequest struct {
	Name *string `json:"name"`
}

type QueryTagDTO struct {
	Name     string `json:"name"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type TagResponse struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddCardTagRequest struct {
	TagId int `json:"tag_id" binding:"required"`
}

type ReplaceCardTagsRequest struct {
	TagIds []int `json:"tag_ids"`
}

type QueryCardTagDTO struct {
	CardId   int `json:"card_id"`
	TagId    int `json:"tag_id"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type CardTagResponse struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	CardId    int       `json:"card_id"`
	TagId     int       `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
