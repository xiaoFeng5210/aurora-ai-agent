package handler

import (
	"aurora-agent/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UID_IN_TOKEN = "uid"
	UID_IN_CTX   = "uid"
	COOKIE_NAME  = "jwt"
)

var (
	KeyConfig = utils.InitViper("conf", "jwt", "yaml")
)

// 从cookie里取出jwt，从而取出uid
func GetLoginUid(ctx *gin.Context) int {
	//依靠浏览器自动回传的cookie，提取出jwt token
	token := ""
	for _, cookie := range ctx.Request.Cookies() {
		if cookie.Name == COOKIE_NAME {
			token = cookie.Value
		}
	}
	return GetUidFromJwt(token)
}

// 从jwt里取出uid
func GetUidFromJwt(token string) int {
	_, payload, err := utils.VerifyJwt(token, KeyConfig.GetString("secret"))
	if err != nil {
		return 0
	}
	for k, v := range payload.UserDefined {
		if k == UID_IN_TOKEN {
			return int(v.(float64))
		}
	}
	return 0
}

// 身份认证中间件，先确保是登录状态
func Auth(ctx *gin.Context) {
	loginUid := GetLoginUid(ctx)
	if loginUid <= 0 {
		// ctx.String(http.StatusForbidden, "auth failed") //返回403
		ctx.Redirect(http.StatusTemporaryRedirect, "/login") //重定向到登录页面
		ctx.Abort()                                          //中断。通过Abort()使中间件后面的handler不再执行，但是本handler还是会继续执行。所以下一行代码需要显式return
	} else {
		ctx.Set(UID_IN_CTX, loginUid) //把登录的uid放入ctx中
	}
}
