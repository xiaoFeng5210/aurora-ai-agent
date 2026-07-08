package service

import (
	"aurora-agent/database/aliyunossvector"
	"aurora-agent/middleware"
	"aurora-agent/service/embedding"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UpsertQdrantByMDText(ctx *gin.Context, mdText string) (*aliyunossvector.UpdateResult, error) {
	uid := ctx.GetInt(middleware.UID_IN_CTX)
	md := &embedding.MdDocument{
		Content: mdText,
	}
	md.Chunk()
	_, err := md.Embedding()
	if err != nil {
		return nil, err
	}
	logger.Info("Upsert vector UID:", zap.Any("uid", uid))
	updateResult, err := md.UpsertQdrantVector(uid, "")
	if err != nil {
		logger.Error("Upsert vector by MD text failed", zap.Error(err))
		return nil, err
	}
	logger.Info("Upsert vector by MD text success", zap.Any("updateResult", updateResult))
	return updateResult, nil
}

func QueryQdrantVector(ctx *gin.Context, prompt string) ([]aliyunossvector.ScoredVector, error) {
	uid := ctx.GetInt(middleware.UID_IN_CTX)
	fmt.Println("Query vector UID:", uid)
	queryVector, err := embedding.Embed(prompt, 1024)
	if err != nil {
		return nil, err
	}
	searchResult, err := aliyunossvector.QueryRagVector(queryVector, uid)
	if err != nil {
		logger.Error("Query vector failed", zap.Error(err))
		return nil, err
	}
	logger.Info("Query vector success", zap.Any("searchResult", searchResult))
	return searchResult, nil
}

func DeleteQdrantVectorByFilenames(uid int, filenames []string) (*aliyunossvector.UpdateResult, error) {
	normalizedFilenames := normalizeQdrantFilenames(filenames)
	if len(normalizedFilenames) == 0 {
		return nil, fmt.Errorf("filename is required")
	}

	result, err := aliyunossvector.DeleteRagVectorByFilenames(uid, normalizedFilenames)
	if err != nil {
		logger.Error("Delete vectors by filenames failed",
			zap.Int("uid", uid),
			zap.Strings("filenames", normalizedFilenames),
			zap.Error(err),
		)
		return nil, err
	}
	logger.Info("Delete vectors by filenames success",
		zap.Int("uid", uid),
		zap.Strings("filenames", normalizedFilenames),
		zap.Any("result", result),
	)
	return result, nil
}

func normalizeQdrantFilenames(filenames []string) []string {
	seen := make(map[string]struct{}, len(filenames))
	normalized := make([]string, 0, len(filenames))
	for _, filename := range filenames {
		filename = strings.TrimSpace(filename)
		if filename == "" {
			continue
		}
		if _, ok := seen[filename]; ok {
			continue
		}
		seen[filename] = struct{}{}
		normalized = append(normalized, filename)
	}
	return normalized
}
