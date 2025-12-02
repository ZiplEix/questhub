package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"questhub/database"
	model "questhub/models/database"
	"time"

	"github.com/google/uuid"
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

func GetGameCharacters(gameID string) ([]model.Character, error) {
	characters := []model.Character{}
	query := `
		SELECT c.id, c.game_id, c.user_id, c.name, c.race, c.max_hp, c.current_hp, c.avatar_url, c.stats, c.inventory, c.is_npc, c.created_at, u.name
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
		var character model.Character
		var playerName sql.NullString
		if err := rows.Scan(
			&character.ID,
			&character.GameID,
			&character.UserID,
			&character.Name,
			&character.Race,
			&character.MaxHP,
			&character.CurrentHP,
			&character.AvatarURL,
			&character.Stats,
			&character.Inventory,
			&character.IsNPC,
			&character.CreatedAt,
			&playerName,
		); err != nil {
			return nil, err
		}
		if playerName.Valid {
			character.PlayerName = &playerName.String
		}
		characters = append(characters, character)
	}

	return characters, nil
}

type InventoryItem struct {
	Name     string  `json:"name"`
	Quantity string  `json:"quantity"`
	ImageURL *string `json:"image_url,omitempty"`
	IconName string  `json:"icon_name,omitempty"`
}

func CreateCharacter(gameID, name, race string, maxHP int, isNPC bool, avatarFile *multipart.FileHeader, avatarURL string, stats json.RawMessage, inventoryItems []InventoryItem, inventoryImages map[int]*multipart.FileHeader) (*model.Character, error) {
	var finalAvatarURL *string

	// Create uploads directory if it doesn't exist
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// Handle Avatar
	if avatarFile != nil {
		src, err := avatarFile.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		ext := filepath.Ext(avatarFile.Filename)
		filename := fmt.Sprintf("%d_avatar%s", time.Now().UnixNano(), ext)
		filepath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(filepath)
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("/uploads/%s", filename)
		finalAvatarURL = &url
	} else if avatarURL != "" {
		finalAvatarURL = &avatarURL
	}

	// Handle Inventory Images
	for i := range inventoryItems {
		if file, ok := inventoryImages[i]; ok {
			src, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer src.Close()

			ext := filepath.Ext(file.Filename)
			filename := fmt.Sprintf("%d_inv_%d%s", time.Now().UnixNano(), i, ext)
			filepath := filepath.Join(uploadDir, filename)

			dst, err := os.Create(filepath)
			if err != nil {
				return nil, err
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				return nil, err
			}

			url := fmt.Sprintf("/uploads/%s", filename)
			inventoryItems[i].ImageURL = &url
		}
	}

	inventoryJSON, err := json.Marshal(inventoryItems)
	if err != nil {
		return nil, err
	}

	// Insert into Database
	character := model.Character{
		ID:        uuid.New().String(),
		GameID:    gameID,
		Name:      name,
		Race:      race,
		MaxHP:     maxHP,
		CurrentHP: maxHP, // Start with full HP
		AvatarURL: finalAvatarURL,
		Stats:     stats,
		Inventory: inventoryJSON,
		IsNPC:     isNPC,
		CreatedAt: time.Now(),
	}

	query := `
		INSERT INTO characters (id, game_id, name, race, max_hp, current_hp, avatar_url, stats, inventory, is_npc, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = database.DB.Exec(context.Background(), query,
		character.ID,
		character.GameID,
		character.Name,
		character.Race,
		character.MaxHP,
		character.CurrentHP,
		character.AvatarURL,
		character.Stats,
		character.Inventory,
		character.IsNPC,
		character.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &character, nil
}

func UpdateCharacter(characterID, gameID, name, race string, maxHP int, isNPC bool, avatarFile *multipart.FileHeader, avatarURL string, stats json.RawMessage, inventoryItems []InventoryItem, inventoryImages map[int]*multipart.FileHeader) (*model.Character, error) {
	var currentAvatarURL *string
	var currentInventoryJSON []byte

	queryGet := `SELECT avatar_url, inventory FROM characters WHERE id = $1 AND game_id = $2`
	err := database.DB.QueryRow(context.Background(), queryGet, characterID, gameID).Scan(&currentAvatarURL, &currentInventoryJSON)
	if err != nil {
		return nil, err
	}

	var finalAvatarURL *string = currentAvatarURL

	// Create uploads directory if it doesn't exist
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// Handle Avatar Update
	if avatarFile != nil {
		src, err := avatarFile.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		ext := filepath.Ext(avatarFile.Filename)
		filename := fmt.Sprintf("%d_avatar%s", time.Now().UnixNano(), ext)
		filepath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(filepath)
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("/uploads/%s", filename)
		finalAvatarURL = &url
	} else if avatarURL != "" {
		finalAvatarURL = &avatarURL
	}

	// Handle Inventory Images
	// We need to preserve existing images if not replaced.
	// The inventoryItems passed here are the NEW state.
	// If an item has an image_url already set (from frontend), we keep it.
	// If it has a new file in inventoryImages, we upload and set it.

	for i := range inventoryItems {
		if file, ok := inventoryImages[i]; ok {
			src, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer src.Close()

			ext := filepath.Ext(file.Filename)
			filename := fmt.Sprintf("%d_inv_%d%s", time.Now().UnixNano(), i, ext)
			filepath := filepath.Join(uploadDir, filename)

			dst, err := os.Create(filepath)
			if err != nil {
				return nil, err
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				return nil, err
			}

			url := fmt.Sprintf("/uploads/%s", filename)
			inventoryItems[i].ImageURL = &url
		}
	}

	inventoryJSON, err := json.Marshal(inventoryItems)
	if err != nil {
		return nil, err
	}

	// Update Database
	queryUpdate := `
		UPDATE characters
		SET name = $1, race = $2, max_hp = $3, is_npc = $4, avatar_url = $5, stats = $6, inventory = $7
		WHERE id = $8 AND game_id = $9
		RETURNING id, game_id, user_id, name, race, max_hp, current_hp, avatar_url, stats, inventory, is_npc, created_at
	`

	var updatedChar model.Character
	err = database.DB.QueryRow(context.Background(), queryUpdate,
		name, race, maxHP, isNPC, finalAvatarURL, stats, inventoryJSON, characterID, gameID,
	).Scan(
		&updatedChar.ID,
		&updatedChar.GameID,
		&updatedChar.UserID,
		&updatedChar.Name,
		&updatedChar.Race,
		&updatedChar.MaxHP,
		&updatedChar.CurrentHP,
		&updatedChar.AvatarURL,
		&updatedChar.Stats,
		&updatedChar.Inventory,
		&updatedChar.IsNPC,
		&updatedChar.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updatedChar, nil
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
	queryUpdate := `UPDATE characters SET user_id = $1 WHERE id = $2 AND game_id = $3`
	_, err = database.DB.Exec(context.Background(), queryUpdate, playerID, characterID, gameID)
	return err
}
