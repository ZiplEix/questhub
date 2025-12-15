package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"questhub/database"
	model "questhub/models/database"
	"time"

	"github.com/jackc/pgx/v5"
)

func GetUserCharacter(gameID, userID string) (*model.Character, error) {
	char := &model.Character{}
	// Join with game_characters to filtering by game_id and user_id
	query := `
		SELECT c.id, gc.game_id, gc.user_id, c.name, c.avatar_url, c.stats, c.inventory, c.is_npc, c.created_at,
		       c.initiative, c.age, c.height, c.weight, c.max_spells, c.spells, c.abilities, c.experience, c.armor_class, c.speed,
		       c.type, c.sub_race
		FROM characters c
		JOIN game_characters gc ON c.id = gc.character_id
		WHERE gc.game_id = $1 AND gc.user_id = $2
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
		&char.Type,
		&char.SubRace,
	)
	if err != nil {
		if err == sql.ErrNoRows || err == pgx.ErrNoRows {
			return nil, nil // No character found, not an error
		}
		log.Printf("[GetUserCharacter] Scan error: %v\n", err)
		return nil, err
	}
	return char, nil
}

func GetCharacter(gameID, charID string) (*model.Character, error) {
	char := &model.Character{}
	// Join with game_characters to enforce game context and fill game_id/user_id
	query := `
		SELECT c.id, gc.game_id, gc.user_id, c.name, c.race, c.max_hp, c.current_hp, COALESCE(c.avatar_url, ''), c.stats, c.inventory, c.is_npc, c.money, c.created_at,
		       c.initiative, c.age, c.height, c.weight, c.max_spells, c.spells, c.abilities, c.experience, c.type, c.sub_race, c.armor_class, c.speed
		FROM characters c
		JOIN game_characters gc ON c.id = gc.character_id
		WHERE gc.game_id = $1 AND c.id = $2
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
		SELECT c.id, gc.game_id, gc.user_id, c.name, c.race, c.max_hp, c.current_hp, COALESCE(c.avatar_url, ''), c.stats, c.inventory, c.is_npc, c.money, c.created_at, u.name,
		       c.initiative, c.age, c.height, c.weight, c.max_spells, c.spells, c.abilities, c.experience, c.type, c.sub_race, c.armor_class, c.speed
		FROM characters c
		JOIN game_characters gc ON c.id = gc.character_id
		LEFT JOIN "user" u ON gc.user_id = u.id
		WHERE gc.game_id = $1 AND (c.type = 'PLAYER' OR c.type = 'NPC')
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
		SELECT c.id, gc.game_id, gc.user_id, c.name, c.race, c.max_hp, c.current_hp, COALESCE(c.avatar_url, ''), c.stats, c.inventory, c.is_npc, c.money, c.created_at,
		       c.initiative, c.age, c.height, c.weight, c.max_spells, c.spells, c.abilities, c.experience, c.type, c.sub_race, c.armor_class, c.speed
		FROM characters c
		JOIN game_characters gc ON c.id = gc.character_id
		WHERE gc.game_id = $1 AND c.type = 'MONSTER'
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

	// Use Transaction for multiple inserts
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	// 1. Insert into characters (stats/info only)
	queryChar := `
		INSERT INTO characters (name, race, max_hp, current_hp, is_npc, avatar_url, stats, inventory, money, created_at,
		                        initiative, age, height, weight, max_spells, spells, abilities, experience, type, sub_race, armor_class, speed)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), $7, $8, $9, $10,
		        $11, $12, $13, $14, $15, $16, $17, $18, $19, NULLIF($20, ''), $21, $22)
		RETURNING id
	`

	err = tx.QueryRow(context.Background(), queryChar,
		char.Name, char.Race, char.MaxHP, char.CurrentHP, char.IsNPC, avatarURL, string(char.Stats), string(char.Inventory), char.Money, char.CreatedAt,
		char.Initiative, char.Age, char.Height, char.Weight, char.MaxSpells, string(char.Spells), char.Abilities, char.Experience, char.Type, subRace, char.ArmorClass, char.Speed,
	).Scan(&char.ID)

	if err != nil {
		log.Printf("[CreateCharacter] DB Error insert char: %v\n", err)
		return nil, err
	}

	// 2. Insert into game_characters (link)
	queryLink := `
		INSERT INTO game_characters (game_id, character_id, user_id)
		VALUES ($1, $2, NULLIF($3, ''))
	`
	_, err = tx.Exec(context.Background(), queryLink, gameID, char.ID, userID)
	if err != nil {
		log.Printf("[CreateCharacter] DB Error insert link: %v\n", err)
		return nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return char, nil
}

func UpdateCharacter(id, gameID, name, race string, maxHP int, isNPC bool, avatarURL string, stats, inventory []byte, money, initiative int, age, height, weight string, maxSpells int, spells []byte, abilities string, experience int, charType, subRace string, armorClass, speed int) (*model.Character, error) {
	// join with game_characters to enforce game_id check
	query := `
		UPDATE characters c
		SET name = $1, race = $2, max_hp = $3, is_npc = $4, avatar_url = NULLIF($5, ''), stats = $6, inventory = $7, money = $8,
		    initiative = $9, age = $10, height = $11, weight = $12, max_spells = $13, spells = $14, abilities = $15, experience = $16,
			type = $17, sub_race = NULLIF($18, ''), armor_class = $19, speed = $20
		FROM game_characters gc
		WHERE c.id = gc.character_id AND c.id = $21 AND gc.game_id = $22
		RETURNING c.id, gc.game_id, gc.user_id, c.name, c.race, c.max_hp, c.current_hp, COALESCE(c.avatar_url, ''), c.stats, c.inventory, c.is_npc, c.money, c.created_at,
		          c.initiative, c.age, c.height, c.weight, c.max_spells, c.spells, c.abilities, c.experience, c.type, c.sub_race, c.armor_class, c.speed
	`

	char := &model.Character{}
	err := database.DB.QueryRow(context.Background(), query,
		name, race, maxHP, isNPC, avatarURL, string(stats), string(inventory), money,
		initiative, age, height, weight, maxSpells, string(spells), abilities, experience, charType, subRace, armorClass, speed,
		id, gameID,
	).Scan(&char.ID, &char.GameID, &char.UserID, &char.Name, &char.Race, &char.MaxHP, &char.CurrentHP, &char.AvatarURL, &char.Stats, &char.Inventory, &char.IsNPC, &char.Money, &char.CreatedAt,
		&char.Initiative, &char.Age, &char.Height, &char.Weight, &char.MaxSpells, &char.Spells, &char.Abilities, &char.Experience, &char.Type, &char.SubRace, &char.ArmorClass, &char.Speed)

	if err != nil {
		return nil, err
	}

	return char, nil
}

func DeleteCharacter(characterID, gameID string) error {
	// First check/verify it belongs to game via game_characters
	// But actually we can just delete from characters where id in (select character_id from game_characters ...)
	// OR delete from game_characters?
	// If we delete from characters, cascade deletes game_characters.
	// We need to support the use case where multiple games might use same character?
	// User said "dissocier game_id et user_id ... pour que l'on puisse les exporter et importer".
	// Instances are still per game though. "Un joueur ne peut rejoindre qu'une seule fois la même partie".
	// "game_characters" has unique(game_id, character_id).
	// Are characters shared across games? If they are distinct instances (clones via template), then deleting one is fine.
	// If "Templates" are the shared things, then "Characters" table is instances.
	// So deleting a character from "Characters" table is correct.
	// But we must ensure it belongs to the game requesting deletion?

	// We can delete using:
	// DELETE FROM characters WHERE id = $1 AND id IN (SELECT character_id FROM game_characters WHERE game_id = $2)

	query := `DELETE FROM characters WHERE id = $1 AND id IN (SELECT character_id FROM game_characters WHERE game_id = $2)`
	result, err := database.DB.Exec(context.Background(), query, characterID, gameID)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("character not found or not in game")
	}
	return nil
}

