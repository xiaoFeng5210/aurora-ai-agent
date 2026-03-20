package dto

import "time"

type QueryUserDTO struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Phone string `json:"phone"`
	Birthday time.Time `json:"birthday"`

	Page int `json:"page"`
	PageSize int `json:"page_size"`
}
