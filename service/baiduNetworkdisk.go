package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type BaiduTokenResponse struct {
	ExpiresIn int `json:"expires_in"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var (
	baiduTokenResponse BaiduTokenResponse
)



func GetBaiduNetworkdiskToken() (*BaiduTokenResponse, error) {
	clientId := os.Getenv("BAIDU_NETWORKDISK_CLIENT_ID")
	code := "3e5596876c24102217425d64e91566bd"
	clientSecret := os.Getenv("BAIDU_NETWORKDISK_CLIENT_SECRET")
	url := fmt.Sprintf(`
	https://openapi.baidu.com/oauth/2.0/token?
grant_type=authorization_code&
code=%s&
client_id=%s&
client_secret=%s&
redirect_uri=oob
`, code, clientId, clientSecret)


	url = strings.ReplaceAll(url, "\n", "")
	url = strings.ReplaceAll(url, "\t", "")


	resp, err := http.Get(strings.TrimSpace(url))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("body: " + string(body))

	
  err = json.Unmarshal(body, &baiduTokenResponse)
	if err != nil {
		return nil, err
	}

	fmt.Println("token: " + baiduTokenResponse.AccessToken)

	return &baiduTokenResponse, nil
}
