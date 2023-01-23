package main

import (
	"fmt"
	"log"

	"github.com/sdpatel1026/doc-search/configs"
)

func main() {
	log.Println("doc-search service up and running...")
	initializeRoutes()
	appPort := configs.GetEnvWithKey("APP_PORT", "8080")
	serverAddress := fmt.Sprintf(":%s", appPort)
	router.Run(serverAddress)
}
