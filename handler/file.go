package handler

import (
	"aurora-agent/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBaiduNetworkdiskToken(ctx *gin.Context) {
	resp, err := service.GetBaiduNetworkdiskTokenWeb()
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Content-Type", "text/html")
	ctx.Writer.Write(resp)
}
