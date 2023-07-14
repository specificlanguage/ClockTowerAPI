package main

import (
	"ClockTowerAPI/db"
	"ClockTowerAPI/game"
	"ClockTowerAPI/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Trusted Proxies
	router.ForwardedByClientIP = true
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic(err.Error())
	}

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"*"}
	router.Use(cors.New(config))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ClockTower Backend",
			"version": "0.0.1",
		})
	})

	gameGroup := router.Group("/game/")
	gameGroup.Use(middleware.UUIDRequired())
	{
		gameGroup.POST("/create", game.CreateGameEndpoint)
	}

	return router
}

func main() {

	// TODO: Load database
	db.Init()
	// Now available as GameDB from this point forward

	// Initialize webserver
	r := SetupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
