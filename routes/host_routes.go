package routes

import (
	"github.com/dannonb/go-network-monitor/controllers"
	"github.com/gin-gonic/gin"
)

func HostRoute(router *gin.Engine) {
	router.POST("/host/:userId", controllers.AddHost())
	router.GET("/hosts", controllers.GetAllHosts())
}