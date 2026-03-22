package llm

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Printf("Loaded .env file")
}

func TestGLM(t *testing.T) {
	apiKey := os.Getenv("GLM_API_KEY")
	if apiKey == "" {
		log.Fatalf("GLM_API_KEY is not set")
	}

	glm := &GLM{
		APIKey:   apiKey,
		Model:    "glm-4.7",
		MaxToken: 65536,
		Stream:   true,
	}

	glm.ChatWithGLMInStream()

}
