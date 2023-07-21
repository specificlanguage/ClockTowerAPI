package http

import (
	"ClockTowerAPI/game"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	"net/http"
)

var Mel = melody.New()
var Clients = make(map[uuid.UUID]*Client)
var Games = make(map[string]*game.GameSess)

func UUIDRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid := ctx.GetHeader("Authorization")
		ctx.Set("uuid", uuid)
	}
}

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

	// Websocket handlers here
	Mel.HandleConnect(func(s *melody.Session) {
		clientUUID := uuid.MustParse(s.MustGet("uuid").(string)) // Convert to string
		gid := s.MustGet("gid").(string)
		cl := Client{clientUUID, gid, s}
		Clients[clientUUID] = &cl
		// gameCoord := Games[gid]
		// gameCoord.Clients[clientUUID.String()] = clientUUID
		gameSessions := GetSessionsFromGameID(gid)
		Mel.BroadcastMultiple([]byte(fmt.Sprintf("%s connected to game %s", clientUUID, gid)), gameSessions)
	})

	Mel.HandleMessage(func(s *melody.Session, msg []byte) {
		clientUUID := uuid.MustParse(s.MustGet("uuid").(string)) // Convert to string
		gid := s.MustGet("gid").(string)
		gameCoord := Games[gid]
		sentMsg := make(map[string]interface{})
		err := json.Unmarshal(msg, &sentMsg)

		if err != nil {
			errMsg := gin.H{"message": "Could not parse message"}
			raw, _ := json.Marshal(errMsg)
			gameCoord.OutChannel <- game.MessageToClient{
				Message: game.Message{Type: "ERROR", Message: raw},
				UUIDs:   []uuid.UUID{clientUUID}}
		}

		// Just acknowledgement of message if not done.
		if gin.Mode() == gin.TestMode {
			Mel.BroadcastMultiple([]byte(fmt.Sprintf("Message acknowledged")), []*melody.Session{s})
		}
		gameCoord.InChannel <- sentMsg
	})

	return router
}
