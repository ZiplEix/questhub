package database

import (
	"time"
)

type ChatMessage struct {
	ID         string    `json:"id"`
	GameID     string    `json:"game_id"`
	SenderID   string    `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	Content    string    `json:"content"`
	Type       string    `json:"type"` // "CHAT_GLOBAL", "CHAT_PRIVATE", "EVENT"
	TargetID   *string   `json:"target_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
