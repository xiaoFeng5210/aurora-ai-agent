package handler

import (
	"aurora-agent/handler/vo"
	"aurora-agent/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBaiduNetworkdiskToken(ctx *gin.Context) {
	resp, err := service.GetBaiduNetworkdiskToken()
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    *resp,
	})
}
