package routes

import (
	"github.com/dannonb/go-network-monitor/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine)  {
	router.GET("/user/:userId", controllers.GetUser()) 
	router.PUT("/user/:userId", controllers.EditUser())
	router.DELETE("/user/:userId", controllers.DeleteUser())
	router.GET("/users", controllers.GetAllUsers())
}

func UserAuthRoute(router *gin.Engine) {
	router.POST("/user/sign-up", controllers.SignUp())
	router.POST("/user/login", controllers.Login())
}