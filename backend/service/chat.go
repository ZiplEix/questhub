package service

import (
	"context"
	"questhub/database"
	model "questhub/models/database"
)

func SaveMessage(msg model.ChatMessage) error {
	_, err := database.DB.Exec(context.Background(),
		`INSERT INTO messages (game_id, sender_id, sender_name, content, type, target_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		msg.GameID, msg.SenderID, msg.SenderName, msg.Content, msg.Type, msg.TargetID, msg.CreatedAt)
	return err
}

func GetGameMessages(gameID string) ([]model.ChatMessage, error) {
	rows, err := database.DB.Query(context.Background(),
		`SELECT id, game_id, sender_id, sender_name, content, type, target_id, created_at
		FROM messages
		WHERE game_id = $1
		ORDER BY created_at ASC
		LIMIT 100`, // Limit to last 100 messages for now
		gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []model.ChatMessage{}
	for rows.Next() {
		var msg model.ChatMessage
		err := rows.Scan(&msg.ID, &msg.GameID, &msg.SenderID, &msg.SenderName, &msg.Content, &msg.Type, &msg.TargetID, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
