package main

import (
	"aurora-agent/database"
	"aurora-agent/router"
	"aurora-agent/utils"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger = utils.InitZap("./log/zap")
	database.DBConnect()
}

func main() {

	r := router.SetupRouter()
	logger.Info("Server started on port 1119")
	err := r.Run(":1119")
	if err != nil {
		logger.Error("Failed to start server", zap.Error(err))
	}
}
