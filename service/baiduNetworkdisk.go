package service

import (
	"aurora-agent/database/redis"
	"aurora-agent/handler/vo"
	"aurora-agent/utils"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"crypto/md5"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	baseUrl = "https://pan.baidu.com"
	headers = map[string]string{
		"User-Agent": "pan.baidu.com",
	}
	capacityUnit float64 = 1024 * 1024 * 1024
)

const baiduNetworkdiskAppRoot = "/apps/aurora-ai-agent知识库"
const baiduNetworkdiskTokenKey = "baidu_networkdisk:token"
const baiduNetworkdiskTokenTTL = 21 * 24 * time.Hour

type baiduNetworkdiskFileManagerItem struct {
	Path string `json:"path"`
}

type baiduNetworkdiskFileMetaResponse struct {
	Errno  int                        `json:"errno"`
	Errmsg string                     `json:"errmsg,omitempty"`
	List   []baiduNetworkdiskFileMeta `json:"list"`
}

type baiduNetworkdiskFileMeta struct {
	Category int    `json:"category"`
	Dlink    string `json:"dlink"`
	FsId     int64  `json:"fs_id"`
	Filename string `json:"filename"`
	Isdir    int    `json:"isdir"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
}

var ErrBaiduNetworkdiskInvalidDeleteRequest = errors.New("invalid baidu networkdisk delete request")

var logger *zap.Logger

func init() {
	logger = utils.Logger
}

// 删除文件或文件夹。path 可以是列表接口返回的绝对路径，也可以是应用目录下的相对路径。
func DeleteBaiduNetworkdiskFiles(paths []string, async int) (*vo.BaiduNetworkdiskDeleteFileResponse, error) {
	if async < 0 || async > 2 {
		return nil, fmt.Errorf("%w: async must be 0, 1, or 2", ErrBaiduNetworkdiskInvalidDeleteRequest)
	}

	normalizedPaths, err := normalizeBaiduNetworkdiskDeletePaths(paths)
	if err != nil {
		return nil, err
	}

	access_token, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {
		return nil, err
	}

	filelist := make([]baiduNetworkdiskFileManagerItem, 0, len(normalizedPaths))
	for _, path := range normalizedPaths {
		filelist = append(filelist, baiduNetworkdiskFileManagerItem{Path: path})
	}

	filelistBytes, err := json.Marshal(filelist)
	if err != nil {
		return nil, err
	}

	httpUrl := baseUrl + "/rest/2.0/xpan/file" + fmt.Sprintf("?access_token=%s&method=filemanager&opera=delete", access_token)
	formData := url.Values{
		"async":    {strconv.Itoa(async)},
		"filelist": {string(filelistBytes)},
	}

	req, err := http.NewRequest(http.MethodPost, httpUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", headers["User-Agent"])
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result vo.BaiduNetworkdiskDeleteFileResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.Errno != 0 {
		return nil, errors.New(utils.DeleteFileErrorCodeBaiduNetworkdisk(result.Errno))
	}

	return &result, nil
}

func normalizeBaiduNetworkdiskDeletePaths(paths []string) ([]string, error) {
	normalized := make([]string, 0, len(paths))
	seen := make(map[string]struct{}, len(paths))

	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}

		if strings.HasPrefix(path, "/") {
			if path == baiduNetworkdiskAppRoot {
				return nil, fmt.Errorf("%w: 不能删除百度网盘应用根目录", ErrBaiduNetworkdiskInvalidDeleteRequest)
			}
			if !strings.HasPrefix(path, baiduNetworkdiskAppRoot+"/") {
				return nil, fmt.Errorf("%w: 只能删除%s目录下的文件或文件夹", ErrBaiduNetworkdiskInvalidDeleteRequest, baiduNetworkdiskAppRoot)
			}
		} else {
			path = baiduNetworkdiskAppRoot + "/" + strings.TrimLeft(path, "/")
		}

		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		normalized = append(normalized, path)
	}

	if len(normalized) == 0 {
		return nil, fmt.Errorf("%w: path or paths is required", ErrBaiduNetworkdiskInvalidDeleteRequest)
	}

	return normalized, nil
}

func BaiduNetworkdiskFilenamesFromPaths(paths []string) []string {
	filenames := make([]string, 0, len(paths))
	seen := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		filename := filepath.Base(strings.TrimRight(strings.TrimSpace(path), "/"))
		if filename == "." || filename == "/" || filename == "" {
			continue
		}
		if _, ok := seen[filename]; ok {
			continue
		}
		seen[filename] = struct{}{}
		filenames = append(filenames, filename)
	}
	return filenames
}

// 通用上传文件或文件夹
func GMBaiduNetworkdiskUpload(fileParam multipart.File, paramPath string, isdir string, username string) (*vo.CreateFileOrDirResponse, error) {
	isdirInt, err := strconv.Atoi(isdir)
	if err != nil {
		return nil, errors.New("isdir is not a valid integer: " + err.Error())
	}
	precreateInfo, fileData, blockList, err := PrecreateUpload(paramPath, isdirInt, fileParam, username)
	if err != nil {
		logger.Error("PrecreateUpload failed", zap.Error(err))
		return nil, err
	}

	const chunkSize = 1024 * 1024 * 4
	for _, blockIdx := range precreateInfo.BlockList {
		start := blockIdx * chunkSize
		end := start + chunkSize
		if end > len(fileData) {
			end = len(fileData)
		}
		partData := fileData[start:end]
		_, err = SplitUpload(precreateInfo, partData, blockIdx)
		if err != nil {
			logger.Error("Upload failed", zap.Error(err))
			return nil, err
		}
	}

	result, err := CreateFileOrDir(precreateInfo.Path, 0, precreateInfo.Uploadid, blockList, precreateInfo.Size)
	if err != nil {
		logger.Error("CreateFileOrDir failed", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func CreateFileOrDir(path string, isdir int, uploadid string, blockList string, size int) (*vo.CreateFileOrDirResponse, error) {
	access_token, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {
		return nil, err
	}
	method := "create"

	if isdir == 1 {
		size = 0
	}
	httpUrl := baseUrl + "/rest/2.0/xpan/file" + fmt.Sprintf("?access_token=%s&method=%s", access_token, method)

	formData := url.Values{
		"path":       {path},
		"size":       {strconv.Itoa(size)},
		"isdir":      {strconv.Itoa(isdir)},
		"rtype":      {"3"},
		"uploadid":   {uploadid},
		"block_list": {blockList},
	}
	encodedFormData := formData.Encode()

	req, err := http.NewRequest(http.MethodPost, httpUrl, strings.NewReader(encodedFormData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", headers["User-Agent"])
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result vo.CreateFileOrDirResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.Errno != 0 {
		return nil, errors.New(utils.CreateFileOrDirErrorCodeBaiduNetworkdisk(result.Errno))
	}
	return &result, nil
}

// 分片上传
func SplitUpload(precreateInfo *vo.PrecreateUploadResponse, fileData []byte, partseqIdx int) (*vo.UploadResponse, error) {

	access_token, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {
		return nil, err
	}
	method := "upload"
	path := precreateInfo.Path
	uploadid := precreateInfo.Uploadid
	partseq := partseqIdx
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
		return nil, errors.New(utils.SplitUploadErrorCodeBaiduNetworkdisk(result.Errno))
	}
	return &result, nil
}

// 预上传
func PrecreateUpload(paramPath string, isdir int, fileParam multipart.File, username string) (*vo.PrecreateUploadResponse, []byte, string, error) {
	access_token, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {
		return nil, nil, "", err
	}
	httpUrl := baseUrl + "/rest/2.0/xpan/file?" + fmt.Sprintf("method=precreate&access_token=%s", access_token)
	appName := "/aurora-ai-agent知识库"
	var parent string
	if username == "zhangqingfeng" {
		parent = "/apps" + appName + "/super" + "/"
	} else {
		parent = "/apps" + appName + "/users-data/" + username + "/"
	}
	path := parent + paramPath

	blockList, size, fileData, err := HandleFile(fileParam)
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
		return nil, nil, "", errors.New(utils.PrecreateUploadErrorCodeBaiduNetworkdisk(result.Errno))
	}

	result.Path = path
	result.Size = int(size)
	return &result, fileData, string(blockListBytes), nil
}

func HandleFile(fileParam multipart.File) (blockList []string, size int64, fileData []byte, err error) {
	fileData = []byte{}

	var buffer = make([]byte, 1024*1024*4)
	for {
		n, err := fileParam.Read(buffer)
		if n > 0 {
			md5DataContainer := make([]byte, 16)
			md5Data := md5.Sum(buffer[:n])
			copy(md5DataContainer, md5Data[:])
			block := hex.EncodeToString(md5DataContainer)
			blockList = append(blockList, block)
			fileData = append(fileData, buffer[:n]...)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, 0, nil, errors.New("读取上传文件失败: " + err.Error())
		}

	}

	size = int64(len(fileData))

	return blockList, size, fileData, nil
}

// 获取文件或文件夹列表
func GetBaiduNetworkdiskFileList(dir string) (*vo.BaiduNetworkdiskFileListResponse, error) {
	access_token, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {
		return nil, err
	}

	allDir := "/apps/aurora-ai-agent知识库/" + dir

	method := "list"
	url := baseUrl + "/rest/2.0/xpan/file?" + fmt.Sprintf("access_token=%s&method=%s&dir=%s", access_token, method, allDir)
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
	access_token, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {
		return nil, err
	}
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
	client := redis.Client()
	if client == nil {
		return "", errors.New("redis client is not initialized")
	}

	value, err := client.Get(context.Background(), baiduNetworkdiskTokenKey).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", errors.New("没有百度网盘token了")
		}
		return "", err
	}

	var token vo.BaiduTokenResponse
	if err = json.Unmarshal([]byte(value), &token); err == nil && token.AccessToken != "" {
		return token.AccessToken, nil
	}

	if strings.TrimSpace(value) == "" {
		return "", errors.New("没有百度网盘token了")
	}

	return value, nil
}

func GetBaiduNetworkdiskToken(code string) (*vo.BaiduTokenResponse, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, errors.New("code is required")
	}

	clientId := os.Getenv("BAIDU_NETWORKDISK_CLIENT_ID")
	clientSecret := os.Getenv("BAIDU_NETWORKDISK_CLIENT_SECRET")
	if clientId == "" || clientSecret == "" {
		return nil, errors.New("baidu networkdisk client id or client secret is not configured")
	}

	query := url.Values{}
	query.Set("grant_type", "authorization_code")
	query.Set("code", code)
	query.Set("client_id", clientId)
	query.Set("client_secret", clientSecret)
	query.Set("redirect_uri", "oob")

	tokenURL := "https://openapi.baidu.com/oauth/2.0/token?" + query.Encode()
	resp, err := http.Get(tokenURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result vo.BaiduTokenResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.Error != "" {
		return nil, fmt.Errorf("baidu networkdisk token exchange failed: %s", result.ErrorDescription)
	}
	if result.AccessToken == "" {
		return nil, errors.New("baidu networkdisk token response missing access_token")
	}

	if err = saveBaiduNetworkdiskTokenToRedis(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func saveBaiduNetworkdiskTokenToRedis(token *vo.BaiduTokenResponse) error {
	client := redis.Client()
	if client == nil {
		return errors.New("redis client is not initialized")
	}

	value, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return client.Set(context.Background(), baiduNetworkdiskTokenKey, string(value), baiduNetworkdiskTokenTTL).Err()
}

// 下载文件到本地临时文件夹下
func DownloadFile2TempFolder(fileCloudPath string) (string, error) {
	fileData, err := DonwloadFileFromBaiduNetworkdisk(fileCloudPath)
	if err != nil {
		return "", err
	}
	fileCloudPathSlice := strings.Split(fileCloudPath, "/")
	tempFile, err := os.CreateTemp("", fileCloudPathSlice[len(fileCloudPathSlice)-1])
	if err != nil {
		return "", err
	}

	defer tempFile.Close()

	_, err = tempFile.Write(fileData)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func DonwloadFileFromBaiduNetworkdisk(fileCloudPath string) ([]byte, error) {
	cloudPath, err := normalizeBaiduNetworkdiskFilePath(fileCloudPath)
	if err != nil {
		return nil, err
	}

	accessToken, err := GetBaiduNetworkdiskTokenFromRedis()
	if err != nil {
		return nil, err
	}

	fileInfo, err := getBaiduNetworkdiskFileInfoByPath(accessToken, cloudPath)
	if err != nil {
		return nil, err
	}
	if fileInfo.Isdir == 1 {
		return nil, fmt.Errorf("baidu networkdisk path is a directory: %s", cloudPath)
	}

	fileMeta, err := getBaiduNetworkdiskFileMeta(accessToken, fileInfo.FsId)
	if err != nil {
		return nil, err
	}
	if fileMeta.Dlink == "" {
		return nil, fmt.Errorf("baidu networkdisk file dlink is empty: %s", cloudPath)
	}

	downloadURL, err := appendBaiduNetworkdiskAccessToken(fileMeta.Dlink, accessToken)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", headers["User-Agent"])

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("baidu networkdisk download failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	return body, nil
}

func GetBaiduNetworkdiskFileData(fileCloudPath string) ([]byte, error) {
	return DonwloadFileFromBaiduNetworkdisk(fileCloudPath)
}

func normalizeBaiduNetworkdiskFilePath(fileCloudPath string) (string, error) {
	fileCloudPath = strings.TrimSpace(fileCloudPath)
	if fileCloudPath == "" {
		return "", errors.New("baidu networkdisk file path is required")
	}

	if strings.HasPrefix(fileCloudPath, "/") {
		cleanPath := path.Clean(fileCloudPath)
		if cleanPath == baiduNetworkdiskAppRoot {
			return "", errors.New("baidu networkdisk file path cannot be app root")
		}
		if !strings.HasPrefix(cleanPath, baiduNetworkdiskAppRoot+"/") {
			return "", fmt.Errorf("baidu networkdisk file path must be under %s", baiduNetworkdiskAppRoot)
		}
		return cleanPath, nil
	}

	cleanPath := path.Clean(baiduNetworkdiskAppRoot + "/" + strings.TrimLeft(fileCloudPath, "/"))
	if !strings.HasPrefix(cleanPath, baiduNetworkdiskAppRoot+"/") {
		return "", fmt.Errorf("baidu networkdisk file path must be under %s", baiduNetworkdiskAppRoot)
	}
	return cleanPath, nil
}

func getBaiduNetworkdiskFileInfoByPath(accessToken string, fileCloudPath string) (*vo.Info, error) {
	dir := path.Dir(fileCloudPath)
	fileList, err := getBaiduNetworkdiskFileListByAbsoluteDir(accessToken, dir)
	if err != nil {
		return nil, err
	}

	for idx := range fileList.List {
		item := fileList.List[idx]
		if item.Path == fileCloudPath {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("baidu networkdisk file not found: %s", fileCloudPath)
}

func getBaiduNetworkdiskFileListByAbsoluteDir(accessToken string, dir string) (*vo.BaiduNetworkdiskFileListResponse, error) {
	query := url.Values{}
	query.Set("access_token", accessToken)
	query.Set("method", "list")
	query.Set("dir", dir)

	req, err := http.NewRequest(http.MethodGet, baseUrl+"/rest/2.0/xpan/file?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", headers["User-Agent"])

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("baidu networkdisk file list failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var result vo.BaiduNetworkdiskFileListResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.Errno != 0 {
		return nil, errors.New(utils.CommonErrorCodeBaiduNetworkdisk(result.Errno))
	}

	return &result, nil
}

func getBaiduNetworkdiskFileMeta(accessToken string, fsID int) (*baiduNetworkdiskFileMeta, error) {
	query := url.Values{}
	query.Set("method", "filemetas")
	query.Set("access_token", accessToken)
	query.Set("fsids", fmt.Sprintf("[%d]", fsID))
	query.Set("dlink", "1")

	req, err := http.NewRequest(http.MethodGet, baseUrl+"/rest/2.0/xpan/multimedia?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", headers["User-Agent"])

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("baidu networkdisk file meta failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var result baiduNetworkdiskFileMetaResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.Errno != 0 {
		if result.Errmsg != "" {
			return nil, fmt.Errorf("baidu networkdisk file meta failed: errno=%d errmsg=%s", result.Errno, result.Errmsg)
		}
		return nil, errors.New(utils.CommonErrorCodeBaiduNetworkdisk(result.Errno))
	}
	if len(result.List) == 0 {
		return nil, fmt.Errorf("baidu networkdisk file meta not found: fs_id=%d", fsID)
	}

	return &result.List[0], nil
}

func appendBaiduNetworkdiskAccessToken(dlink string, accessToken string) (string, error) {
	downloadURL, err := url.Parse(dlink)
	if err != nil {
		return "", err
	}

	query := downloadURL.Query()
	query.Set("access_token", accessToken)
	downloadURL.RawQuery = query.Encode()

	return downloadURL.String(), nil
}
