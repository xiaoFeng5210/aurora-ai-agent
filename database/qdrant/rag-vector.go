package qdrant_db

import (
	"context"

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
  if exists {
    return nil
  }
	return qdrantClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: ragCollectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1024,  // Vector size is defined by used model
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

// upsert engine, 每个文本chunk组成一个point
func UpsertRagVector(chunks []string, vectors [][]float32, payload map[string]any) (*qdrant.UpdateResult, error) {
	points := make([]*qdrant.PointStruct, len(chunks))
	// userID := payload["user_id"]
	// fileName := payload["file_name"]
	for idx := range chunks {
		points[idx] = &qdrant.PointStruct{
			Id: qdrant.NewIDNum(uint64(idx)),
			Vectors: qdrant.NewVectors(vectors[idx]...),
			Payload: qdrant.NewValueMap(payload),
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
func QueryRagVector(queryVector []float32) ([]*qdrant.ScoredPoint, error) {
	searchResult, err := qdrantClient.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: ragCollectionName,
		Query:          qdrant.NewQuery(queryVector...),
		WithPayload:    qdrant.NewWithPayload(true),
		// Filter: &qdrant.Filter{
		// 	Must: []*qdrant.Condition{
		// 		qdrant.NewMatch("user_id", strconv.Itoa(userID)),
		// 	},
		// },
	})
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}
