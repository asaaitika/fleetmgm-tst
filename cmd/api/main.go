package main

import (
	"fmt"
	"log"

	"github.com/asaaitika/fleetmgm-tst/internal/config"
	"github.com/asaaitika/fleetmgm-tst/internal/handlers"
	"github.com/asaaitika/fleetmgm-tst/internal/repositories"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("✅ Connected to PostgreSQL")

	vehicleRepo := repositories.NewVehicleRepository(db)

	vehicleHandler := handlers.NewVehicleHandler(vehicleRepo)

	router := gin.Default()

	router.GET("/vehicles/:vehicle_id/location", vehicleHandler.GetLastLocation)
	router.GET("/vehicles/:vehicle_id/history", vehicleHandler.GetLocationHistory)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 API Server starting on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
