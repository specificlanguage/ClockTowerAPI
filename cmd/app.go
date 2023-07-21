package cmd

import (
	"ClockTowerAPI/db"
	"ClockTowerAPI/http"
	"github.com/olahol/melody"
)

var Mel = melody.New()

func Run() {
	db.Init()
	// Now available as GameDB from this point forward

	// Initialize webserver
	r := http.SetupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
