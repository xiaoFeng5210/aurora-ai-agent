package main

import (
	"aurora-agent/database"
	"aurora-agent/router"
	utils "aurora-agent/utils"

	"go.uber.org/zap"
)


func init() {
	database.DBConnect()
}

func main() {
	r := router.SetupRouter()
	utils.Logger.Info("Server started on port 1119")
	err := r.Run(":1119")
	if err != nil {
		utils.Logger.Error("Failed to start server", zap.Error(err))
	}
}
