package game

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
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
	storytellerUUID := ctx.MustGet("uuid")
	fmt.Println(storytellerUUID)
	// TODO: add game to DB
}
