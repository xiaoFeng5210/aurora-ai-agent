package user

import "time"

type UserDTO struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Birthday   time.Time `json:"birthday"`
	UserPrompt string    `json:"user_prompt"`
}
