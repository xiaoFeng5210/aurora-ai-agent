package handler

import (
	"aurora-agent/handler/dto"
	"aurora-agent/handler/vo"
	"aurora-agent/middleware"
	"aurora-agent/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GMBaiduNetworkdiskUpload(ctx *gin.Context) {
	filename := ctx.PostForm("filename")
	isdir := ctx.PostForm("isdir")

	fh, err := ctx.FormFile("file")
	fileParam, err := fh.Open()

	err = service.GMBaiduNetworkdiskUpload(fileParam, filename, isdir)
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}
	vo.RespondSuccess(ctx, nil)
}

func GetBaiduNetworkdiskCapacity(ctx *gin.Context) {
	resp, err := service.GetBaiduNetworkdiskCapacity()
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}
	vo.RespondSuccess(ctx, resp)
}

func GetBaiduNetworkdiskToken(ctx *gin.Context) {
	resp, err := service.GetBaiduNetworkdiskToken()
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    *resp,
	})
}

func GetBaiduNetworkdiskFileList(ctx *gin.Context) {
	dir := ctx.Query("dir")
	resp, err := service.GetBaiduNetworkdiskFileList(dir)
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}
	vo.RespondSuccess(ctx, resp)
}

func DeleteBaiduNetworkdiskFile(ctx *gin.Context) {
	var req dto.DeleteBaiduNetworkdiskFileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	paths := make([]string, 0, len(req.Paths)+1)
	if req.Path != "" {
		paths = append(paths, req.Path)
	}
	paths = append(paths, req.Paths...)
	if len(paths) == 0 {
		vo.RespondError(ctx, http.StatusBadRequest, errors.New("path or paths is required"))
		return
	}
	if req.Async < 0 || req.Async > 2 {
		vo.RespondError(ctx, http.StatusBadRequest, errors.New("async must be 0, 1, or 2"))
		return
	}

	resp, err := service.DeleteBaiduNetworkdiskFiles(paths, req.Async)
	if err != nil {
		if errors.Is(err, service.ErrBaiduNetworkdiskInvalidDeleteRequest) {
			vo.RespondError(ctx, http.StatusBadRequest, err)
			return
		}
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}

	qdrantResult, err := service.DeleteQdrantVectorByFilenames(
		ctx.GetInt(middleware.UID_IN_CTX),
		service.BaiduNetworkdiskFilenamesFromPaths(paths),
	)
	if err != nil {
		vo.RespondError(ctx, http.StatusInternalServerError, err)
		return
	}

	vo.RespondSuccess(ctx, gin.H{
		"baidu_networkdisk": resp,
		"qdrant":            qdrantResult,
	})
}
