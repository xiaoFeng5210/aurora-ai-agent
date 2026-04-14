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

