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
	GameCode string `json:"gameCode"`
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

	gameCode := createGame.GameCode
	if gameCode == "" {
		gameCode = generateGameCode()
	}

	if _, ok := Games[gameCode]; !ok {
		gameCode = generateGameCode()
	}

	// Regenerate game code if already using, somehow
	game_db := db.Game{Code: gameCode, ScriptID: createGame.ScriptId, StorytellerUUID: uuid.MustParse(storyUUID), DBUUID: uuid.New()}

	if result := db.GameDB.Create(&game_db); result.Error != nil {
		log.Printf("%sDB Write Error: %s", logger.Red, result.Error.Error())

		// TODO: delete entry from table if game is too old (~1 day, say?)

		gameCode = generateGameCode()
		// Attempt reassigning game code
		if _, ok := Games[gameCode]; ok {
			game_db = db.Game{Code: gameCode, ScriptID: createGame.ScriptId, StorytellerUUID: uuid.MustParse(storyUUID), DBUUID: uuid.New()}
			if result := db.GameDB.Create(&game_db); result.Error != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create game"})
				return
			}
		}
	}

	inChannel := make(chan game.MessageFromClient)
	outChannel := OutChannel

	sess := game.GameSess{Code: gameCode, Clients: make(map[uuid.UUID]*game.Player), InChannel: inChannel, OutChannel: outChannel, Phase: game.GAME_LOBBY, DBUUID: game_db.DBUUID}

	// Start one thread for game logic
	go game.GameHandler(sess)

	// Make router recognize this game exists
	Games[gameCode] = sess

	ctx.JSON(http.StatusOK, gin.H{"code": gameCode, "script": "Trouble Brewing", "storyteller": storyUUID})
}
