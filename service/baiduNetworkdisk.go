package service

import (
	"aurora-agent/handler/vo"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	baiduTokenResponse vo.BaiduTokenResponse
	baseUrl = "https://pan.baidu.com/api"
	headers = map[string]string{
		"User-Agent": "pan.baidu.com",
	}
	capacityUnit float64 = 1024 * 1024 * 1024
)


// 获取百度网盘容量
func GetBaiduNetworkdiskCapacity() (*vo.BaiduNetworkdiskCapacityResponse, error) {
	access_token, _ := GetBaiduNetworkdiskTokenFromRedis()
	checkfree := 1
	url := baseUrl + "/quota" + fmt.Sprintf("?access_token=%s&checkfree=%d", access_token, checkfree)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var baiduNetworkdiskCapacityResponse vo.BaiduNetworkdiskCapacityResponse
	err = json.Unmarshal(body, &baiduNetworkdiskCapacityResponse)
	if err != nil {
		return nil, err
	}

	baiduNetworkdiskCapacityResponse.Total = float64(baiduNetworkdiskCapacityResponse.Total) / capacityUnit / 1000
	baiduNetworkdiskCapacityResponse.Used = float64(baiduNetworkdiskCapacityResponse.Used) / capacityUnit / 1000
	baiduNetworkdiskCapacityResponse.Free = float64(baiduNetworkdiskCapacityResponse.Free) / capacityUnit / 1000

	return &baiduNetworkdiskCapacityResponse, nil
}


// 获取存着的百度网盘token
func GetBaiduNetworkdiskTokenFromRedis() (string, error) {
	// TODO
	return "121.9c59616adbca06490171624ad5e0144e.Ysba4kKcKqJDh1aOOj-NribdDuYz6M4c9_M0S9Y.sUIB3A", nil
}



func GetBaiduNetworkdiskToken() (*vo.BaiduTokenResponse, error) {
	clientId := os.Getenv("BAIDU_NETWORKDISK_CLIENT_ID")
	code := "78c9184f32a08bd4e54bbcaff2b6e49f"
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
