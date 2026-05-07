package main

import (
	"aurora-agent/database"
	qdrant_db "aurora-agent/database/qdrant"
	rabbitmq_module "aurora-agent/database/rabbitmq"
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
	qdrant_db.QdrantConnect()
	err := qdrant_db.CreateRagCollection()
	if err != nil {
		panic(err)
	}
	if _, err = rabbitmq_module.Connect(); err != nil {
		panic(err)
	}
}

func main() {
	r := router.SetupRouter()
	utils.Logger.Info("Server started on port 1119")
	err := r.Run(":1119")
	if err != nil {
		utils.Logger.Error("Failed to start server", zap.Error(err))
	}
}
