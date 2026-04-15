package handler

import (
	"aurora-agent/handler/dto"
	"aurora-agent/handler/vo"
	"aurora-agent/middleware"
	"aurora-agent/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func parseHistoryDocumentID(ctx *gin.Context) (int, error) {
	documentID := ctx.Param("id")
	if documentID == "" {
		documentID = ctx.Param("document_id")
	}

	return strconv.Atoi(documentID)
}

func CreateHistoryMessage(ctx *gin.Context) {
	documentID, err := parseHistoryDocumentID(ctx)
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	var req dto.CreateHistoryMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	message, err := service.CreateHistoryMessage(ctx.GetInt(middleware.UID_IN_CTX), documentID, req)
	if err != nil {
		logger.Error("create history message failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, message)
}

func GetHistoryMessageByMessageID(ctx *gin.Context) {
	documentID, err := parseHistoryDocumentID(ctx)
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	message, err := service.GetHistoryMessageByMessageID(ctx.GetInt(middleware.UID_IN_CTX), documentID, ctx.Param("message_id"))
	if err != nil {
		logger.Error("get history message failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, message)
}

func QueryHistoryMessages(ctx *gin.Context) {
	documentID, err := parseHistoryDocumentID(ctx)
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	var req dto.QueryHistoryMessagesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	messages, err := service.QueryHistoryMessages(ctx.GetInt(middleware.UID_IN_CTX), documentID, req)
	if err != nil {
		logger.Error("query history messages failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, messages)
}

func ProxyQueryHistoryMessages(ctx *gin.Context) {
	documentID, err := parseHistoryDocumentID(ctx)
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	messages, err := service.ProxyQueryHistoryMessages(ctx.GetInt(middleware.UID_IN_CTX), documentID)
	if err != nil {
		logger.Error("proxy query history messages failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, messages)
}

func UpdateHistoryMessage(ctx *gin.Context) {
	documentID, err := parseHistoryDocumentID(ctx)
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	var req dto.UpdateHistoryMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	message, err := service.UpdateHistoryMessage(ctx.GetInt(middleware.UID_IN_CTX), documentID, ctx.Param("message_id"), req)
	if err != nil {
		logger.Error("update history message failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, message)
}

func DeleteHistoryMessage(ctx *gin.Context) {
	documentID, err := parseHistoryDocumentID(ctx)
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := service.DeleteHistoryMessage(ctx.GetInt(middleware.UID_IN_CTX), documentID, ctx.Param("message_id")); err != nil {
		logger.Error("delete history message failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
