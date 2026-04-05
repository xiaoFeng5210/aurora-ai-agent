package dto

type GetBaiduNetworkdiskFileListRequest struct {
	Method string `json:"method" binding:"required"`
  AccessToken string `json:"access_token" binding:"required"`
	ParentPath string `json:"parent_path" binding:"required"`
}
