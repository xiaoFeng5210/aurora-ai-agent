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

// func TestGetBaiduNetworkdiskToken(t *testing.T) {
// 	resp, err := GetBaiduNetworkdiskTokenWeb()
// 	if err != nil {
// 		t.Fatalf("GetBaiduNetworkdiskToken failed: %v", err)
// 	}
// 	t.Logf("resp: %s", string(resp))
// }
