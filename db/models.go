package db

import (
	"github.com/google/uuid"
	"time"
)

type Game struct {
	ID              uint
	Code            string `gorm:"unique"`
	ScriptID        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	StorytellerUUID uuid.UUID
	Players         []GamePlayer `gorm:"foreignKey:GameID"`
	Phase           uint         `gorm:"default:0"` // Will be a lookup in the phase list
	// TODO: possibly convert Phase into a Redis cache for easy updating on night phases
}

type GamePlayer struct {
	ID     uint
	UUID   uuid.UUID
	GameID uint
	Role   uint `gorm:"default:0"` // This will eventually just be a lookup in the role table, with each role having an ID
}
