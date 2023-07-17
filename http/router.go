package http

import (
	"ClockTowerAPI/game"
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
			"message": "ClockTowerAPI",
			"version": "0.0.1",
		})
	})

	router.GET("/script", game.GetScriptInfoEndpoint)

	gameGroup := router.Group("/game/")
	gameGroup.Use(UUIDRequired())
	{
		gameGroup.POST("/create", game.CreateGameEndpoint)
	}

	return router
}
