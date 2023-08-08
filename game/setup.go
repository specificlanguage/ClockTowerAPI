package game

import (
	"encoding/json"
	"errors"
	"log"
)

type GameSetupMessage struct {
	NumPlayers int      `json:"numPlayers"`
	Roles      []string `json:"roles"`
}

func setupGame(message json.RawMessage, sess GameSess) error {
	var setupData = GameSetupMessage{}
	err := json.Unmarshal(message, &setupData)
	if err != nil {
		log.Fatalln("error:", err)
	}

	if len(sess.Clients)-1 != setupData.NumPlayers {
		return errors.New("Not enough players")
	}

	ShuffleRoles(setupData.Roles, MapValsToPointerList(GetNonStorytellerPlayers(sess)))

	return nil
}
