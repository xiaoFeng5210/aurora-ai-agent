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
