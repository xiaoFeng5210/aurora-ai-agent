package dto

import "time"

type CreateCardRequest struct {
	Title         string   `json:"title"`
	Content       string   `json:"content" binding:"required"`
	Tags          []string `json:"tags"`
	TagIds        []int    `json:"tag_ids"`
	ExternalLinks []string `json:"external_links"`
	InternalLinks []string `json:"internal_links"`
}

type UpdateCardRequest struct {
	Title         *string   `json:"title"`
	Content       *string   `json:"content"`
	Tags          *[]string `json:"tags"`
	TagIds        *[]int    `json:"tag_ids"`
	ExternalLinks *[]string `json:"external_links"`
	InternalLinks *[]string `json:"internal_links"`
}

type QueryCardDTO struct {
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
	TagIds   []int    `json:"tag_ids"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
}

type CardResponse struct {
	Id            int       `json:"id"`
	UserId        int       `json:"user_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Tags          []string  `json:"tags"`
	TagIds        []int     `json:"tag_ids,omitempty"`
	ExternalLinks []string  `json:"external_links"`
	InternalLinks []string  `json:"internal_links"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
