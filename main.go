package main

import (
	"fmt"
	"time"
	"log"

	"github.com/dannonb/go-network-monitor/config"
	"github.com/dannonb/go-network-monitor/helpers"
	"github.com/dannonb/go-network-monitor/middleware"
	"github.com/dannonb/go-network-monitor/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	// connect database
	config.ConnectDB()
	helpers.ClearCache()
	helpers.UpdateCache()
}

func main() {

	router := gin.Default()

	router.Use(gin.Logger())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"data": "Hello World",
		})
	})

	routes.UserAuthRoute(router)

	router.Use(middleware.Authentication())

	routes.UserRoute(router)
	routes.HostRoute(router)

	// start server
	go func() {
		router.Run()
	}()

	//hosts := []string{"mail.google.com", "example.com"}

	hosts, err := helpers.GetHostsFromCache()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HOSTS: ", hosts)
	for {
		fmt.Println("Monitoring...")
		for _, host := range hosts {
			go helpers.Monitor(host)
		}

		time.Sleep(helpers.Interval)
	}
}
