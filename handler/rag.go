package handler

import (
	"aurora-agent/handler/vo"
	"aurora-agent/middleware"
	"aurora-agent/service"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 支持上传的文件类型白名单
// 新增格式时需同步更新 service/ragService.go 中的 fileParsers 注册表
var allowedFileExtensions = map[string]bool{
	".txt": true,
	".md":  true,
	".csv": true,
	// TODO: 待对应解析器实现后启用
	// ".pdf":  true,
	// ".docx": true,
}

// 上传文件并且将文件向量化上传到qdrant, 只是将任务发送给消息队列
func RagVectorizeTask(ctx *gin.Context) {
	fh, err := ctx.FormFile("file")
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedFileExtensions[ext] {
		vo.RespondError(ctx, http.StatusBadRequest,
			vo.ErrUnsupportedFileType)
		return
	}

	file, err := fh.Open()
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}
	defer file.Close()

	uid := ctx.GetInt(middleware.UID_IN_CTX)

	err = service.SendRagVectorizeTask(uid, file, fh.Filename)
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}

	vo.RespondSuccess(ctx, gin.H{
		"filename": fh.Filename,
		"status":   "pending",
	})
}

// CreateRag 上传文件并写入个人 RAG 知识库
// 流程：接收文件 → 校验类型 → 读取文本内容 → 分块+向量化 → 写入 Qdrant
// 每个用户的知识库通过 user_id 隔离，互不可见
func CreateRag(ctx *gin.Context) {
	// 1. 从 multipart/form-data 中获取上传文件
	fh, err := ctx.FormFile("file")
	if err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedFileExtensions[ext] {
		vo.RespondError(ctx, http.StatusBadRequest,
			vo.ErrUnsupportedFileType)
		return
	}

	// 3. 打开文件流并读取全部内容
	file, err := fh.Open()
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}
	defer file.Close()

	// 4. 调用 service 层完成：文本提取 → 分块 → 向量化 → 写入 Qdrant
	uid := ctx.GetInt(middleware.UID_IN_CTX)
	result, err := service.CreateRag(uid, file, fh.Filename)
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}

	vo.RespondSuccess(ctx, result)
}
