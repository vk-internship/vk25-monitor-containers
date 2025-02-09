package main

import (
	"log"
	"monitoring-backend/controllers"
	"monitoring-backend/database"
	"monitoring-backend/repositories"
	"monitoring-backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := database.Connect()
	defer db.Close()

	if err := database.InitDB(db); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	pingRepo := repositories.NewPingRepository(db)
	pingService := services.NewPingService(pingRepo)
	pingController := controllers.NewPingController(pingService)

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/pings", pingController.GetPings)
	router.POST("/pings", pingController.CreateOrUpdatePing)

	router.Run(":8080")
}
