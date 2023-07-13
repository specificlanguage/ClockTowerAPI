package db

import "time"

type Game struct {
	ID              uint
	Code            string
	ScriptID        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	StorytellerUUID string
	Players         []GamePlayer `gorm:"foreignkey:UUID"`
	Phase           uint         `gorm:"default:0"` // Will be a lookup in the phase list
	// TODO: possibly convert Phase into a Redis cache for easy updating on night phases
}

type GamePlayer struct {
	ID   uint
	UUID string
	Role uint `gorm:"default:0"` // This will eventually just be a lookup in the role table, with each role having an ID
}

type Action struct {
	ID           uint
	Game         Game
	Type         uint `gorm:"default:0"` // TODO: implement ActionTypes // Lookup in Action list
	TargetPlayer []GamePlayer
	TargetRole   []uint
	Description  string // Give a description of what the "action" was that night
}
