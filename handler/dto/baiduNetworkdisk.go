package dto

type GetBaiduNetworkdiskFileListRequest struct {
	Method string `json:"method" binding:"required"`
  AccessToken string `json:"access_token" binding:"required"`
	Dir string `json:"dir" binding:"required"`
}


type PrecreateUploadRequest struct {
	Method string `json:"method" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	Path string `json:"path" binding:"required"`
	Size int `json:"size" binding:"required"`
	Isdir int `json:"isdir" binding:"required"`
	BlockList []string `json:"block_list" binding:"required"`
	Autoinit int `json:"autoinit" binding:"required"`
}
