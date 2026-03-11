package main

import (
	"aurora-agent/router"
	"log"
)

func main() {
	r := router.SetupRouter()
	err := r.Run(":1119")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
