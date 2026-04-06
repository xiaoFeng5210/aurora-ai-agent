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

func TestUpload(t *testing.T) {
	precreateInfo, fileData, blockList, err := PrecreateUpload("/test.md", 0)
	if err != nil {
		t.Fatalf("PrecreateUpload failed: %v", err)
	}

	_, err = Upload(precreateInfo, fileData)
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	err = CreateFileOrDir(precreateInfo.Path, 0, precreateInfo.Uploadid, blockList, precreateInfo.Size)
	if err != nil {
		t.Fatalf("CreateFileOrDir failed: %v", err)
	}

	t.Logf("CreateFileOrDir success")
}


func TestHandleFile(t *testing.T) {
	blockList, size, _, err := HandleFile()
	if err != nil {
		t.Fatalf("HandleFile failed: %v", err)
	}
	t.Logf("blockList: %v", blockList)
	t.Logf("size: %d", size)
}
