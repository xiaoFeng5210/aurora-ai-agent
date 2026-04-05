package vo


type BaiduTokenResponse struct {
	ExpiresIn int `json:"expires_in"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// 单位 T
type BaiduNetworkdiskCapacityResponse struct {
	Total float64 `json:"total"`
	Expire bool `json:"expire"`
	Used float64 `json:"used"`
	Free float64 `json:"free"`
}

type FileCategory int
const (
	fileVideo FileCategory = 1
  fileAudio FileCategory = 2
	fileImage FileCategory = 3
	fileDocument FileCategory = 4
	fileApp FileCategory = 5
)

type Info struct {
	Category int `json:"category"`
	FsId int `json:"fs_id"`
	Path string `json:"path"`   // 绝对路径
	ServerFilename string `json:"server_filename"` // 文件名称
  Isdir int `json:"isdir"`
	LocalCtime int `json:"local_ctime"`
	LocalMtime int `json:"local_mtime"`
	ServerCtime int `json:"server_ctime"`
	ServerMtime int `json:"server_mtime"`
	Size int `json:"size"`
}

type BaiduNetworkdiskFileListResponse struct {
	List []Info `json:"list"`
	Errno int `json:"errno"`
}


type PrecreateUploadResponse struct {
	Errno int `json:"errno"`
	Path string `json:"path"`
	
}


