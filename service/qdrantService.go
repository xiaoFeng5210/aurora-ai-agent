package service

import (
	qdrant_db "aurora-agent/database/qdrant"
	"aurora-agent/service/embedding"

	"github.com/qdrant/go-client/qdrant"
	"go.uber.org/zap"
)

func UpsertQdrantByMDText(mdText string) (*qdrant.UpdateResult, error) {
	md := &embedding.MdDocument{
		Content: mdText,
	}
	md.Chunk()
	_, err := md.Embedding()
	if err != nil {
		return nil, err
	}
	updateResult, err := md.UpsertQdrantVector()
	if err != nil {
		logger.Error("Upsert Qdrant by MD text failed", zap.Error(err))
		return nil, err
	}
	logger.Info("Upsert Qdrant by MD text success", zap.Any("updateResult", updateResult))
	return updateResult, nil
}

func QueryQdrantVector(queryText string) ([]*qdrant.ScoredPoint, error) {
	queryVector, err := embedding.Embed(queryText, 1024)
	if err != nil {
		return nil, err
	}
	searchResult, err := qdrant_db.QueryRagVector(queryVector)
	if err != nil {
		return nil, err
	}
	logger.Info("Query Qdrant vector success", zap.Any("searchResult", searchResult))
	return searchResult, nil
}
