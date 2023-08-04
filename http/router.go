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
var Games = make(map[string]game.GameSess)
var OutChannel = make(chan game.MessageToClient)

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
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowHeaders = []string{"*"}
	router.Use(cors.New(config))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ClockTowerAPI",
			"version": "0.0.1",
		})
	})

	router.GET("/script", game.GetScriptInfoEndpoint)
	router.GET("/game/:id", InteractEndpoint)

	gameGroup := router.Group("/game/")
	gameGroup.Use(UUIDRequired())
	{
		gameGroup.POST("/create", CreateGameEndpoint)
	}

	go Dispatcher()

	// Websocket handlers here
	Mel.HandleConnect(func(s *melody.Session) {
		go func() {
			fmt.Println(s.MustGet("uuid"))
			clientUUID := s.MustGet("uuid").(uuid.UUID) // Convert to string
			gid := s.MustGet("gid").(string)
			name := s.MustGet("name").(string)
			gameSess := s.MustGet("game").(game.GameSess)

			// Send notification to all existing clients
			joinMessage := gin.H{"name": name, "uuid": clientUUID, "gameID": gid}
			if gameSess.Clients[clientUUID].IsStoryteller {
				joinMessage["isStoryteller"] = true
			}

			gameSess.OutChannel <- game.M(
				game.CLIENT_JOIN,
				joinMessage,
				game.GetConnectedClientsUUIDs(gameSess),
				gid)

			// Send info about all players to the newly joined client
			gameSess.OutChannel <- game.M(
				game.GAME_INFO,
				gin.H{"players": game.GetConnectedPlayers(gameSess)},
				game.SingleToMap(clientUUID),
				gameSess.Code)
		}()
	})

	Mel.HandleMessage(func(s *melody.Session, msg []byte) {
		go func() {
			clientUUID := s.MustGet("uuid").(uuid.UUID) // Convert to string
			gameSess := s.MustGet("game").(game.GameSess)
			val := make(map[string]interface{})
			if err := json.Unmarshal(msg, &val); err != nil {
				gameSess.OutChannel <- game.M(
					game.ERROR,
					gin.H{"message": "Could not parse JSON"},
					game.SingleToMap(clientUUID),
					gameSess.Code)
				return
			}

			gameSess.InChannel <- val
		}()
	})

	// Notify all clients that a disconnect occurred. We don't remove it from gameSess in the case of a reconnect
	Mel.HandleDisconnect(func(s *melody.Session) {
		go func() {
			gameSess := s.MustGet("game").(game.GameSess)
			clientUUID := s.MustGet("uuid").(uuid.UUID) // Convert to string
			// Set disconnected
			player := gameSess.Clients[clientUUID]
			player.IsConnected = false
			gameSess.Clients[clientUUID] = player

			gameSess.OutChannel <- game.M(
				game.CLIENT_DISCONNECT,
				gin.H{"name": player.Name, "uuid": clientUUID, "gameID": player.GameID},
				game.GetConnectedClientsUUIDs(gameSess),
				gameSess.Code,
			)
		}()
	})

	return router
}

// Dispatcher - Goroutine that handles all outgoing messages
func Dispatcher() {
	for {
		select {
		case msg := <-OutChannel:
			msgJSON, _ := json.Marshal(msg.Message)
			Mel.BroadcastFilter(msgJSON,
				func(session *melody.Session) bool {
					clientUUID := session.MustGet("uuid").(uuid.UUID) // Convert to string
					gid := session.MustGet("gid").(string)
					if _, ok := msg.UUIDs[clientUUID]; ok && msg.GameID == gid && !session.IsClosed() {
						return true
					}
					return false
				},
			)
		}
	}
}
