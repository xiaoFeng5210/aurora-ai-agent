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

// 

// upsert engine, 每个文本chunk组成一个point
func UploadRagVector(chunks []string, userID int, fileName string) {
	points := []*qdrant.PointStruct{}
	for idx, chunk := range chunks {
		points[idx] = &qdrant.PointStruct{
			Id: qdrant.NewIDNum(uint64(idx)),
			Vectors: qdrant.NewVectorsDocument(&qdrant.Document{
				Text:  chunk,
				Model: embeddingModel,
			}),
			Payload: qdrant.NewValueMap(map[string]any{
				"user_id": userID,
				"text": chunk,
				"file_name": fileName,
			}),
		}
	}
	qdrantClient.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: ragCollectionName,
		Points:         points,
	})
}


// query engine
// func QueryRagVector(query string, userID int) {
// 	queryResult, err := qdrantClient.Query(context.Background(), &qdrant.QueryPoints{
// 		CollectionName: ragCollectionName,
// 		Query: qdrant.NewQueryDocument(&qdrant.Document{
// 			Text:  query,
// 			Model: embeddingModel,
// 		}),
// 		Limit: qdrant.PtrOf(uint64(3)),
// 		Filter: &qdrant.Filter{
// 			Must: []*qdrant.Condition{
// 				qdrant.NewRange("user_id", &qdrant.Range{
// 					Match: qdrant.PtrOf(userID),
// 				}),
// 			},
// 		},
// 	})
// }
