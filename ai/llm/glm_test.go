package llm

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../.env")
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

	// glm := InitModel("glm-4.7")
	// glm.SendUserPrompt("你好，今天上海天气怎么样？")

	// glm.ChatWithGLMInStream([]MessageRequest{
	// 	{
	// 		Role:    "user",
	// 		Content: "你好，今天上海天气怎么样？",
	// 	},
	// })

}
