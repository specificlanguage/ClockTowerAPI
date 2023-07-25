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
	ERROR             = "ERROR"
	GAME_INFO         = "GAME_INFO"
)

type Player struct {
	UUID        uuid.UUID
	Name        string
	GameID      string
	Role        roles.Role // please note that "nil" means that the client is the storyteller.
	IsConnected bool
}

type GameSess struct {
	Code       string
	Clients    map[uuid.UUID]Player
	InChannel  chan map[string]interface{} // Generic info channel to send information to send. Will specify type later.
	OutChannel chan MessageToClient
}

type MessageToClient struct {
	Message Message
	GameID  string
	UUIDs   map[uuid.UUID]any
}

type Message struct {
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
}

func GameHandler(gh GameSess) {
	for {
		select {
		case msg := <-gh.InChannel:
			fmt.Println("Inbound", msg)
			gh.OutChannel <- M(MESSAGE, msg, MapToAnyMap(gh.Clients), gh.Code)
		}
	}
}

// M - Alias for MakeMessage
func M(msgType string, message map[string]any, clients map[uuid.UUID]any, gameID string) MessageToClient {
	return MakeMessage(msgType, message, clients, gameID)
}

// MakeMessage - Creates a MessageToClient item for use
func MakeMessage(msgType string, message map[string]any, clients map[uuid.UUID]any, gameID string) MessageToClient {
	msgBytes, _ := json.Marshal(message)
	return MessageToClient{
		Message: Message{Type: msgType, Message: msgBytes},
		GameID:  gameID,
		UUIDs:   clients,
	}
}

func GetConnectedClientsUUIDs(sess GameSess) map[uuid.UUID]any {
	uuids := make(map[uuid.UUID]any)
	for _, player := range sess.Clients {
		if player.IsConnected {
			uuids[player.UUID] = true
		}
	}
	return uuids
}
