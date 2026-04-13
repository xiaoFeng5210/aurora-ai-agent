package main

import (
	"aurora-agent/database"
	qdrant_db "aurora-agent/database/qdrant"
	"aurora-agent/router"
	utils "aurora-agent/utils"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	godotenv.Load()
	database.DBConnect()
	qdrant_db.QdrantConnect()
}

func main() {
	r := router.SetupRouter()
	utils.Logger.Info("Server started on port 1119")
	err := r.Run(":1119")
	if err != nil {
		utils.Logger.Error("Failed to start server", zap.Error(err))
	}
}
