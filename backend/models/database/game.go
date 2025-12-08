package database

import "time"

type Game struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	GmID       string    `json:"gm_id"`
	GmName     string    `json:"gm_name"`
	InviteCode string    `json:"invite_code"`
	IsActive   bool      `json:"is_active"`
	ImageURL   string    `json:"image_url"`
	Notes      string    `json:"notes"`
	State      string    `json:"state"` // "ongoing" or "paused"
	CreatedAt  time.Time `json:"created_at"`
}
