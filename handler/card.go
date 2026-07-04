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

func CreateCard(ctx *gin.Context) {
	var req dto.CreateCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	card, err := service.CreateCard(ctx.GetInt(middleware.UID_IN_CTX), req)
	if err != nil {
		logger.Error("create card failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, card)
}

func GetCardById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	card, err := service.GetCardByID(ctx.GetInt(middleware.UID_IN_CTX), id)
	if err != nil {
		logger.Error("get card by id failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, card)
}

func QueryCard(ctx *gin.Context) {
	var filter dto.QueryCardDTO
	if err := ctx.ShouldBindJSON(&filter); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	cards, err := service.QueryCards(ctx.GetInt(middleware.UID_IN_CTX), filter)
	if err != nil {
		logger.Error("query card failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, cards)
}

func UpdateCard(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	var req dto.UpdateCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	card, err := service.UpdateCard(ctx.GetInt(middleware.UID_IN_CTX), id, req)
	if err != nil {
		logger.Error("update card failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, card)
}

func DeleteCard(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := service.DeleteCard(ctx.GetInt(middleware.UID_IN_CTX), id); err != nil {
		logger.Error("delete card failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
