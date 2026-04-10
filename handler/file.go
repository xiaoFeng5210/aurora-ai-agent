package handler

import (
	"aurora-agent/handler/vo"
	"aurora-agent/service"
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
