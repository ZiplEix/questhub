package database

import "time"

type Invitation struct {
	ID        string    `json:"id"`
	GameID    string    `json:"game_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}
