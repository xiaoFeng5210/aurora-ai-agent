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

		apiv1.GET("/test_jwt", middleware.Auth, func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test jwt success"})
		})
	}

	userGroup := apiv1.Group("/users")
	{
		userGroup.GET("", middleware.Auth, handler.GetAllUsers)
		userGroup.GET("/me", middleware.Auth, handler.GetCurrentUser)
		userGroup.PUT("/me", middleware.Auth, handler.UpdateCurrentUser)
		userGroup.PUT("/me/password", middleware.Auth, handler.ChangeCurrentUserPassword)
		userGroup.DELETE("/me", middleware.Auth, handler.DeleteCurrentUser)
		userGroup.GET("/:id", middleware.Auth, handler.GetUserById)
		userGroup.POST("/query", middleware.Auth, handler.QueryUser)
	}

	return r
}
