package main

import (
	"config"
	"log"
	"marketplace/server"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting marketplace app...")
	log.Println("Initializing configuration")

	config := config.InitConfig("appConfig")
	log.Println("Initializing DB")
	dbHandler := server.InitDatabase(config)

	log.Println("Initializing HTTP server")
	httpServer := server.InitHttpServer(config, dbHandler)

	httpServer.Start()
}
