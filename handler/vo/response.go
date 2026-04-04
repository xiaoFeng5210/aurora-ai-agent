package vo

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUsernameExists       = errors.New("username already exists")
	ErrEmailExists          = errors.New("email already exists")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrOldPasswordIncorrect = errors.New("old password is incorrect")
	ErrPasswordTooShort     = errors.New("password must be at least 6 characters")
	ErrBirthdayFormat       = errors.New("birthday must be in YYYY-MM-DD format")
	ErrNoFieldsToUpdate     = errors.New("no fields to update")
)

func RespondSuccess(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}


func RespondError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, gin.H{
		"code":    -1,
		"message": err.Error(),
	})
}

func RespondWithServiceError(ctx *gin.Context, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, ErrUserNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrUsernameExists),
		errors.Is(err, ErrEmailExists),
		errors.Is(err, ErrInvalidCredentials),
		errors.Is(err, ErrOldPasswordIncorrect),
		errors.Is(err, ErrPasswordTooShort),
		errors.Is(err, ErrBirthdayFormat),
		errors.Is(err, ErrNoFieldsToUpdate),
		errors.Is(err, gorm.ErrInvalidData):
		status = http.StatusBadRequest
	}

	RespondError(ctx, status, err)
}
