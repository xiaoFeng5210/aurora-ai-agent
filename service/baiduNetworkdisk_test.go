package service

import (
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

// func TestGetBaiduNetworkdiskToken(t *testing.T) {
// 	resp, err := GetBaiduNetworkdiskTokenWeb()
// 	if err != nil {
// 		t.Fatalf("GetBaiduNetworkdiskToken failed: %v", err)
// 	}
// 	t.Logf("resp: %s", string(resp))
// }
