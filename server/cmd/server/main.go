package main

import (
	"log"
	"library/config"
	"library/repository/mysql"
	"library/database"
	"library/router"
	"library/service"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database
	if err := database.InitMySQL(); err != nil {
		log.Fatalf("Error initializing MySQL: %v", err)
	}

	// Create MySQL factory
	mysqlFactory := mysql.NewFactory(database.DB)

	// Create service factory
	factory := service.NewFactory(mysqlFactory)

	// Set up the router
	r := router.SetupRouter(factory)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
