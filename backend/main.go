package main

import (
	"monitoring-backend/controllers"
	"monitoring-backend/database"
	"monitoring-backend/repositories"
	"monitoring-backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.Connect()
	defer db.Close()

	pingRepo := repositories.NewPingRepository(db)
	pingService := services.NewPingService(pingRepo)
	pingController := controllers.NewPingController(pingService)

	router := gin.Default()
	router.GET("/pings", pingController.GetPings)
	router.POST("/pings", pingController.CreatePing)

	router.Run(":8080")
}
