package dto

type GetBaiduNetworkdiskFileListRequest struct {
	Method string `json:"method" binding:"required"`
  AccessToken string `json:"access_token" binding:"required"`
	Dir string `json:"dir" binding:"required"`
}


type PrecreateUploadRequest struct {
	Method string `json:"method"`
	AccessToken string `json:"access_token"`
	Path string `json:"path"`
	Size int `json:"size"`
	Isdir int `json:"isdir"`
	BlockList []string `json:"block_list"`
	Autoinit int `json:"autoinit"`
}
