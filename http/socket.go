package http

import (
	"ClockTowerAPI/game/roles"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	"time"
)

const timeoutPeriod = 20 * time.Minute

type Client struct {
	uuid       uuid.UUID
	name       string
	websession *melody.Session
	role       roles.Role // please note that "nil" means that the client is the storyteller.
}

type GameSess struct {
	code       string
	clients    map[string]Client
	inChannel  chan map[string]interface{} // Generic info channel to send information to send. Will specify type later.
	outChannel chan MessageToClient
}

type MessageToClient struct {
	msg         Message
	clientUUIDs []uuid.UUID
}

type Message struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"message"`
}

// GameSessHandler is a Goroutine that should only handle one game at a time. Please be careful!
// The intention is that this will autodelete after an hr or so, but for now it will continue to keep going.
func GameSessHandler(gh *GameSess) {

	timer := time.NewTicker(timeoutPeriod)

	go func(gh *GameSess) {

		defer func() {
			msg, err := json.Marshal(Message{
				Type:  "DISCONNECT",
				Value: json.RawMessage("Room timed out, disconnect"),
			})

			if err != nil {
				fmt.Println()
				for _, sess := range GetClientSessions(gh.clients) {
					sess.Close()
				}
			} else {
				for _, sess := range GetClientSessions(gh.clients) {
					sess.CloseWithMsg(msg)
				}
			}
		}()

		for {
			select {
			case msg := <-gh.inChannel: // Should come from the Melody handler
				timer.Reset(timeoutPeriod) // Reset timeout timer, something happened.
				fmt.Println(msg)
			case msg := <-gh.outChannel: // Should come from the game loop
				timer.Reset(timeoutPeriod)
				fmt.Println(msg)
			case <-timer.C: // No activity, timing out.
				return
			}
		}

	}(gh)

}

func GetClientSessions(clients map[string]Client) []melody.Session {
	sessions := make([]melody.Session, len(clients))
	i := 0
	for _, client := range clients {
		sessions[i] = *client.websession
	}
	return sessions
}
