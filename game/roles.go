package game

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
)

// Team signifiers
const (
	TOWNSFOLK = iota
	OUTSIDER
	MINION
	DEMON
	TRAVELER
)

// RoleHandler To be filled in later with every role, probably inputting a GameState & information
type RoleHandler func()

// Role Used for backend purposes to abstract out role information
type Role struct {
	RoleName           string `json:"role_name"`
	Description        string `json:"description"`
	Team               int    `json:"team"`
	FirstNightPriority int    `json:"firstNightPriority"`
	OtherNightPriority int    `json:"otherNightPriority"`
	Handler            map[string]interface{}
}

// RoleView used for processing in/out role information. This should not be used for internal purposes.
type RoleView struct {
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
	Team        int    `json:"team"`
}

// GetRoleDescription - Gets the description of one role in list.
// Should generally be used for custom scripts when we get there
func GetRoleDescription(role string) (*RoleView, error) {
	roles, err := GetAllRoleDescriptions()
	if err != nil {
		return nil, err
	}

	for _, element := range roles {
		if element.RoleName == role {
			return &element, nil
		}
	}
	return nil, errors.New("could not find role")
}

// GetAllRoleDescriptions - Gets all descriptions for a single script, currently only works for TB
func GetAllRoleDescriptions() ([]RoleView, error) {
	var roleData []RoleView
	all_roles, readErr := os.ReadFile("./scripts/roles.json")
	if readErr != nil {
		return nil, readErr
	}
	unmarshalErr := json.Unmarshal(all_roles, &roleData)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return roleData, nil
}

func GetAllRoles() ([]Role, error) {
	var roleData []Role
	all_roles, readErr := os.ReadFile("./scripts/roles.json")
	if readErr != nil {
		return nil, readErr
	}
	unmarshalErr := json.Unmarshal(all_roles, &roleData)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return roleData, nil
}

// ShuffleRoles - Shuffles roles and gives each role to a different player.
func ShuffleRoles(roleNames []string, sess GameSess) {

	n := len(roleNames)
	allRoles, err := GetAllRoles()
	usedRoles := make([]Role, n)

	if err != nil {
		log.Printf("Error when shuffling roleNames: %s", err)
	}

	// Associate role names with roles
	for i, name := range roleNames {
		for _, role := range allRoles {

			if role.RoleName == name {
				usedRoles[i] = role
				break
			}
		}
	}

	// Shuffle roles in list
	for i := 0; i < n-2; i++ {
		j := rand.Intn(n-i) + i
		usedRoles[i], usedRoles[j] = usedRoles[j], usedRoles[i]
	}

	i := 0
	// Distribute roles to players
	for _, player := range GetNonStorytellerPlayers(sess) {
		player.Role = usedRoles[i]
		i++
	}
}
