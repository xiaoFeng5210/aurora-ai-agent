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


// "payload": {
//                 "text": {
//                     "Kind": {
//                         "StringValue": "小松鼠阿豆发现，松果掉落时总是尖头朝下。它问森林里最老的橡树为什么。橡树说松果重心在底部，落下时自然翻转，就像不倒翁一样。阿豆试着把松果扔进水里， 果然每次都是尖头先沉。它恍然大悟，原来大自然处处藏着重心的秘密，鸟蛋圆润一头尖也是同样的道理，滚动时会转圈而不会滚远，保护了窝里的蛋不掉落。"
//                     }
//                 }
//             },
//             "score": 0.5964849,
//             "version": 5

func QueryQdrantVector(prompt string) ([]*qdrant.ScoredPoint, error) {
	queryVector, err := embedding.Embed(prompt, 1024)
	if err != nil {
		return nil, err
	}
	searchResult, err := qdrant_db.QueryRagVector(queryVector)
	if err != nil {
		logger.Error("Query Qdrant vector failed", zap.Error(err))
		return nil, err
	}
	logger.Info("Query Qdrant vector success", zap.Any("searchResult", searchResult))
	return searchResult, nil
}
