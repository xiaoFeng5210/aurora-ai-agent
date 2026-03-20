package handler

import (
	"aurora-agent/database"
	"aurora-agent/database/model"
	"aurora-agent/handler/dto"
	"aurora-agent/handler/model/user"
	"aurora-agent/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger = utils.InitZap("log/zap")
}

func CreateUser(ctx *gin.Context) {
	var userBody user.UserDTO
	if err := ctx.ShouldBindJSON(&userBody); err != nil {
		logger.Error("bind user body failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	t, err := time.Parse("2006-01-02", userBody.Birthday)
	if err != nil {
		logger.Error("parse birthday failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	userDTO := model.User{
		Username:   userBody.Username,
		Password:   userBody.Password,
		Email:      userBody.Email,
		Phone:      userBody.Phone,
		Birthday:   t,
		UserPrompt: userBody.UserPrompt,
	}

	err = database.CreateUser(userDTO)
	if err != nil {
		logger.Error("create user failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func GetAllUsers(ctx *gin.Context) {
	r, err := database.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    r,
	})
}

// 根据ID获取用户信息
func GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)
	user, err := database.GetUserById(idInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}
	users, err := database.QueryUser(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    users,
	})
}
