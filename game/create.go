package game

import (
	"ClockTowerAPI/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"log"
	"math/rand"
	"net/http"
)

var letters = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateGameCode() string {

	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)

}

func CreateGameEndpoint(ctx *gin.Context) {
	storyUUID := ctx.GetString("uuid")
	gameCode := generateGameCode()
	game := db.Game{Code: gameCode, StorytellerUUID: storyUUID}
	result := db.GameDB.Create(&game)
	if result.Error != nil {
		log.Printf("%sDB Write Error: %s", logger.Red, result.Error.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create game"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"code": gameCode, "script": "Trouble Brewing", "storyteller": storyUUID})
	}
}
