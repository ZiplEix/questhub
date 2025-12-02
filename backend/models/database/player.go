package database

import "time"

type Player struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	IsGM      bool      `json:"is_gm"`
	JoinedAt  time.Time `json:"joined_at"`
}
