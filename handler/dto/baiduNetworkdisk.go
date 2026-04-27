package dto

type GetBaiduNetworkdiskFileListRequest struct {
	Method      string `json:"method" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	Dir         string `json:"dir" binding:"required"`
}

type PrecreateUploadRequest struct {
	Path      string   `json:"path"`
	Size      int      `json:"size"`
	Isdir     int      `json:"isdir"`
	BlockList []string `json:"block_list"`
	Autoinit  int      `json:"autoinit"`
	Rtype     int      `json:"rtype"`
}

type CreateFileOrDirRequest struct {
	Path      string `json:"path"`
	Isdir     int    `json:"isdir"`
	Uploadid  string `json:"uploadid"`
	BlockList string `json:"block_list"`
	Size      int    `json:"size"`
}

type GMBaiduNetworkdiskUploadRequest struct {
	Path  string `json:"path"`
	Isdir int    `json:"isdir"`
}

type DeleteBaiduNetworkdiskFileRequest struct {
	Path  string   `json:"path"`
	Paths []string `json:"paths"`
	Async int      `json:"async"`
}

type ExchangeBaiduNetworkdiskTokenRequest struct {
	Code string `json:"code" binding:"required"`
}
