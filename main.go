package main

import (
	"aurora-agent/bootstrap"
	"aurora-agent/database"
	"aurora-agent/database/aliyunossvector"
	redis_db "aurora-agent/database/redis"
	"aurora-agent/router"
	utils "aurora-agent/utils"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	godotenv.Load()
	database.DBConnect()
	if _, err := redis_db.RedisConnect(); err != nil {
		panic(err)
	}
	if err := aliyunossvector.Connect(); err != nil {
		panic(err)
	}
	if err := aliyunossvector.EnsureRagResources(); err != nil {
		panic(err)
	}
}

func main() {
	bootstrap.StartConsumer()
	r := router.SetupRouter()
	utils.Logger.Info("Server started on port 1119")
	err := r.Run(":1119")
	if err != nil {
		utils.Logger.Error("Failed to start server", zap.Error(err))
	}
}
