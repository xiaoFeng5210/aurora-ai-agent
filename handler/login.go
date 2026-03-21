package handler

import (
	"aurora-agent/handler/dto"
	"aurora-agent/middleware"
	"aurora-agent/service"
	"aurora-agent/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Login(ctx *gin.Context) {
	var req dto.LoginRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	queryUser, err := service.AuthenticateUser(req)
	if err != nil {
		respondWithServiceError(ctx, err)
		return
	}

	header := utils.DefautHeader
	payload := utils.JwtPayload{
		Issue:       "dual_token",
		IssueAt:     time.Now().Unix(),
		Expiration:  time.Now().Add(3 * 24 * time.Hour).Unix(),
		UserDefined: map[string]any{"user_id": strconv.Itoa(queryUser.Id), "user_name": queryUser.Username},
	}

	token, err := utils.GenJWT(header, payload, utils.InitViper("conf", "jwt", "yaml").GetString("secret"))
	if err != nil {
		logger.Error("generate token failed", zap.Error(err))
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.SetCookie(
		middleware.COOKIE_NAME,
		token,
		int(3*24*time.Hour/time.Second),
		"/",
		"localhost",
		false,
		true,
	)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "login success",
	})
}
