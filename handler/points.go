package handler

import (
	"aurora-agent/handler/vo"
	"aurora-agent/middleware"
	"aurora-agent/service/points"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCurrentUserPoints(ctx *gin.Context) {
	balance, err := points.GetBalance(ctx.GetInt(middleware.UID_IN_CTX))
	if err != nil {
		vo.RespondWithServiceError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{"balance": balance},
	})
}
