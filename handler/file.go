package handler

import (
	"aurora-agent/handler/dto"
	"aurora-agent/handler/vo"
	"aurora-agent/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GMBaiduNetworkdiskUpload(ctx *gin.Context) {
	var request dto.GMBaiduNetworkdiskUploadRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	file, err := ctx.FormFile("file")

	err = service.GMBaiduNetworkdiskUpload(request.Path, request.Isdir)
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
