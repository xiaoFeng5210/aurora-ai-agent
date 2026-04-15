package router

import (
	"aurora-agent/handler"

	"aurora-agent/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/register", handler.CreateUser)

		apiv1.POST("/login", handler.Login)
		apiv1.POST("/chat/glm/stream/:document_id", handler.StreamChatWithGLMController)

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
	documentGroup.POST("/:document_id/messages", handler.CreateHistoryMessage)
	documentGroup.GET("/:document_id/messages/:message_id", handler.GetHistoryMessageByMessageID)
	documentGroup.POST("/:document_id/messages/query", handler.QueryHistoryMessages)
	documentGroup.GET("/:document_id/messages/proxy/history", handler.ProxyQueryHistoryMessages)
	documentGroup.PUT("/:document_id/messages/:message_id", handler.UpdateHistoryMessage)
	documentGroup.DELETE("/:document_id/messages/:message_id", handler.DeleteHistoryMessage)

	baiduNetworkdiskGroup := apiv1.Group("/file")
	baiduNetworkdiskGroup.GET("/baidu_networkdisk/token", handler.GetBaiduNetworkdiskToken)
	baiduNetworkdiskGroup.GET("/baidu_networkdisk/capacity", handler.GetBaiduNetworkdiskCapacity)
	baiduNetworkdiskGroup.GET("/baidu_networkdisk/file_list", handler.GetBaiduNetworkdiskFileList)
	baiduNetworkdiskGroup.POST("/baidu_networkdisk/upload", handler.GMBaiduNetworkdiskUpload)

	qdrantGroup := apiv1.Group("/qdrant")
	qdrantGroup.POST("/upsert", handler.UpsertQdrantTest)
	qdrantGroup.POST("/query", handler.QueryQdrantVector)
	return r
}
