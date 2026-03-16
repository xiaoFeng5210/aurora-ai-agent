package user

type UserDTO struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Birthday   string `json:"birthday"`
	UserPrompt string `json:"user_prompt"`
}
