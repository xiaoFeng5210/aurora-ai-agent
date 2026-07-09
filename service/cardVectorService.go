package service

import (
	"aurora-agent/database/aliyunossvector"
	"aurora-agent/database/model"
	"aurora-agent/service/embedding"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func syncCardVectorBestEffort(card model.Card) {
	if _, err := upsertCardVector(card); err != nil {
		logger.Error("sync card vector failed",
			zap.Int("uid", card.UserId),
			zap.Int("card_id", card.Id),
			zap.Error(err),
		)
	}
}

func deleteCardVectorBestEffort(uid int, cardID int) {
	if _, err := aliyunossvector.DeleteCardVector(uid, strconv.Itoa(cardID)); err != nil {
		logger.Error("delete card vector failed",
			zap.Int("uid", uid),
			zap.Int("card_id", cardID),
			zap.Error(err),
		)
	}
}

func upsertCardVector(card model.Card) (*aliyunossvector.UpdateResult, error) {
	text := buildCardVectorText(card.Title, card.Content)
	if text == "" {
		return nil, fmt.Errorf("card content is empty")
	}

	md := &embedding.MdDocument{Content: text}
	chunks := md.Chunk()
	vectors, err := md.Embedding()
	if err != nil {
		return nil, err
	}
	if len(chunks) == 0 || len(vectors) == 0 {
		return nil, fmt.Errorf("card chunks or vectors is empty")
	}

	result, err := aliyunossvector.UpsertCardVector(chunks, vectors, map[string]interface{}{
		"user_id":    card.UserId,
		"card_id":    card.Id,
		"card_title": card.Title,
	})
	if err != nil {
		return nil, err
	}
	logger.Info("upsert card vector success",
		zap.Int("uid", card.UserId),
		zap.Int("card_id", card.Id),
		zap.Any("result", result),
	)
	return result, nil
}

func buildCardVectorText(title, content string) string {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	switch {
	case title != "" && content != "":
		return title + "\n\n" + content
	case content != "":
		return content
	default:
		return title
	}
}

func cardVectorFieldsChanged(updates map[string]any) bool {
	if updates == nil {
		return false
	}
	if _, ok := updates["title"]; ok {
		return true
	}
	if _, ok := updates["content"]; ok {
		return true
	}
	return false
}
