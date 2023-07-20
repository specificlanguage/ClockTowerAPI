package http

import (
	"ClockTowerAPI/game"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"net/http"
)

var Mel = melody.New()

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
		gameGroup.GET("/:id", InteractEndpoint)
		gameGroup.POST("/create", CreateGameEndpoint)
	}

	Mel.HandleConnect(func(s *melody.Session) {
		clientUUID := s.MustGet("uuid")
		gid := s.MustGet("gid")
		Mel.Broadcast([]byte(fmt.Sprintf("%s connected to game %s", clientUUID, gid)))
	})

	Mel.HandleMessage(func(s *melody.Session, msg []byte) {
		Mel.Broadcast(msg)
	})

	return router
}
