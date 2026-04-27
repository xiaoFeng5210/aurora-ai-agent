package qdrant_db

import (
	"aurora-agent/utils"
	"fmt"
	"os"
	"sync"

	"github.com/qdrant/go-client/qdrant"
	"go.uber.org/zap"
)

var (
	qdrantClient *qdrant.Client
	qdrantOnce   sync.Once
)

func QdrantConnect() {
	logger := utils.InitZap("log/zap")
	qdrantOnce.Do(func() {
	host := os.Getenv("QDRANT_HOST")
	apiKey := os.Getenv("QDRANT_API_KEY")
	if host == "" || apiKey == "" {
		fmt.Println("QDRANT_HOST or QDRANT_API_KEY is not set")
		return
	}
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   host,
		Port:   6334,
		APIKey: apiKey,
		UseTLS: true,
	})
	if err != nil {
		logger.Error("Failed to create qdrant client", zap.Error(err))
		panic(err)
	}
	qdrantClient = client
	})
}
