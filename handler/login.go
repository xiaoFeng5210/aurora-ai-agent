package handler

import (
	"aurora-agent/database/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
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

	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "email and password are required",
		})
		return
	}

}
