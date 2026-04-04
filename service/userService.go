package service

import (
	"aurora-agent/database"
	"aurora-agent/database/model"
	"aurora-agent/handler/dto"
	"aurora-agent/handler/vo"
	"errors"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	minPasswordLength = 6
	defaultPageSize   = 10
	maxPageSize       = 100
	birthdayLayout    = "2006-01-02"
)



func CreateUser(req dto.CreateUserRequest) error {
	username := strings.TrimSpace(req.Username)
	email := strings.TrimSpace(req.Email)
	phone := strings.TrimSpace(req.Phone)

	if username == "" || email == "" {
		return gorm.ErrInvalidData
	}
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validatePassword(req.Password); err != nil {
		return err
	}

	birthday, err := parseBirthday(req.Birthday)
	if err != nil {
		return err
	}

	exists, err := database.UsernameExists(username, 0)
	if err != nil {
		return err
	}
	if exists {
		return vo.ErrUsernameExists
	}

	exists, err = database.EmailExists(email, 0)
	if err != nil {
		return err
	}
	if exists {
		return vo.ErrEmailExists
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	user := model.User{
		Username:   username,
		Password:   hashedPassword,
		Email:      email,
		Phone:      phone,
		Birthday:   birthday,
		UserPrompt: req.UserPrompt,
	}

	return database.CreateUser(user)
}

func AuthenticateUser(req dto.LoginRequest) (model.User, error) {
	email := strings.TrimSpace(req.Email)
	user, err := database.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, vo.ErrInvalidCredentials
		}
		return model.User{}, err
	}

	if err := verifyPassword(user.Password, req.Password); err != nil {
		return model.User{}, vo.ErrInvalidCredentials
	}

	return user, nil
}

func GetAllUsers() ([]dto.UserResponse, error) {
	users, err := database.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return toUserResponses(users), nil
}

func GetUserByID(id int) (dto.UserResponse, error) {
	user, err := database.GetUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponse{}, vo.ErrUserNotFound
		}
		return dto.UserResponse{}, err
	}
	return toUserResponse(user), nil
}

func GetCurrentUser(uid int) (dto.UserResponse, error) {
	return GetUserByID(uid)
}

func QueryUsers(filter dto.QueryUserDTO) ([]dto.UserResponse, error) {
	birthday, err := parseBirthday(filter.Birthday)
	if err != nil {
		return nil, err
	}

	pageSize := filter.PageSize
	switch {
	case pageSize <= 0:
		pageSize = defaultPageSize
	case pageSize > maxPageSize:
		pageSize = maxPageSize
	}

	page := filter.Page
	if page < 0 {
		page = 0
	}

	users, err := database.QueryUsers(database.UserQueryFilter{
		Email:    strings.TrimSpace(filter.Email),
		Username: strings.TrimSpace(filter.Username),
		Phone:    strings.TrimSpace(filter.Phone),
		Birthday: birthday,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	return toUserResponses(users), nil
}

func UpdateCurrentUser(uid int, req dto.UpdateCurrentUserRequest) (dto.UserResponse, error) {
	updates := make(map[string]any)

	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if username == "" {
			return dto.UserResponse{}, gorm.ErrInvalidData
		}
		exists, err := database.UsernameExists(username, uid)
		if err != nil {
			return dto.UserResponse{}, err
		}
		if exists {
			return dto.UserResponse{}, vo.ErrUsernameExists
		}
		updates["username"] = username
	}

	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email == "" {
			return dto.UserResponse{}, gorm.ErrInvalidData
		}
		if err := validateEmail(email); err != nil {
			return dto.UserResponse{}, err
		}
		exists, err := database.EmailExists(email, uid)
		if err != nil {
			return dto.UserResponse{}, err
		}
		if exists {
			return dto.UserResponse{}, vo.ErrEmailExists
		}
		updates["email"] = email
	}

	if req.Phone != nil {
		updates["phone"] = strings.TrimSpace(*req.Phone)
	}

	if req.Birthday != nil {
		birthday, err := parseBirthday(*req.Birthday)
		if err != nil {
			return dto.UserResponse{}, err
		}
		updates["birthday"] = birthday
	}

	if req.UserPrompt != nil {
		updates["user_prompt"] = *req.UserPrompt
	}

	if len(updates) == 0 {
		return dto.UserResponse{}, vo.ErrNoFieldsToUpdate
	}

	if err := database.UpdateUserByID(uid, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponse{}, vo.ErrUserNotFound
		}
		return dto.UserResponse{}, err
	}

	return GetCurrentUser(uid)
}

func ChangeCurrentUserPassword(uid int, req dto.ChangePasswordRequest) error {
	if err := validatePassword(req.NewPassword); err != nil {
		return err
	}

	user, err := database.GetUserById(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vo.ErrUserNotFound
		}
		return err
	}

	if err := verifyPassword(user.Password, req.OldPassword); err != nil {
		return vo.ErrOldPasswordIncorrect
	}

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	if err := database.UpdateUserByID(uid, map[string]any{"password": hashedPassword}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vo.ErrUserNotFound
		}
		return err
	}

	return nil
}

func DeleteCurrentUser(uid int) error {
	if err := database.SoftDeleteUser(uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vo.ErrUserNotFound
		}
		return err
	}
	return nil
}

func parseBirthday(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	birthday, err := time.Parse(birthdayLayout, value)
	if err != nil {
		return nil, vo.ErrBirthdayFormat
	}

	return &birthday, nil
}

func validatePassword(password string) error {
	if len(password) < minPasswordLength {
		return vo.ErrPasswordTooShort
	}
	return nil
}

func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return gorm.ErrInvalidData
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func toUserResponses(users []model.User) []dto.UserResponse {
	resp := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		resp = append(resp, toUserResponse(user))
	}
	return resp
}

func toUserResponse(user model.User) dto.UserResponse {
	birthday := ""
	if user.Birthday != nil {
		birthday = user.Birthday.Format(birthdayLayout)
	}

	return dto.UserResponse{
		Id:         user.Id,
		Username:   user.Username,
		Email:      user.Email,
		Phone:      user.Phone,
		Birthday:   birthday,
		UserPrompt: user.UserPrompt,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}
