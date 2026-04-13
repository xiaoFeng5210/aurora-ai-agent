package embedding

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

func TestEmbed(t *testing.T) {
	embedded, err := Embed("一段测试文本", 2048)
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}
	t.Logf("embedded: %v", embedded)
}
