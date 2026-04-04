package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetBaiduNetworkdiskTokenWeb() ([]byte, error) {
	clientId := os.Getenv("BAIDU_NETWORKDISK_CLIENT_ID")
	// deviceId := os.Getenv("BAIDU_NETWORKDISK_DEVICE_ID")
	url := fmt.Sprintf(`
	https://openapi.baidu.com/oauth/2.0/authorize?
response_type=code&
client_id=%s&
redirect_uri=oob&
scope=basic,netdisk
`, clientId)


	url = strings.ReplaceAll(url, "\n", "")
	url = strings.ReplaceAll(url, "\t", "")
  fmt.Println(url)


	resp, err := http.Get(strings.TrimSpace(url))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
