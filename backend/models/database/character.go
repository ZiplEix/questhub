package database

import (
	"encoding/json"
	"time"
)

type Character struct {
	ID         string          `json:"id"`
	GameID     string          `json:"game_id"`
	UserID     *string         `json:"user_id"` // Pointer because it can be NULL
	Name       string          `json:"name"`
	Race       string          `json:"race"`
	MaxHP      int             `json:"max_hp"`
	CurrentHP  int             `json:"current_hp"`
	AvatarURL  *string         `json:"avatar_url"`
	Stats      json.RawMessage `json:"stats"`
	Inventory  json.RawMessage `json:"inventory"`
	IsNPC      bool            `json:"is_npc"`
	Money      int             `json:"money"`
	CreatedAt  time.Time       `json:"created_at"`
	PlayerName *string         `json:"player_name,omitempty"` // Populated when joining with user table

	// New fields
	Initiative int             `json:"initiative"`
	Age        string          `json:"age"`
	Height     string          `json:"height"`
	Weight     string          `json:"weight"`
	MaxSpells  int             `json:"max_spells"`
	Spells     json.RawMessage `json:"spells"`
	Abilities  string          `json:"abilities"`
	Experience int             `json:"experience"`
}
