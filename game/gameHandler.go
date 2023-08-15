package game

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

const (
	MESSAGE           = "MESSAGE"
	CLIENT_JOIN       = "CLIENT_JOIN"
	CLIENT_DISCONNECT = "CLIENT_DISCONNECT"
	ERROR             = "ERROR"
	GAME_INFO         = "GAME_INFO"   // Used for one-way info from S2C, like a phase change or a list of players on entry
	GAME_SETUP        = "GAME_SETUP"  // Only used for starting games, should not be used in any other circumstance
	GAME_ACTION       = "GAME_ACTION" // Used for two-way info, first to ask Client for selections, and then response in same way for Server to process selections
)

const (
	LOADING      = "LOADING" // Default when entering room
	GAME_LOBBY   = "LOBBY"   // Landing page on lobby
	GAME_START   = "GAME_START"
	NIGHT        = "NIGHT"        // Most logic will happen at night
	POSTNIGHT    = "POSTNIGHT"    // Should be used for post
	DAY          = "DAY"          // General deliberation for the day
	NOMINATE     = "NOMINATE"     // Nominating voters
	DELIBERATION = "DELIBERATION" // Should be used as a stopgap between nominations and voting
	VOTING       = "VOTING"
	POSTVOTE     = "POSTVOTE" // Should be used for resolution events regarding the voting, and special roles like Vizier
	POSTDAY      = "POSTDAY"
)

// Player represents a player within the game with:
// - a unique UUID,
// - player Name,
// - associated GameID,
// - a Role that's being taken,
// - whether the player IsStoryteller, IsConnected to the client, IsDrunk, or IsPoisoned with respect to the game.
type Player struct {
	UUID          uuid.UUID
	Name          string
	GameID        string
	Role          Role // please note that "nil" means that the client is the storyteller.
	IsStoryteller bool
	IsConnected   bool
	IsDrunk       bool
	IsPoisoned    bool
}

// GameSess represents game-wide information within the game that's not specific to a certain player with:
// - a game Code (6 char string)
// - a map of players listed as Clients, with the key being the client's uuids.
// - a channel for incoming messages, which should be incoming from the client.
// - a channel for outgoing messages, which should only be connected to the Dispatcher
type GameSess struct {
	Code        string
	Clients     map[uuid.UUID]*Player
	InChannel   chan MessageFromClient // Generic info channel to send information to send. Will specify type later.
	OutChannel  chan MessageToClient
	Phase       string
	GameHandler map[string]interface{}
	DBUUID      uuid.UUID
}

// MessageToClient represents an outgoing message from the server to the client.
// - The Message to be sent
// - The game ID of the GameSess that it should be sent to
// - A map of UUIDs to send to. It's listed as a map for convenience as GameSess player lists are maps by default.
// An easy way to construct a MessageToClient struct is to use the game.M function.
type MessageToClient struct {
	Message Message
	GameID  string
	UUIDs   map[uuid.UUID]any
}

// MessageFromClient represents an incoming message from the client to the server.
// - The incoming Message
// - The GameID of the GameSess that it arrives from
// - The GameSess that the client is associated with
// - The UUID of the client.
// An easy way to construct one of these when parsing is to use the game.P function.
type MessageFromClient struct {
	Message  Message
	GameSess GameSess
	UUID     uuid.UUID
}

// Message represents the message format being sent back and forth.
// - The Type of the message - e.g. CLIENT_JOIN, GAME_INFO, GAME_SETUP, etc.
// - The Message itself - containing the actual info.
type Message struct {
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
}

// GameHandler is a goroutine that handles any GameSess actions.
func GameHandler(sess GameSess) {
	for {
		select {
		case msg := <-sess.InChannel:
			// fmt.Println("Inbound", msg)
			if msg.Message.Type == "GAME_SETUP" {
				if err := setupGame(msg.Message.Message, sess); err != nil {
					fmt.Println("Error: ", err)
				}
			} else { // For debugging purposes
				sess.OutChannel <- M(MESSAGE, gin.H{"type": msg.Message.Type, "message": msg.Message.Message}, MapToAnyMap(sess.Clients), sess.Code)
			}
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

// P - Alias for ParseIncomingMessage
func P(rawMsg []byte, sess GameSess, uuid uuid.UUID) (*MessageFromClient, error) {
	return ParseIncomingMessage(rawMsg, sess, uuid)
}

// ParseIncomingMessage - Creates a MessageFromClient by parsing the rawMessage.
func ParseIncomingMessage(rawMsg []byte, sess GameSess, uuid uuid.UUID) (*MessageFromClient, error) {

	msgPtr := new(Message)
	if err := json.Unmarshal(rawMsg, msgPtr); err != nil {
		log.Printf("Could not parse message from client %s:\n %s", uuid.String(), err.Error())
		return nil, err
	}

	return &MessageFromClient{
		Message:  *msgPtr,
		GameSess: sess,
		UUID:     uuid,
	}, nil

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

func GetNonStorytellerPlayers(sess GameSess) map[uuid.UUID]*Player {
	players := make(map[uuid.UUID]*Player)
	for _, player := range sess.Clients {
		if player.IsStoryteller {
			continue
		} else if player.IsConnected {
			players[player.UUID] = player
		}
	}
	return players
}

func GetStoryteller(sess GameSess) map[uuid.UUID]any {
	players := make(map[uuid.UUID]any)
	for _, player := range sess.Clients {
		if player.IsStoryteller {
			players[player.UUID] = *player
			return players
		}
	}
	return players
}
