package main

import (
	"config"
	"fmt"
	"log"
	"marketplace/server"
	"os"

	_ "github.com/lib/pq"
)

func getConfigFileName() string {
	env := os.Getenv("ENV")
	if env != "" {
		return "appConfig-" + env
	}
	return "appConfig"
}

func main() {
	log.Println("Starting marketplace app...")
	log.Println("Initializing configuration")

	configName := getConfigFileName()
	fmt.Println(configName)
	config := config.InitConfig(configName)
	log.Println("Initializing DB")
	dbHandler := server.InitDatabase(config)

	log.Println("Initializing HTTP server")
	httpServer := server.InitHttpServer(config, dbHandler)

	httpServer.Start()
}
