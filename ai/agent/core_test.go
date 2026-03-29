package agent

import (
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
	result, err := agent.RunAgent("你好，今天上海天气怎么样？")
	if err != nil {
		t.Fatalf("RunAgent failed: %v", err)
	}
	t.Logf("result: %v", result.content)
}
