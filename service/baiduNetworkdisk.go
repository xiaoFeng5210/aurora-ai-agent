package service

import (
	"aurora-agent/handler/vo"
	"aurora-agent/utils"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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


func CreateFileOrDir(path string, isdir int, uploadid string, blockList string, size int) error {
	access_token, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {
		return err
	}
	method := "create"

	if isdir == 1 {
		size = 0
	}
	httpUrl := baseUrl + "/rest/2.0/xpan/file" + fmt.Sprintf("?access_token=%s&method=%s", access_token, method)

	formData := url.Values{
		"path": {path},
		"size": {strconv.Itoa(size)},
		"isdir": {strconv.Itoa(isdir)},
		"rtype": {"3"},
		"uploadid": {uploadid},
		"block_list": {blockList},
	}
	encodedFormData := formData.Encode()

	req, err := http.NewRequest(http.MethodPost, httpUrl, strings.NewReader(encodedFormData))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", headers["User-Agent"])
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result vo.CreateFileOrDirResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	fmt.Println("result: " + string(body))
	if result.Errno != 0 {
		return errors.New("百度网盘创建文件或文件夹失败: " + strconv.Itoa(result.Errno))
	}
	return nil
}

// 分片上传
func Upload(precreateInfo *vo.PrecreateUploadResponse, fileData []byte) (*vo.UploadResponse, error) {
	access_token, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             
		return nil, err
	}
	method := "upload"
	path := precreateInfo.Path
	uploadid := precreateInfo.Uploadid
	partseq := 0
	httpUrl := "https://c3.pcs.baidu.com/rest/2.0/pcs/superfile2" + fmt.Sprintf("?access_token=%s&method=%s&path=%s&uploadid=%s&partseq=%d&type=tmpfile", access_token, method, url.QueryEscape(path), uploadid, partseq)

	var postData bytes.Buffer
	writer := multipart.NewWriter(&postData)
	fileWriter, err := writer.CreateFormFile("file", "upload-part")
	if err != nil {
		return nil, err
	}
	if _, err = fileWriter.Write(fileData); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, httpUrl, &postData)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", headers["User-Agent"])
	req.Header.Set("Content-Type", writer.FormDataContentType())
  
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result vo.UploadResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.Errno != 0 {
		return nil, errors.New("百度网盘上传失败: " + strconv.Itoa(result.Errno))
	}
	return &result, nil
}

// 预上传
func PrecreateUpload(path string, isdir int) (*vo.PrecreateUploadResponse, []byte, string, error) {
	access_token, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {
		return nil, nil, "", err
	}
	httpUrl := baseUrl + "/rest/2.0/xpan/file?" + fmt.Sprintf("method=precreate&access_token=%s", access_token)

	appName := "/aurora-ai-agent知识库"
	username := "super"  // TODO
	parent := "/apps" + appName + "/" + username
	path = parent + path

	blockList, size, fileData, err := HandleFile()
	if err != nil {
		return nil, nil, "", err
	}

	blockListBytes, err := json.Marshal(blockList)
	if err != nil {
		return nil, nil, "", err
	}

	formData := url.Values{}
	formData.Set("path", path)
	formData.Set("size", strconv.FormatInt(size, 10))
	formData.Set("isdir", strconv.Itoa(isdir))
	formData.Set("autoinit", "1")
	formData.Set("rtype", "3")
	formData.Set("block_list", string(blockListBytes))

	encodedFormData := formData.Encode()

	req, err := http.NewRequest(http.MethodPost, httpUrl, strings.NewReader(encodedFormData))
	if err != nil {
		return nil, nil, "", err
	}
	req.Header.Set("User-Agent", headers["User-Agent"])
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, "", err
	}
	var result vo.PrecreateUploadResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, nil, "", err
	}
	if result.Errno != 0 {
		return nil, nil, "", errors.New("百度网盘预上传失败: " + strconv.Itoa(result.Errno))
	}
	result.Size = int(size)
	return &result, fileData, string(blockListBytes), nil
}


// TODO 处理文件，这里测试我们使用本地
func HandleFile() (blockList []string, size int64, fileData []byte, err error) {
	filePath := "test.md"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, 0, nil, err
	}
	size = fileInfo.Size()
	defer file.Close()

	if err != nil {
		return nil, 0, nil, err
	}

	fileData = []byte{}
	var buffer = make([]byte, 1024 * 1024 * 4)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, 0, nil, err
			}

		md5DataContainer := make([]byte, 16)
		md5Data := md5.Sum(buffer[:n])
		copy(md5DataContainer, md5Data[:])
		block := hex.EncodeToString(md5DataContainer)
		blockList = append(blockList, block)
		fileData = append(fileData, buffer[:n]...)
	}

	return blockList, size, fileData, nil		
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
	return "121.fd4aaafe9c485ee574ede6d425b7dc44.Y3yDmZzKQMW56XOwD8Sg3yXuEci7UV5_aFOPY--.mq6tWg", nil
}



func GetBaiduNetworkdiskToken() (*vo.BaiduTokenResponse, error) {
	clientId := os.Getenv("BAIDU_NETWORKDISK_CLIENT_ID")
	code := "30cbd944cc1706524de096e358d7c138"
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

	
  err = json.Unmarshal(body, &baiduTokenResponse)
	if err != nil {
		return nil, err
	}

	return &baiduTokenResponse, nil
}
