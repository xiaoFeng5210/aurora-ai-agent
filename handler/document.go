package handler

import (
	"aurora-agent/database"
	"aurora-agent/handler/dto"
	"aurora-agent/middleware"
	"aurora-agent/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateDocument(ctx *gin.Context) {
	var req dto.CreateDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	document, err := service.CreateDocument(ctx.GetInt(middleware.UID_IN_CTX), req)
	if err != nil {
		logger.Error("create document failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    document,
	})
}

func GetDocumentById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	document, err := service.GetDocumentByID(ctx.GetInt(middleware.UID_IN_CTX), id)
	if err != nil {
		logger.Error("get document by id failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    document,
	})
}

func QueryDocument(ctx *gin.Context) {
	var filter dto.QueryDocumentDTO
	if err := ctx.ShouldBindJSON(&filter); err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	documents, err := service.QueryDocuments(ctx.GetInt(middleware.UID_IN_CTX), filter)
	if err != nil {
		logger.Error("query document failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    documents,
	})
}

func UpdateDocument(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	var req dto.UpdateDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	document, err := service.UpdateDocument(ctx.GetInt(middleware.UID_IN_CTX), id, req)
	if err != nil {
		logger.Error("update document failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    document,
	})
}

func DeleteDocument(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := service.DeleteDocument(ctx.GetInt(middleware.UID_IN_CTX), id); err != nil {
		logger.Error("delete document failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func GetAllDocuments(ctx *gin.Context) {
	documents, err := database.GetAllDocumentsByUserID(ctx.GetInt(middleware.UID_IN_CTX))
	if err != nil {
		logger.Error("get all documents failed", zap.Error(err))
		respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    documents,
	})
}
