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
		SELECT id, game_id, user_id, name, avatar_url, stats, inventory, is_npc, created_at,
		       initiative, age, height, weight, max_spells, spells, abilities, experience, armor_class, speed
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
		&char.Initiative,
		&char.Age,
		&char.Height,
		&char.Weight,
		&char.MaxSpells,
		&char.Spells,
		&char.Abilities,
		&char.Experience,
		&char.ArmorClass,
		&char.Speed,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No character found, not an error
		}
		return nil, err
	}
	return char, nil
}

func GetCharacter(gameID, charID string) (*model.Character, error) {
	char := &model.Character{}
	query := `
		SELECT id, game_id, user_id, name, race, max_hp, current_hp, COALESCE(avatar_url, ''), stats, inventory, is_npc, money, created_at,
		       initiative, age, height, weight, max_spells, spells, abilities, experience, type, sub_race, armor_class, speed
		FROM characters
		WHERE game_id = $1 AND id = $2
	`
	err := database.DB.QueryRow(context.Background(), query, gameID, charID).Scan(
		&char.ID,
		&char.GameID,
		&char.UserID,
		&char.Name,
		&char.Race,
		&char.MaxHP,
		&char.CurrentHP,
		&char.AvatarURL,
		&char.Stats,
		&char.Inventory,
		&char.IsNPC,
		&char.Money,
		&char.CreatedAt,
		&char.Initiative,
		&char.Age,
		&char.Height,
		&char.Weight,
		&char.MaxSpells,
		&char.Spells,
		&char.Abilities,
		&char.Experience,
		&char.Type,
		&char.SubRace,
		&char.ArmorClass,
		&char.Speed,
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
		SELECT c.id, c.game_id, c.user_id, c.name, c.race, c.max_hp, c.current_hp, COALESCE(c.avatar_url, ''), c.stats, c.inventory, c.is_npc, c.money, c.created_at, u.name,
		       c.initiative, c.age, c.height, c.weight, c.max_spells, c.spells, c.abilities, c.experience, c.type, c.sub_race, c.armor_class, c.speed
		FROM characters c
		LEFT JOIN "user" u ON c.user_id = u.id
		WHERE c.game_id = $1 AND (c.type = 'PLAYER' OR c.type = 'NPC')
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
		if err := rows.Scan(&char.ID, &char.GameID, &char.UserID, &char.Name, &char.Race, &char.MaxHP, &char.CurrentHP, &char.AvatarURL, &char.Stats, &char.Inventory, &char.IsNPC, &char.Money, &char.CreatedAt, &playerName,
			&char.Initiative, &char.Age, &char.Height, &char.Weight, &char.MaxSpells, &char.Spells, &char.Abilities, &char.Experience, &char.Type, &char.SubRace, &char.ArmorClass, &char.Speed); err != nil {
			return nil, err
		}
		if playerName.Valid {
			char.PlayerName = &playerName.String
		}
		characters = append(characters, char)
	}

	return characters, nil
}

func GetGameMonsters(gameID string) ([]model.Character, error) {
	characters := []model.Character{}
	query := `
		SELECT c.id, c.game_id, c.user_id, c.name, c.race, c.max_hp, c.current_hp, COALESCE(c.avatar_url, ''), c.stats, c.inventory, c.is_npc, c.money, c.created_at,
		       c.initiative, c.age, c.height, c.weight, c.max_spells, c.spells, c.abilities, c.experience, c.type, c.sub_race, c.armor_class, c.speed
		FROM characters c
		WHERE c.game_id = $1 AND c.type = 'MONSTER'
		ORDER BY c.created_at DESC
	`
	rows, err := database.DB.Query(context.Background(), query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var char model.Character
		if err := rows.Scan(&char.ID, &char.GameID, &char.UserID, &char.Name, &char.Race, &char.MaxHP, &char.CurrentHP, &char.AvatarURL, &char.Stats, &char.Inventory, &char.IsNPC, &char.Money, &char.CreatedAt,
			&char.Initiative, &char.Age, &char.Height, &char.Weight, &char.MaxSpells, &char.Spells, &char.Abilities, &char.Experience, &char.Type, &char.SubRace, &char.ArmorClass, &char.Speed); err != nil {
			return nil, err
		}
		characters = append(characters, char)
	}

	return characters, nil
}

