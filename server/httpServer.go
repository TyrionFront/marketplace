package server

import (
	"database/sql"
	"log"
	"repositories"
	"services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"controllers"
)

type HttpServer struct {
	config          *viper.Viper
	router          *gin.Engine
	statsController *controllers.StatsController
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	statsRepository := repositories.NewStatsRepository(dbHandler)
	statsService := services.NewStatsService(statsRepository)
	statsController := controllers.NewStatsController(statsService)

	router := gin.Default()
	router.POST("/points", statsController.SaveStats)

	return HttpServer{
		config:          config,
		router:          router,
		statsController: statsController,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.GetString("http.server_address"))

	if err != nil {
		log.Printf("Error while starting HTTP server: %v", err)
	}
}
