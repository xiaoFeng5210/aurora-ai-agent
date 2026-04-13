package qdrant_db

import (
	"errors"
	"fmt"
	"os"

	"github.com/qdrant/go-client/qdrant"
)


func Connect() (*qdrant.Client, error) {
	host := os.Getenv("QDRANT_HOST")
	apiKey := os.Getenv("QDRANT_API_KEY")
	if host == "" || apiKey == "" {
		return nil, errors.New("QDRANT_HOST or QDRANT_API_KEY is not set")
	}
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   host,
		APIKey: apiKey,
		UseTLS: true,
	})
	if err != nil {
		fmt.Println("Failed to create qdrant client", err)
		return nil, err
	}
	return client, nil
}
