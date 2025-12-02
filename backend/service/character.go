package service

import (
	"context"
	"database/sql"
	"fmt"
	"questhub/database"
	model "questhub/models/database"
	"time"
)

func GetUserCharacter(gameID, userID string) (*model.Character, error) {
	char := &model.Character{}
	query := `
		SELECT id, game_id, user_id, name, avatar_url, stats, inventory, is_npc, created_at
		FROM characters
		WHERE game_id = $1 AND user_id = $2
	`
	err := database.DB.QueryRow(context.Background(), query, gameID, userID).Scan(
		&char.ID,
		&char.GameID,
		&char.UserID,
		&char.Name,
		&char.AvatarURL,
		&char.Stats,
		&char.Inventory,
		&char.IsNPC,
		&char.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No character found, not an error
		}
		return nil, err
	}
	return char, nil
}

func GetGameCharacters(gameID string) ([]model.Character, error) {
	characters := []model.Character{}
	query := `
		SELECT c.id, c.game_id, c.user_id, c.name, c.race, c.max_hp, c.current_hp, COALESCE(c.avatar_url, ''), c.stats, c.inventory, c.is_npc, c.money, c.created_at, u.name
		FROM characters c
		LEFT JOIN "user" u ON c.user_id = u.id
		WHERE c.game_id = $1
		ORDER BY c.created_at DESC
	`
	rows, err := database.DB.Query(context.Background(), query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var char model.Character
		var playerName sql.NullString
		if err := rows.Scan(&char.ID, &char.GameID, &char.UserID, &char.Name, &char.Race, &char.MaxHP, &char.CurrentHP, &char.AvatarURL, &char.Stats, &char.Inventory, &char.IsNPC, &char.Money, &char.CreatedAt, &playerName); err != nil {
			return nil, err
		}
		if playerName.Valid {
			char.PlayerName = &playerName.String
		}
		characters = append(characters, char)
	}

	return characters, nil
}

type InventoryItem struct {
	Name     string  `json:"name"`
	Quantity string  `json:"quantity"`
	ImageURL *string `json:"image_url,omitempty"`
	IconName string  `json:"icon_name,omitempty"`
}

func CreateCharacter(gameID, userID, name, race string, maxHP int, isNPC bool, avatarURL string, stats, inventory []byte, money int) (*model.Character, error) {
	char := &model.Character{
		GameID:    gameID,
		Name:      name,
		Race:      race,
		MaxHP:     maxHP,
		CurrentHP: maxHP,
		IsNPC:     isNPC,
		Stats:     stats,
		Inventory: inventory,
		Money:     money,
		CreatedAt: time.Now(),
	}

	if userID != "" {
		char.UserID = &userID
	}
	if avatarURL != "" {
		char.AvatarURL = &avatarURL
	}

	query := `
		INSERT INTO characters (game_id, user_id, name, race, max_hp, current_hp, is_npc, avatar_url, stats, inventory, money, created_at)
		VALUES ($1, NULLIF($2, ''), $3, $4, $5, $6, $7, NULLIF($8, ''), $9, $10, $11, $12)
		RETURNING id
	`

	err := database.DB.QueryRow(context.Background(), query,
		char.GameID, userID, char.Name, char.Race, char.MaxHP, char.CurrentHP, char.IsNPC, avatarURL, char.Stats, char.Inventory, char.Money, char.CreatedAt,
	).Scan(&char.ID)

	if err != nil {
		return nil, err
	}

	return char, nil
}

func UpdateCharacter(id, gameID, name, race string, maxHP int, isNPC bool, avatarURL string, stats, inventory []byte, money int) (*model.Character, error) {
	query := `
		UPDATE characters
		SET name = $1, race = $2, max_hp = $3, is_npc = $4, avatar_url = NULLIF($5, ''), stats = $6, inventory = $7, money = $8
		WHERE id = $9 AND game_id = $10
		RETURNING id, game_id, user_id, name, race, max_hp, current_hp, COALESCE(avatar_url, ''), stats, inventory, is_npc, money, created_at
	`

	char := &model.Character{}
	err := database.DB.QueryRow(context.Background(), query,
		name, race, maxHP, isNPC, avatarURL, stats, inventory, money, id, gameID,
	).Scan(&char.ID, &char.GameID, &char.UserID, &char.Name, &char.Race, &char.MaxHP, &char.CurrentHP, &char.AvatarURL, &char.Stats, &char.Inventory, &char.IsNPC, &char.Money, &char.CreatedAt)

	if err != nil {
		return nil, err
	}

	return char, nil
}

func DeleteCharacter(characterID, gameID string) error {
	query := `DELETE FROM characters WHERE id = $1 AND game_id = $2`
	result, err := database.DB.Exec(context.Background(), query, characterID, gameID)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("character not found")
	}
	return nil
}

func AssignCharacterToPlayer(gameID, characterID, playerID string) error {
	// Verify character exists and belongs to game
	var currentUserID sql.NullString
	queryCheck := `SELECT user_id FROM characters WHERE id = $1 AND game_id = $2`
	err := database.DB.QueryRow(context.Background(), queryCheck, characterID, gameID).Scan(&currentUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("character not found")
		}
		return err
	}

	// Update user_id
	queryUpdate := `UPDATE characters SET user_id = NULLIF($1, '') WHERE id = $2 AND game_id = $3`
	_, err = database.DB.Exec(context.Background(), queryUpdate, playerID, characterID, gameID)
	return err
}
