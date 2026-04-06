package service

import (
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestBaiduNetworkdiskUpload(t *testing.T) {
	err := GMBaiduNetworkdiskUpload("/test.md", 0)
	if err != nil {
		t.Fatalf("BaiduNetworkdiskUpload failed: %v", err)
	}
	t.Logf("BaiduNetworkdiskUpload success")
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


func TestHandleFile(t *testing.T) {
	blockList, size, _, err := HandleFile()
	if err != nil {
		t.Fatalf("HandleFile failed: %v", err)
	}
	t.Logf("blockList: %v", blockList)
	t.Logf("size: %d", size)
}
