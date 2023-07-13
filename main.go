package main

import (
	"ClockTowerAPI/db"
	"ClockTowerAPI/game"
	"ClockTowerAPI/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.ForwardedByClientIP = true
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic(err.Error())
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ClockTower Backend",
			"version": "0.0.1",
		})
	})

	gameGroup := r.Group("/game/")
	gameGroup.Use(middleware.UUIDRequired())
	{
		gameGroup.POST("/create", game.CreateGameEndpoint)
	}

	return r
}

func main() {

	// TODO: Load database
	db.Init()
	// Now available as GameDB from this point forward

	// Initialize webserver
	r := SetupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
