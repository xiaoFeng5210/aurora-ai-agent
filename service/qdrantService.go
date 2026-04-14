package service

import (
	"aurora-agent/service/embedding"

	"github.com/qdrant/go-client/qdrant"
	"go.uber.org/zap"
)

func UpsertQdrantByMDText(mdText string) (*qdrant.UpdateResult, error) {
	md := &embedding.MdDocument{
		Content: mdText,
	}
	md.Chunk()
	md.Embedding()
	updateResult, err := md.UpsertQdrantVector()
	if err != nil {
		logger.Error("Upsert Qdrant by MD text failed", zap.Error(err))
		return nil, err
	}
	logger.Info("Upsert Qdrant by MD text success", zap.Any("updateResult", updateResult))
	return updateResult, nil
}
