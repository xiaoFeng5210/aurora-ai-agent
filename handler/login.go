package handler

import (
	"aurora-agent/database"
	"aurora-agent/database/model"
	"aurora-agent/middleware"
	"aurora-agent/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Login(ctx *gin.Context) {
	logger := utils.InitZap("log/zap")
	var user model.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	if user.Email != "" && user.Password != "" {
		queryUser, err := database.GetUserByEmail(user.Email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    -1,
				"message": err.Error(),
			})
			return
		}
		// 查密码不一致。
		if user.Password != queryUser.Password {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    -1,
				"message": "password is incorrect",
			})
			return
		}

		header := utils.DefautHeader
		payload := utils.JwtPayload{
			Issue:       "dual_token",
			IssueAt:     time.Now().Unix(),
			Expiration:  time.Now().Add(3 * 24 * time.Hour).Unix(), //3天过期
			UserDefined: map[string]any{"user_id": strconv.Itoa(queryUser.Id), "user_name": queryUser.Username},
		}

		token, err := utils.GenJWT(header, payload, utils.InitViper("conf", "jwt", "yaml").GetString("secret"))
		if err != nil {
			logger.Error("generate token failed", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "generate token failed",
			})
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
		return
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "email and password are required",
		})
		return
	}

}
