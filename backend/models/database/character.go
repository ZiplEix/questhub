package database

import (
	"encoding/json"
	"time"
)

type Character struct {
	ID        string          `json:"id"`
	GameID    string          `json:"game_id"`
	UserID    *string         `json:"user_id"` // Pointer because it can be NULL
	Name      string          `json:"name"`
	AvatarURL *string         `json:"avatar_url"`
	Stats     json.RawMessage `json:"stats"`
	Inventory json.RawMessage `json:"inventory"`
	IsNPC     bool            `json:"is_npc"`
	CreatedAt time.Time       `json:"created_at"`
}