func AssignCharacterToPlayer(gameID, characterID, playerID string) error {
	// Verify character exists and belongs to game in game_characters
	// And update user_id in game_characters

	queryUpdate := `UPDATE game_characters SET user_id = NULLIF($1, '') WHERE character_id = $2 AND game_id = $3`
	result, err := database.DB.Exec(context.Background(), queryUpdate, playerID, characterID, gameID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("character not found in game")
	}
	return nil
}

func GetNotes(gameID, userID string) (string, error) {
	var content string
	query := `SELECT content FROM notes WHERE game_id = $1 AND user_id = $2`
	err := database.DB.QueryRow(context.Background(), query, gameID, userID).Scan(&content)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No notes found, return empty string
		}
		return "", err
	}
	return content, nil
}

func UpdateNotes(gameID, userID, content string) error {
	query := `
		INSERT INTO notes (game_id, user_id, content)
		VALUES ($1, $2, $3)
		ON CONFLICT (game_id, user_id) DO UPDATE
		SET content = EXCLUDED.content
	`
	_, err := database.DB.Exec(context.Background(), query, gameID, userID, content)
	return err
}

func EnsureGMCharacter(gameID, manualGMID string) (*model.Character, error) {
	log.Printf("[EnsureGMCharacter] Checking for GM char. GameID=%s, UserID=%s\n", gameID, manualGMID)
	// Check if character already exists for GM
	char, err := GetUserCharacter(gameID, manualGMID)
	if err != nil {
		log.Printf("[EnsureGMCharacter] Error getting user character: %v\n", err)
		return nil, err
	}
	if char != nil {
		log.Printf("[EnsureGMCharacter] Found existing character: %s\n", char.ID)
		return char, nil
	}

	log.Printf("[EnsureGMCharacter] Creating new hidden GM character...\n")

	// Create hidden character for GM
	// We use "GM_HIDDEN" as type to easily filter it out if needed, though GetGameCharacters filters by PLAYER/NPC usually.
	newChar, err := CreateCharacter(
		gameID,
		manualGMID,
		"GM Notes",
		"Game Master",
		1,            // HP
		false,        // isNPC (technically no, but...)
		"",           // avatar
		[]byte("{}"), // stats (NOT NULL constraint)
		[]byte("[]"), // inventory (NOT NULL constraint)
		0,            // money
		0,            // initiative
		"",           // age
		"",           // height
		"",           // weight
		0,            // maxSpells
		[]byte("{}"), // spells
		"",           // abilities
		0,            // experience
		"GM_HIDDEN",  // type
		"",           // subRace
		10,           // armorClass
		30,           // speed
	)
	if err != nil {
		log.Printf("[EnsureGMCharacter] Error creating character: %v\n", err)
		return nil, err
	}
	log.Printf("[EnsureGMCharacter] Created character: %s\n", newChar.ID)
	return newChar, nil
}
