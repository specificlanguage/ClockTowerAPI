package roles

import (
	"encoding/json"
	"errors"
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
	roleName    string
	description string
	team        int
	priority    int // Use for priority when processing events at night
	handler     RoleHandler
}

// RoleView used for processing in/out role information. This should not be used for internal purposes.
type RoleView struct {
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
	Team        uint   `json:"team"`
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
