package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func InteractEndpoint(ctx *gin.Context) {
	gameID := ctx.Param("id")
	if gameID[len(gameID)-1:] == "/" {
		gameID = gameID[:len(gameID)-1]
	}

	if gameID == "" {
		ctx.JSON(400, gin.H{"message": "Did not provide game code"})
	}

	// TODO: query for game in database

	clientUUID, uuidErr := uuid.Parse(ctx.GetString("uuid"))
	if uuidErr != nil {
		ctx.JSON(403, gin.H{"message": "Did not provide UUID in headers"})
	}

	// TODO: check if players are full in game

	// Upgrades request
	err := Mel.HandleRequestWithKeys(ctx.Writer, ctx.Request, gin.H{"uuid": clientUUID, "gid": gameID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not instantiate Websocket"})
	}
}
