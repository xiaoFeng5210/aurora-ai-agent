package service

import (
	qdrant_db "aurora-agent/database/qdrant"
	"aurora-agent/middleware"
	"aurora-agent/service/embedding"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/qdrant/go-client/qdrant"
	"go.uber.org/zap"
)

func UpsertQdrantByMDText(ctx *gin.Context, mdText string) (*qdrant.UpdateResult, error) {
	uid := ctx.GetInt(middleware.UID_IN_CTX)
	md := &embedding.MdDocument{
		Content: mdText,
	}
	md.Chunk()
	_, err := md.Embedding()
	if err != nil {
		return nil, err
	}
	logger.Info("Upsert Qdrant UID:", zap.Any("uid", uid))
	updateResult, err := md.UpsertQdrantVector(uid, "")
	if err != nil {
		logger.Error("Upsert Qdrant by MD text failed", zap.Error(err))
		return nil, err
	}
	logger.Info("Upsert Qdrant by MD text success", zap.Any("updateResult", updateResult))
	return updateResult, nil
}

func QueryQdrantVector(ctx *gin.Context, prompt string) ([]*qdrant.ScoredPoint, error) {
	uid := ctx.GetInt(middleware.UID_IN_CTX)
	fmt.Println("Query Qdrant UID:", uid)
	queryVector, err := embedding.Embed(prompt, 1024)
	if err != nil {
		return nil, err
	}
	searchResult, err := qdrant_db.QueryRagVector(queryVector, uid)
	if err != nil {
		logger.Error("Query Qdrant vector failed", zap.Error(err))
		return nil, err
	}
	logger.Info("Query Qdrant vector success", zap.Any("searchResult", searchResult))
	return searchResult, nil
}
