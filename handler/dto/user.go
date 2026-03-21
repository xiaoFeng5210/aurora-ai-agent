package dto

import "time"

type CreateUserRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone"`
	Birthday   string `json:"birthday"`
	UserPrompt string `json:"user_prompt"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type QueryUserDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Birthday string `json:"birthday"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type UpdateCurrentUserRequest struct {
	Username   *string `json:"username"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Birthday   *string `json:"birthday"`
	UserPrompt *string `json:"user_prompt"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UserResponse struct {
	Id         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Birthday   string    `json:"birthday"`
	UserPrompt string    `json:"user_prompt"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
