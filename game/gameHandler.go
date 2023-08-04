package game

import (
	"ClockTowerAPI/game/roles"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	MESSAGE           = "MESSAGE"
	CLIENT_JOIN       = "CLIENT_JOIN"
	CLIENT_DISCONNECT = "CLIENT_DISCONNECT"
	ERROR             = "ERROR"
	GAME_INFO         = "GAME_INFO"
)

const (
	LOADING      = "LOADING"      // Default when entering room
	GAME_LOBBY   = "LOBBY"        // Landing page on lobby
	NIGHT        = "NIGHT"        // Most logic will happen at night
	POSTNIGHT    = "POSTNIGHT"    // Should be used for post
	DAY          = "DAY"          // General deliberation for the day
	NOMINATE     = "NOMINATE"     // Nominating voters
	DELIBERATION = "DELIBERATION" // Should be used as a stopgap between nominations and voting
	VOTING       = "VOTING"
	POSTVOTE     = "POSTVOTE" // Should be used for resolution events regarding the voting, and special roles like Vizier
	POSTDAY      = "POSTDAY"
)

type Player struct {
	UUID          uuid.UUID
	Name          string
	GameID        string
	Role          roles.Role // please note that "nil" means that the client is the storyteller.
	IsStoryteller bool
	IsConnected   bool
	IsDrunk       bool
	IsPoisoned    bool
}

type GameSess struct {
	Code       string
	Clients    map[uuid.UUID]Player
	InChannel  chan map[string]interface{} // Generic info channel to send information to send. Will specify type later.
	OutChannel chan MessageToClient
	Phase      string
}

type GameState struct {
	Phase string
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

// GetConnectedPlayers - Returns a list of all connected players with only name, uuid, and isStoryteller.
// Should only be used for non-Storyteller information
func GetConnectedPlayers(sess GameSess) map[string]any {
	players := make(map[string]any)
	for _, player := range sess.Clients {
		playerRedacted := gin.H{"uuid": player.UUID, "name": player.Name}
		if player.IsStoryteller {
			playerRedacted["isStoryteller"] = true
		}
		if player.IsConnected {
			players[player.UUID.String()] = playerRedacted
		}
	}
	return players
}
