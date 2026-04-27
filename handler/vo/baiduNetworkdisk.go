package vo

type BaiduTokenResponse struct {
	ExpiresIn        int    `json:"expires_in,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// 单位 T
type BaiduNetworkdiskCapacityResponse struct {
	Total  float64 `json:"total"`
	Expire bool    `json:"expire"`
	Used   float64 `json:"used"`
	Free   float64 `json:"free"`
}

type FileCategory int

const (
	fileVideo    FileCategory = 1
	fileAudio    FileCategory = 2
	fileImage    FileCategory = 3
	fileDocument FileCategory = 4
	fileApp      FileCategory = 5
)

type Info struct {
	Category       int    `json:"category"`
	FsId           int    `json:"fs_id"`
	Path           string `json:"path"`            // 绝对路径
	ServerFilename string `json:"server_filename"` // 文件名称
	Isdir          int    `json:"isdir"`
	LocalCtime     int    `json:"local_ctime"`
	LocalMtime     int    `json:"local_mtime"`
	ServerCtime    int    `json:"server_ctime"`
	ServerMtime    int    `json:"server_mtime"`
	Size           int    `json:"size"`
}

type BaiduNetworkdiskFileListResponse struct {
	List  []Info `json:"list"`
	Errno int    `json:"errno"`
}

type PrecreateUploadResponse struct {
	Errno     int    `json:"errno"`
	Path      string `json:"path"`
	Uploadid  string `json:"uploadid"`
	Partseq   int    `json:"partseq"`
	Size      int    `json:"size"`
	BlockList []int  `json:"block_list"`
}

type UploadResponse struct {
	Errno int    `json:"errno"`
	MD5   string `json:"md5"`
}

type CreateFileOrDirResponse struct {
	Errno int `json:"errno"`
	FsId  int `json:"fs_id"`
}

type BaiduNetworkdiskFileManagerInfo struct {
	Errno  int    `json:"errno,omitempty"`
	Errmsg string `json:"errmsg,omitempty"`
	Path   string `json:"path,omitempty"`
	ToPath string `json:"to_path,omitempty"`
}

type BaiduNetworkdiskDeleteFileResponse struct {
	Errno            int                               `json:"errno"`
	Errmsg           string                            `json:"errmsg,omitempty"`
	Error            string                            `json:"error,omitempty"`
	ErrorDescription string                            `json:"error_description,omitempty"`
	Info             []BaiduNetworkdiskFileManagerInfo `json:"info,omitempty"`
	RequestId        int64                             `json:"request_id,omitempty"`
	Taskid           int64                             `json:"taskid,omitempty"`
}
