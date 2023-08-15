package game

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
)

type GameSetupMessage struct {
	NumPlayers int      `json:"numPlayers"`
	Roles      []string `json:"roles"`
}

type RoleInfo struct {
	PlayerUUID string
	RoleName   string
}

func setupGame(message json.RawMessage, sess GameSess) error {
	var setupData = GameSetupMessage{}
	err := json.Unmarshal(message, &setupData)
	if err != nil {
		log.Fatalln("error:", err)
	}

	ShuffleRoles(setupData.Roles, sess)

	// Post role shuffle to check other things like Drunks, other items

	// Send message to all players of their roles, also sets up message to storyteller of whose role is who
	roleInfoList := make([]RoleInfo, setupData.NumPlayers)
	i := 0
	for _, player := range sess.Clients {
		// fmt.Println(player.Role.RoleName)
		if !player.IsStoryteller && player.IsConnected {
			roleInfoList[i] = RoleInfo{player.UUID.String(), player.Role.RoleName}
			i++
			sess.OutChannel <- M("GAME_SETUP", gin.H{
				"role":    player.Role.RoleName,
				"phase":   GAME_START,
				"players": setupData.NumPlayers,
			}, SingleToMap(player.UUID), sess.Code)
		}
	}

	sess.OutChannel <- M("GAME_SETUP", gin.H{
		"roles": roleInfoList,
	}, GetStoryteller(sess), sess.Code)

	return nil
}
