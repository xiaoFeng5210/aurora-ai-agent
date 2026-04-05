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

type Info struct {
	Category int `json:"category"`
	FsId int `json:"fs_id"`
  Isdir bool `json:"isdir"`
	LocalCtime int `json:"local_ctime"`
	LocalMtime int `json:"local_mtime"`
	ServerCtime int `json:"server_ctime"`
	ServerMtime int `json:"server_mtime"`
	Size int `json:"size"`
}

type BaiduNetworkdiskFileListResponse struct {
	Info []Info `json:"info"`
}
