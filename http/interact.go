package http

import (
	"ClockTowerAPI/game"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func InteractEndpoint(ctx *gin.Context) {
	gameID := ctx.Param("id")
	name := ctx.Query("name")
	if gameID[len(gameID)-1:] == "/" {
		gameID = gameID[:len(gameID)-1]
	} else if gameID == "" {
		ctx.JSON(400, gin.H{"message": "Did not provide gameSess code"})
		return
	}

	if name == "" {
		ctx.JSON(400, gin.H{"message": "Provide name in params"})
		return
	}

	clientUUID, uuidErr := uuid.Parse(ctx.Query("uuid"))
	if uuidErr != nil {
		ctx.JSON(401, gin.H{"message": "Did not provide UUID in headers"})
		return
	}

	gameSess, ok := Games[gameID]
	if !ok {
		ctx.JSON(400, gin.H{"message": "Game does not exist"})
		return
	}

	if cl, ok := gameSess.Clients[clientUUID]; ok {
		// Reconnected
		cl.IsConnected = true
		gameSess.Clients[clientUUID] = cl
	} else {
		// New player
		if len(gameSess.Clients) >= 16 {
			ctx.JSON(401, gin.H{"message": "Game is full"})
			return
		}
		gameSess.Clients[clientUUID] = game.Player{UUID: clientUUID, Name: name, IsConnected: true}
	}

	// Upgrades request
	err := Mel.HandleRequestWithKeys(ctx.Writer, ctx.Request, gin.H{"uuid": clientUUID, "gid": gameID, "game": gameSess, "name": name})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not instantiate Websocket"})
		return
	}
}
