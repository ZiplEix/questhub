package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Character struct {
	ID        string         `json:"id"`
	GameID    string         `json:"game_id"`
	UserID    *string        `json:"user_id"`
	Name      string         `json:"name"`
	AvatarURL string         `json:"avatar_url"`
	IsNPC     bool           `json:"is_npc"`
	Stats     CharacterStats `json:"stats"`
	Inventory Inventory      `json:"inventory"`
	CreatedAt time.Time      `json:"created_at"`
}

// --- CharacterStats ---
type CharacterStats map[string]any

func (a CharacterStats) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *CharacterStats) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// --- Inventory ---

type InventoryItem struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Quantity    int    `json:"quantity"`
}

type Inventory []InventoryItem

func (a Inventory) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Inventory) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
