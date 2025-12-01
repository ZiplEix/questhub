package database

import "time"

type Game struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	GmID       string    `json:"gm_id"`
	InviteCode string    `json:"invite_code"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}
