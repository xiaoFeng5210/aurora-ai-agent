package utils

import "strconv"

// 百度网盘相关错误码
func CommonErrorCodeBaiduNetworkdisk(code int) string {
	switch code {
	case 2:
		return "百度网盘参数错误"
	case -6:
		return "百度网盘token过期"
	case 10:
		return "文件已存在"
	case 11:
		return "用户不存在"
	case 12:
		return "批量转存出错"
	default:
		return "未知错误" + strconv.Itoa(code)
	}
}

// precreate upload 错误码
func PrecreateUploadErrorCodeBaiduNetworkdisk(code int) string {
	switch code {
	case -7:
		return "文件或目录名错误或无权访问"
	case -10:
		return "容量不足"
	default:
		return "未知错误" + strconv.Itoa(code)
	}
}

// split upload 错误码
func SplitUploadErrorCodeBaiduNetworkdisk(code int) string {
	switch code {
	case 31024:
		return "没有申请上传权限"
	case 31299:
		return "第一个分片的大小小于4MB"
	case 31364:
		return "超出分片大小限制"
	case 31363:
		return "分片缺失"
	default:
		return "未知错误" + strconv.Itoa(code)
	}
}

// create file or dir 错误码
func CreateFileOrDirErrorCodeBaiduNetworkdisk(code int) string {
	switch code {
	case -7:
		return "文件或目录名错误或无权访问"
	case -8:
		return "文件或目录已存在"
	case 10:
		return "创建文件失败"
	default:
		return "未知错误" + strconv.Itoa(code)
	}
}

func DeleteFileErrorCodeBaiduNetworkdisk(code int) string {
	switch code {
	case -3, -9, 31066:
		return "文件或目录不存在"
	case -6, 110, 31045:
		return "百度网盘token无效或已过期"
	case -7:
		return "文件或目录名错误或无权访问"
	case 111:
		return "有其他异步任务正在执行"
	default:
		return CommonErrorCodeBaiduNetworkdisk(code)
	}
}
