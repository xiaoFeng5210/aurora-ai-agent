package router

import (
	"aurora-agent/handler"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"aurora-agent/middleware"

	"github.com/gin-gonic/gin"
)

const webDistDir = "web/dist"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/register", handler.CreateUser)

		apiv1.POST("/login", handler.Login)
		apiv1.POST("/chat/glm/stream/:document_id", middleware.Auth, handler.StreamChatWithGLMController)

		apiv1.GET("/test_jwt", middleware.Auth, func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test jwt success"})
		})
	}

	userGroup := apiv1.Group("/users")
	userGroup.Use(middleware.Auth)
	userGroup.GET("", handler.GetAllUsers)
	userGroup.GET("/:id", handler.GetUserById)
	userGroup.POST("/query", handler.QueryUser)
	userGroup.GET("/me", handler.GetCurrentUser)
	userGroup.PUT("/me", handler.UpdateCurrentUser)
	userGroup.PUT("/me/password", handler.ChangeCurrentUserPassword)
	userGroup.DELETE("/me", handler.DeleteCurrentUser)

	documentGroup := apiv1.Group("/documents")
	documentGroup.Use(middleware.Auth)
	documentGroup.GET("", handler.GetAllDocuments)
	documentGroup.POST("", handler.CreateDocument)
	documentGroup.GET("/:id", handler.GetDocumentById)
	documentGroup.POST("/query", handler.QueryDocument)
	documentGroup.PUT("/:id", handler.UpdateDocument)
	documentGroup.DELETE("/:id", handler.DeleteDocument)
	documentGroup.POST("/:id/messages", handler.CreateHistoryMessage)
	documentGroup.GET("/:id/messages/:message_id", handler.GetHistoryMessageByMessageID)
	documentGroup.POST("/:id/messages/query", handler.QueryHistoryMessages)
	documentGroup.GET("/:id/messages/proxy/history", handler.ProxyQueryHistoryMessages)
	documentGroup.PUT("/:id/messages/:message_id", handler.UpdateHistoryMessage)
	documentGroup.DELETE("/:id/messages/:message_id", handler.DeleteHistoryMessage)

	baiduNetworkdiskGroup := apiv1.Group("/file")
	baiduNetworkdiskGroup.GET("/baidu_networkdisk/token", handler.GetBaiduNetworkdiskToken)
	baiduNetworkdiskGroup.GET("/baidu_networkdisk/capacity", handler.GetBaiduNetworkdiskCapacity)
	baiduNetworkdiskGroup.GET("/baidu_networkdisk/file_list", handler.GetBaiduNetworkdiskFileList)
	baiduNetworkdiskGroup.POST("/baidu_networkdisk/upload", handler.GMBaiduNetworkdiskUpload)

	qdrantGroup := apiv1.Group("/qdrant")
	qdrantGroup.Use(middleware.Auth)
	qdrantGroup.POST("/upsert", handler.UpsertQdrantTest)
	qdrantGroup.POST("/query", handler.QueryQdrantVector)





	// RAG 个人知识库：上传文件 → 文本提取 → 向量化 → 写入 Qdrant
	// 每个用户的数据通过 user_id 隔离，互不可见
	ragGroup := apiv1.Group("/rag")
	ragGroup.Use(middleware.Auth)
	ragGroup.POST("/create", handler.CreateRag)

	registerWebStatic(r)

	return r
}

// registerWebStatic 将 web/dist 静态资源挂载到根路径，并对未命中的非 API
// 请求回退到 index.html，以支持 SPA 的前端路由。
func registerWebStatic(r *gin.Engine) {
	indexFile := filepath.Join(webDistDir, "index.html")
	assetsDir := filepath.Join(webDistDir, "assets")

	r.Static("/assets", assetsDir)
	r.StaticFile("/", indexFile)
	r.StaticFile("/index.html", indexFile)
	r.StaticFile("/favicon.svg", filepath.Join(webDistDir, "favicon.svg"))
	r.StaticFile("/icons.svg", filepath.Join(webDistDir, "icons.svg"))

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		fsPath := filepath.Join(webDistDir, filepath.Clean(path))
		if info, err := os.Stat(fsPath); err == nil && !info.IsDir() {
			c.File(fsPath)
			return
		}

		c.File(indexFile)
	})
}