type InventoryItem struct {
	Name        string  `json:"name"`
	Quantity    string  `json:"quantity"`
	Description string  `json:"description,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	IconName    string  `json:"icon_name,omitempty"`
}

func CreateCharacter(gameID, userID, name, race string, maxHP int, isNPC bool, avatarURL string, stats, inventory []byte, money, initiative int, age, height, weight string, maxSpells int, spells []byte, abilities string, experience int, charType, subRace string, armorClass, speed int) (*model.Character, error) {
	char := &model.Character{
		GameID:     gameID,
		Name:       name,
		Race:       race,
		MaxHP:      maxHP,
		CurrentHP:  maxHP,
		IsNPC:      isNPC,
		Stats:      stats,
		Inventory:  inventory,
		Money:      money,
		CreatedAt:  time.Now(),
		Initiative: initiative,
		Age:        age,
		Height:     height,
		Weight:     weight,
		MaxSpells:  maxSpells,
		Spells:     spells,
		Abilities:  abilities,
		Experience: experience,
		Type:       charType,
		ArmorClass: armorClass,
		Speed:      speed,
	}

	if userID != "" {
		char.UserID = &userID
	}
	if avatarURL != "" {
		char.AvatarURL = &avatarURL
	}
	if subRace != "" {
		char.SubRace = &subRace
	}

	query := `
		INSERT INTO characters (game_id, user_id, name, race, max_hp, current_hp, is_npc, avatar_url, stats, inventory, money, created_at,
		                        initiative, age, height, weight, max_spells, spells, abilities, experience, type, sub_race, armor_class, speed)
		VALUES ($1, NULLIF($2, ''), $3, $4, $5, $6, $7, NULLIF($8, ''), $9, $10, $11, $12,
		        $13, $14, $15, $16, $17, $18, $19, $20, $21, NULLIF($22, ''), $23, $24)
		RETURNING id
	`

	err := database.DB.QueryRow(context.Background(), query,
		char.GameID, userID, char.Name, char.Race, char.MaxHP, char.CurrentHP, char.IsNPC, avatarURL, char.Stats, char.Inventory, char.Money, char.CreatedAt,
		char.Initiative, char.Age, char.Height, char.Weight, char.MaxSpells, char.Spells, char.Abilities, char.Experience, char.Type, subRace, char.ArmorClass, char.Speed,
	).Scan(&char.ID)

	if err != nil {
		return nil, err
	}

	return char, nil
}

func UpdateCharacter(id, gameID, name, race string, maxHP int, isNPC bool, avatarURL string, stats, inventory []byte, money, initiative int, age, height, weight string, maxSpells int, spells []byte, abilities string, experience int, charType, subRace string, armorClass, speed int) (*model.Character, error) {
	query := `
		UPDATE characters
		SET name = $1, race = $2, max_hp = $3, is_npc = $4, avatar_url = NULLIF($5, ''), stats = $6, inventory = $7, money = $8,
		    initiative = $9, age = $10, height = $11, weight = $12, max_spells = $13, spells = $14, abilities = $15, experience = $16,
			type = $17, sub_race = NULLIF($18, ''), armor_class = $19, speed = $20
		WHERE id = $21 AND game_id = $22
		RETURNING id, game_id, user_id, name, race, max_hp, current_hp, COALESCE(avatar_url, ''), stats, inventory, is_npc, money, created_at,
		          initiative, age, height, weight, max_spells, spells, abilities, experience, type, sub_race, armor_class, speed
	`

	char := &model.Character{}
	err := database.DB.QueryRow(context.Background(), query,
		name, race, maxHP, isNPC, avatarURL, stats, inventory, money,
		initiative, age, height, weight, maxSpells, spells, abilities, experience, charType, subRace, armorClass, speed,
		id, gameID,
	).Scan(&char.ID, &char.GameID, &char.UserID, &char.Name, &char.Race, &char.MaxHP, &char.CurrentHP, &char.AvatarURL, &char.Stats, &char.Inventory, &char.IsNPC, &char.Money, &char.CreatedAt,
		&char.Initiative, &char.Age, &char.Height, &char.Weight, &char.MaxSpells, &char.Spells, &char.Abilities, &char.Experience, &char.Type, &char.SubRace, &char.ArmorClass, &char.Speed)

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

func GetCharacterNotes(charID string) (string, error) {
	var content string
	query := `SELECT content FROM character_notes WHERE character_id = $1`
	err := database.DB.QueryRow(context.Background(), query, charID).Scan(&content)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No notes found, return empty string
		}
		return "", err
	}
	return content, nil
}

func UpdateCharacterNotes(charID, content string) error {
	query := `
		INSERT INTO character_notes (character_id, content)
		VALUES ($1, $2)
		ON CONFLICT (character_id) DO UPDATE
		SET content = EXCLUDED.content
	`
	_, err := database.DB.Exec(context.Background(), query, charID, content)
	return err
}
