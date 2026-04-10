package service

import (
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestGetBaiduNetworkdiskFileList(t *testing.T) {
	resp, err := GetBaiduNetworkdiskFileList("/oss/")
	if err != nil {
		t.Fatalf("GetBaiduNetworkdiskFileList failed: %v", err)
	}
	if resp.Errno != 0 {
		t.Fatalf("GetBaiduNetworkdiskFileList failed: %v", resp.Errno)
	}
	t.Logf("resp: %v", resp.List)
}
