package main

import (
	"aurora-agent/router"
	"aurora-agent/utils"

	"go.uber.org/zap"
)

func main() {
	logger := utils.InitZap("./log/zap")
	r := router.SetupRouter()
	logger.Info("Server started on port 1119")
	err := r.Run(":1119")
	if err != nil {
		logger.Error("Failed to start server", zap.Error(err))
	}
}
