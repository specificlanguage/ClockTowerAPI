package http

import (
	"ClockTowerAPI/game"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	"time"
)

const timeoutPeriod = 20 * time.Minute

type Client struct {
	UUID       uuid.UUID
	GameID     string
	Websession *melody.Session
}

// Dispatcher is a Goroutine that sends out all messages for a single game. Please be careful!
func Dispatcher(sendChannel chan game.MessageToClient) {

	timer := time.NewTicker(timeoutPeriod)

	for {
		select {
		case msg := <-sendChannel: // Should come from the game loop
			timer.Reset(timeoutPeriod)
			msgJSON, _ := json.Marshal(msg.Message)
			Mel.BroadcastMultiple(msgJSON, GetSessionsFromUUID(msg.UUIDs))
		case <-timer.C: // No activity, timing out.
			return
		}
	}
}

// Utility functions below to get sessions from various methods.

func GetSessionsFromClient(clients []Client) []*melody.Session {
	sessions := make([]*melody.Session, len(clients))
	for i, client := range clients {
		sessions[i] = client.Websession
	}
	return sessions
}

func GetSessionsFromUUID(uuids []uuid.UUID) []*melody.Session {
	sessions := make([]*melody.Session, len(uuids))
	for i, uuid := range uuids {
		sessions[i] = Clients[uuid].Websession
	}
	return sessions
}

func GetSessionsFromGameID(gameID string) []*melody.Session {
	sessions := make([]*melody.Session, 0)
	for _, client := range Clients {
		if client.GameID == gameID {
			sessions = append(sessions, client.Websession)
		}
	}
	return sessions
}
