package service

import (
	"aurora-agent/handler/dto"
	"aurora-agent/handler/vo"
	"aurora-agent/utils"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"crypto/md5"

	"go.uber.org/zap"
)

var (
	baiduTokenResponse vo.BaiduTokenResponse
	baseUrl = "https://pan.baidu.com"
	headers = map[string]string{
		"User-Agent": "pan.baidu.com",
	}
	capacityUnit float64 = 1024 * 1024 * 1024
)

var logger *zap.Logger

func init() {
	logger = utils.Logger
}


// 预上传
func PrecreateUpload(path string, isdir int) (*vo.PrecreateUploadResponse, error) {
	access_token, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {
		return nil, err
	}
	httpUrl := baseUrl + "/rest/2.0/xpan/file"

	appName := "/aurora-ai-agent" + url.PathEscape("知识库")
	username := "super"  // TODO
	parent := "/apps" + appName + "/" + username
	path = parent + path

	blockList, size, err := HandleFile()
	if err != nil {
		return nil, err
	}

	postBytes, err := json.Marshal(
		dto.PrecreateUploadRequest{
			Method: "precreate",
			AccessToken: access_token,
			Path: path,
			Isdir: isdir,
			Autoinit: 1,
			BlockList: blockList,
			Size: int(size),
		},
	)
	if err != nil {
		return nil, err
	}
	postData := bytes.NewBuffer(postBytes)
	resp, err := http.Post(httpUrl, "application/json", postData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result vo.PrecreateUploadResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.Errno != 0 {
		return nil, errors.New("百度网盘预上传失败: " + strconv.Itoa(result.Errno))
	}
	return &result, nil
}


// TODO 处理文件，这里测试我们使用本地
func HandleFile() (blockList []string, size int64, err error) {
	filePath := "test.md"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}
	size = fileInfo.Size()
	defer file.Close()

	var buffer = make([]byte, 1024 * 1024 * 4)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, 0, err
			}

		md5DataContainer := make([]byte, 16)
		md5Data := md5.Sum(buffer[:n])
		copy(md5DataContainer, md5Data[:])
		block := hex.EncodeToString(md5DataContainer)
		blockList = append(blockList, block)
	}

	return blockList, size, nil		
}



// 获取文件或文件夹列表
func GetBaiduNetworkdiskFileList(dir string) (*vo.BaiduNetworkdiskFileListResponse, error) {
	access_token, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {
		return nil, err
	}

	method := "list"
	url := baseUrl + "/rest/2.0/xpan/file?" + fmt.Sprintf("access_token=%s&method=%s&dir=%s", access_token, method, dir)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result vo.BaiduNetworkdiskFileListResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}



// 获取百度网盘容量
func GetBaiduNetworkdiskCapacity() (*vo.BaiduNetworkdiskCapacityResponse, error) {
	access_token, _ := GetBaiduNetworkdiskTokenFromRedis()
	checkfree := 1
	url := baseUrl + "/api/quota" + fmt.Sprintf("?access_token=%s&checkfree=%d", access_token, checkfree)
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
