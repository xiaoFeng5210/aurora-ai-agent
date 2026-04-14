package handler

import (
	"aurora-agent/handler/vo"
	"aurora-agent/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpsertQdrantTestRequest struct {
	Text string `json:"text"`
}

func UpsertQdrantTest(ctx *gin.Context) {
	var req UpsertQdrantTestRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	updateResult, err := service.UpsertQdrantByMDText(req.Text)
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}
	vo.RespondSuccess(ctx, updateResult)
}
