package http

import (
	"ClockTowerAPI/db"
	"ClockTowerAPI/game"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func InteractEndpoint(ctx *gin.Context) {
	gameID := ctx.Param("id")
	if gameID == "" {
		gameID = ctx.Query("id")
	}
	name := ctx.Query("name")
	if gameID[len(gameID)-1:] == "/" {
		gameID = gameID[:len(gameID)-1]
	} else if gameID == "" {
		ctx.JSON(400, gin.H{"message": "Did not provide gameSess code"})
		return
	}
	gameID = strings.ToUpper(gameID) // In case it's a lowercase item

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
		if cl.IsConnected {
			ctx.JSON(400, gin.H{"message": "Already connected"})
		}
		// Reconnected
		cl.IsConnected = true
		gameSess.Clients[clientUUID] = cl
	} else {
		// New player
		if len(gameSess.Clients) > 16 {
			ctx.JSON(401, gin.H{"message": "Game is full"})
			return
		}

		// Check if storyteller, which means going to the DB
		gameEntry := db.GetGameByID(gameID)
		if gameEntry == nil {
			ctx.JSON(400, gin.H{"message": "Game does not exist"})
			return
		}

		if clientUUID != gameEntry.StorytellerUUID {
			gameSess.Clients[clientUUID] = &game.Player{UUID: clientUUID, Name: name, IsConnected: true}
		} else {
			gameSess.Clients[clientUUID] = &game.Player{UUID: clientUUID, Name: name, IsConnected: true, IsStoryteller: true}
		}
	}

	// Upgrades request
	err := Mel.HandleRequestWithKeys(ctx.Writer, ctx.Request, gin.H{"uuid": clientUUID, "gid": gameID, "game": gameSess, "name": name})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not instantiate Websocket"})
		return
	}
}
