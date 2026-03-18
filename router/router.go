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

		apiv1.GET("/users", handler.GetAllUsers)

		apiv1.GET("/test_jwt", middleware.Auth, func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test jwt success"})
		})
	}

	return r
}
