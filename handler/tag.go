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

func CreateTag(ctx *gin.Context) {
	var req dto.CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	tag, err := service.CreateTag(ctx.GetInt(middleware.UID_IN_CTX), req)
	if err != nil {
		logger.Error("create tag failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, tag)
}

func GetTagById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	tag, err := service.GetTagByID(ctx.GetInt(middleware.UID_IN_CTX), id)
	if err != nil {
		logger.Error("get tag by id failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, tag)
}

func QueryTag(ctx *gin.Context) {
	var filter dto.QueryTagDTO
	if err := ctx.ShouldBindJSON(&filter); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	tags, err := service.QueryTags(ctx.GetInt(middleware.UID_IN_CTX), filter)
	if err != nil {
		logger.Error("query tag failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, tags)
}

func UpdateTag(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	var req dto.UpdateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	tag, err := service.UpdateTag(ctx.GetInt(middleware.UID_IN_CTX), id, req)
	if err != nil {
		logger.Error("update tag failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, tag)
}

func DeleteTag(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := service.DeleteTag(ctx.GetInt(middleware.UID_IN_CTX), id); err != nil {
		logger.Error("delete tag failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func AddTagToCard(ctx *gin.Context) {
	cardID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	var req dto.AddCardTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	cardTag, err := service.AddTagToCard(ctx.GetInt(middleware.UID_IN_CTX), cardID, req.TagId)
	if err != nil {
		logger.Error("add tag to card failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, cardTag)
}

func ReplaceCardTags(ctx *gin.Context) {
	cardID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	var req dto.ReplaceCardTagsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	cardTags, err := service.ReplaceCardTags(ctx.GetInt(middleware.UID_IN_CTX), cardID, req.TagIds)
	if err != nil {
		logger.Error("replace card tags failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, cardTags)
}

func DeleteTagFromCard(ctx *gin.Context) {
	cardID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	tagID, err := strconv.Atoi(ctx.Param("tag_id"))
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := service.DeleteTagFromCard(ctx.GetInt(middleware.UID_IN_CTX), cardID, tagID); err != nil {
		logger.Error("delete tag from card failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func QueryCardTags(ctx *gin.Context) {
	var filter dto.QueryCardTagDTO
	if err := ctx.ShouldBindJSON(&filter); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	cardTags, err := service.QueryCardTags(ctx.GetInt(middleware.UID_IN_CTX), filter)
	if err != nil {
		logger.Error("query card tags failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, cardTags)
}

func QueryTagsByCard(ctx *gin.Context) {
	cardID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	tags, err := service.QueryTagsByCard(ctx.GetInt(middleware.UID_IN_CTX), cardID)
	if err != nil {
		logger.Error("query tags by card failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}

	vo.RespondSuccess(ctx, tags)
}
