package controllers

import (
	"monitoring-backend/models"
	"monitoring-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingController struct {
	service *services.PingService
}

func NewPingController(service *services.PingService) *PingController {
	return &PingController{service: service}
}

func (c *PingController) GetPings(ctx *gin.Context) {
	pings, err := c.service.GetAllPings()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pings)
}

func (c *PingController) CreatePing(ctx *gin.Context) {
	var ping models.Ping

	if err := ctx.ShouldBindJSON(&ping); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreatePing(ping); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Пинг создан"})
}
