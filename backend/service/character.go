package service

import (
	"context"
	"database/sql"
	"questhub/database"
	model "questhub/models/database"
)

func GetUserCharacter(gameID, userID string) (*model.Character, error) {
	character := &model.Character{}
	query := `
		SELECT id, game_id, user_id, name, avatar_url, stats, inventory, is_npc, created_at
		FROM characters
		WHERE game_id = $1 AND user_id = $2
	`
	err := database.DB.QueryRow(context.Background(), query, gameID, userID).Scan(
		&character.ID,
		&character.GameID,
		&character.UserID,
		&character.Name,
		&character.AvatarURL,
		&character.Stats,
		&character.Inventory,
		&character.IsNPC,
		&character.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No character found, not an error
		}
		return nil, err
	}
	return character, nil
}
