package http

import (
	"ClockTowerAPI/db"
	game "ClockTowerAPI/game"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm/logger"
	"log"
	"math/rand"
	"net/http"
)

var letters = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")

type CreateGameBody struct {
	ScriptId string `json:"scriptID"`
}

func generateGameCode() string {

	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)

}

func CreateGameEndpoint(ctx *gin.Context) {
	var createGame CreateGameBody
	storyUUID := ctx.GetString("uuid")
	reqErr := ctx.BindJSON(&createGame)
	if reqErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing scriptID"})
		return
	}

	gameCode := generateGameCode()
	game_db := db.Game{Code: gameCode, ScriptID: createGame.ScriptId, StorytellerUUID: uuid.MustParse(storyUUID)}
	result := db.GameDB.Create(&game_db)

	inChannel := make(chan map[string]interface{})
	outChannel := make(chan game.MessageToClient)

	sess := game.GameSess{Code: gameCode, Clients: map[string]game.Player{}, InChannel: inChannel, OutChannel: outChannel}

	// Start one thread for game logic, another for dispatching out to the clients.
	go game.GameHandler(&sess)
	go Dispatcher(outChannel)

	if result.Error != nil {
		log.Printf("%sDB Write Error: %s", logger.Red, result.Error.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create game"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"code": gameCode, "script": "Trouble Brewing", "storyteller": storyUUID})
	}
}
