package server

import (
	"database/sql"
	"log"

	"github.com/spf13/viper"
)

func InitDatabase(config *viper.Viper) *sql.DB {
	connectiongString := config.GetString("database.connection_string")
	maxIdleConnections := config.GetInt("database.max_idle_connections")
	maxOpenConnections := config.GetInt("database.max_open_connections")
	connectionMaxLifetime := config.GetDuration("database.connection_max_lifetime")
	driverName := config.GetString("database.driver_name")

	if connectiongString == "" {
		log.Fatal("Database connection string is missing")
	}

	dbHandler, initErr := sql.Open(driverName, connectiongString)

	if initErr != nil {
		log.Fatalf("Error while initializing database: %v", initErr)
	}
	dbHandler.SetMaxIdleConns(maxIdleConnections)
	dbHandler.SetMaxOpenConns(maxOpenConnections)
	dbHandler.SetConnMaxLifetime(connectionMaxLifetime)

	validationErr := dbHandler.Ping()
	if validationErr != nil {
		dbHandler.Close()
		log.Fatalf("Error while validating database: %v", validationErr)
	}

	return dbHandler
}
