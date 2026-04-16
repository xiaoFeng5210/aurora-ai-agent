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

type QueryQdrantVectorRequest struct {
	Prompt string `json:"prompt"`
}

func UpsertQdrantTest(ctx *gin.Context) {
	var req UpsertQdrantTestRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	updateResult, err := service.UpsertQdrantByMDText(ctx, req.Text)
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}
	vo.RespondSuccess(ctx, updateResult)
}


// 查询qdrant 向量
func QueryQdrantVector(ctx *gin.Context) {
	var req QueryQdrantVectorRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}
	searchResult, err := service.QueryQdrantVector(ctx, req.Prompt)
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}
	vo.RespondSuccess(ctx, searchResult)
}
