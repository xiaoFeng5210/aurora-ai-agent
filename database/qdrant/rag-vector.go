package qdrant_db

import (
	"context"

	"github.com/qdrant/go-client/qdrant"
)

// 创建集合
func CreateRagCollection(collectionName string) error {
	err := qdrantClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1024,  // Vector size is defined by used model
			Distance: qdrant.Distance_Cosine,
		}),
	})
	return err
}

// 上传
func UploadRagVector(collectionName string, docs []string) {
	embeddingModel := "sentence-transformers/all-minilm-l6-v2"
	points := []*qdrant.PointStruct{}
	for idx, doc := range docs {
		points[idx] = &qdrant.PointStruct{
			Id: qdrant.NewIDNum(uint64(idx)),
			Vectors: qdrant.NewVectorsDocument(&qdrant.Document{
				Text:  doc,
				Model: embeddingModel,
			}),
			Payload: qdrant.NewValueMap(map[string]any{
				"text": doc,
			}),
		}
	}
	qdrantClient.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
	})
}

