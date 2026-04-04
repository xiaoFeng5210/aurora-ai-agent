package service

import (
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestGetBaiduNetworkdiskToken(t *testing.T) {
	token, err := GetBaiduNetworkdiskTokenWeb()
	if err != nil {
		t.Fatalf("GetBaiduNetworkdiskToken failed: %v", err)
	}
	t.Logf("token: %s", token)
}
