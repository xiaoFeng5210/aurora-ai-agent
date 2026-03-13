package handler

import (
	"aurora-agent/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
