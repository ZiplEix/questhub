package database

import (
	"encoding/json"
	"time"
)

type Template struct {
	ID          string          `json:"id"`
	CreatedBy   string          `json:"created_by"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Type        string          `json:"type"` // CHARACTER, NPC, MONSTER, CREATURE
	Data        json.RawMessage `json:"data"`
	IsPublic    bool            `json:"is_public"`
	Uses        int             `json:"uses"`
	CreatedAt   time.Time       `json:"created_at"`
}
