package agent

import (
	"aurora-agent/ai"
	"log"
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

func TestAgent(t *testing.T) {
	agent := Agent{}
	agent.NewAgent()
	result, err := agent.RunAgent([]ai.Message{
		{
			Role:    "user",
			Content: "你好，今天上海天气怎么样？",
		},
	}, nil)
	if err != nil {
		t.Fatalf("RunAgent failed: %v", err)
	} 
	t.Logf("result: %v", result.content)
}
