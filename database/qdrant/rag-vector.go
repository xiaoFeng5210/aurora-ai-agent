package qdrant_db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

var ragCollectionName = "rag"
var embeddingModel = "sentence-transformers/all-minilm-l6-v2"


// 创建集合
func CreateRagCollection() error {
	exists, err := qdrantClient.CollectionExists(context.Background(), ragCollectionName)
  if err != nil {
    return err
  }
  if !exists {
		err = qdrantClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
			CollectionName: ragCollectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     1024,  // Vector size is defined by used model
				Distance: qdrant.Distance_Cosine,
			}),
		})
		if err != nil {
			return err
		}
	}
	// 为 user_id 创建 keyword 索引，支持过滤查询（已存在则幂等）
	_, err = qdrantClient.CreateFieldIndex(context.Background(), &qdrant.CreateFieldIndexCollection{
		CollectionName: ragCollectionName,
		FieldName:      "user_id",
		FieldType:      qdrant.PtrOf(qdrant.FieldType_FieldTypeKeyword),
	})
	return err
}

// upsert engine, 每个文本chunk组成一个point
func UpsertRagVector(chunks []string, vectors [][]float32, payload map[string]any) (*qdrant.UpdateResult, error) {
	points := make([]*qdrant.PointStruct, len(chunks))
	userID := fmt.Sprintf("%v", payload["user_id"])
	for idx := range chunks {
		points[idx] = &qdrant.PointStruct{
			Id: qdrant.NewIDUUID(uuid.New().String()),
			Vectors: qdrant.NewVectors(vectors[idx]...),
			Payload: qdrant.NewValueMap(map[string]any{
				"text": chunks[idx],
				"user_id": userID,
			}),
		}
	}
	result, err := qdrantClient.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: ragCollectionName,
		Points:         points,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}


// query engine
func QueryRagVector(queryVector []float32, userID int) ([]*qdrant.ScoredPoint, error) {
	searchResult, err := qdrantClient.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: ragCollectionName,
		Query:          qdrant.NewQuery(queryVector...),
		WithPayload:    qdrant.NewWithPayload(true),
		Filter: &qdrant.Filter{
			Must: []*qdrant.Condition{
				qdrant.NewMatch("user_id", strconv.Itoa(userID)),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}
