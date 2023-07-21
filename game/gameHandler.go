package game

import (
	"ClockTowerAPI/game/roles"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

const (
	MESSAGE           = "MESSAGE"
	CLIENT_JOIN       = "CLIENT_JOIN"
	CLIENT_DISCONNECT = "CLIENT_DISCONNECT"
	DISCONNECT        = "DISCONNECT"
	ERROR             = "ERROR"
)

type Player struct {
	UUID   uuid.UUID
	Name   string
	GameID string
	Role   roles.Role // please note that "nil" means that the client is the storyteller.
}

type GameSess struct {
	Code       string
	Clients    map[string]Player
	InChannel  chan map[string]interface{} // Generic info channel to send information to send. Will specify type later.
	OutChannel chan MessageToClient
}

type MessageToClient struct {
	Message Message
	UUIDs   []uuid.UUID
}

type Message struct {
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
}

func GameHandler(gh *GameSess) {
	for {
		select {
		case msg := <-gh.InChannel:
			fmt.Println(msg)
		}
	}
}

func AddPlayerToGame() {

}

// M - Alias for MakeMessage
func M(msgType string, message map[string]any, clients []uuid.UUID) MessageToClient {
	return MakeMessage(msgType, message, clients)
}

// MakeMessage - Creates a MessageToClient item for use
func MakeMessage(msgType string, message map[string]any, clients []uuid.UUID) MessageToClient {
	msgBytes, _ := json.Marshal(message)
	return MessageToClient{
		Message: Message{Type: msgType, Message: msgBytes},
		UUIDs:   clients,
	}
}

func GetUUIDS(sess *GameSess) *[]uuid.UUID {
	uuids := make([]uuid.UUID, len(sess.Clients))
	i := 0
	for _, player := range sess.Clients {
		uuids[i] = player.UUID
		i += 1
	}
	return &uuids
}
