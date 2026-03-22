package handler

import (
	"aurora-agent/handler/dto"
	"aurora-agent/middleware"
	"aurora-agent/service"
	"aurora-agent/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var logger *zap.Logger

func init() {
	logger = utils.InitZap("log/zap")
}

func CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("bind create user request failed", zap.Error(err))
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := service.CreateUser(req); err != nil {
		logger.Error("create user failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func GetAllUsers(ctx *gin.Context) {
	users, err := service.GetAllUsers()
	if err != nil {
		logger.Error("get all users failed", zap.Error(err))
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    users,
	})
}

func GetUserById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := service.GetUserByID(id)
	if err != nil {
		logger.Error("get user by id failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

func QueryUser(ctx *gin.Context) {
	var filter dto.QueryUserDTO
	if err := ctx.ShouldBindJSON(&filter); err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	users, err := service.QueryUsers(filter)
	if err != nil {
		logger.Error("query user failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    users,
	})
}

func GetCurrentUser(ctx *gin.Context) {
	user, err := service.GetCurrentUser(ctx.GetInt(middleware.UID_IN_CTX))
	if err != nil {
		logger.Error("get current user failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

func UpdateCurrentUser(ctx *gin.Context) {
	var req dto.UpdateCurrentUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := service.UpdateCurrentUser(ctx.GetInt(middleware.UID_IN_CTX), req)
	if err != nil {
		logger.Error("update current user failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

func ChangeCurrentUserPassword(ctx *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := service.ChangeCurrentUserPassword(ctx.GetInt(middleware.UID_IN_CTX), req); err != nil {
		logger.Error("change current user password failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func DeleteCurrentUser(ctx *gin.Context) {
	if err := service.DeleteCurrentUser(ctx.GetInt(middleware.UID_IN_CTX)); err != nil {
		logger.Error("delete current user failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.SetCookie(
		middleware.COOKIE_NAME,
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func respondWithServiceError(ctx *gin.Context, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		status = http.StatusNotFound
	case errors.Is(err, service.ErrDocumentNotFound):
		status = http.StatusNotFound
	case errors.Is(err, service.ErrUsernameExists),
		errors.Is(err, service.ErrEmailExists),
		errors.Is(err, service.ErrInvalidCredentials),
		errors.Is(err, service.ErrOldPasswordIncorrect),
		errors.Is(err, service.ErrPasswordTooShort),
		errors.Is(err, service.ErrBirthdayFormat),
		errors.Is(err, service.ErrDocumentDisplayNameRequired),
		errors.Is(err, service.ErrNoFieldsToUpdate),
		errors.Is(err, gorm.ErrInvalidData):
		status = http.StatusBadRequest
	}

	respondError(ctx, status, err)
}

func respondError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, gin.H{
		"code":    -1,
		"message": err.Error(),
	})
}
