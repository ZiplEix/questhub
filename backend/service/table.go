package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"questhub/database"
	model "questhub/models/database"
)

func generateInviteCode() (string, error) {
	bytes := make([]byte, 10) // 10 bytes = 20 hex chars
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

func CreateTable(name, gmID, imageURL string) (*model.Game, error) {
	inviteCode, err := generateInviteCode()
	if err != nil {
		return nil, err
	}

	game := &model.Game{
		Name:       name,
		GmID:       gmID,
		InviteCode: inviteCode,
		IsActive:   true,
		ImageURL:   imageURL,
		CreatedAt:  time.Now(),
	}

	query := `
		INSERT INTO games (name, gm_id, invite_code, is_active, image_url, created_at)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6)
		RETURNING id
	`

	err = database.DB.QueryRow(context.Background(), query, game.Name, game.GmID, game.InviteCode, game.IsActive, game.ImageURL, game.CreatedAt).Scan(&game.ID)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func GetTable(id string) (*model.Game, error) {
	game := &model.Game{}
	query := `
		SELECT id, name, gm_id, invite_code, is_active, COALESCE(image_url, ''), created_at
		FROM games
		WHERE id = $1
	`
	err := database.DB.QueryRow(context.Background(), query, id).Scan(
		&game.ID, &game.Name, &game.GmID, &game.InviteCode, &game.IsActive, &game.ImageURL, &game.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func GetGames(userID string) ([]model.Game, error) {
	games := []model.Game{}
	query := `
		SELECT g.id, g.name, g.gm_id, u.name, g.invite_code, g.is_active, COALESCE(g.image_url, ''), g.created_at
		FROM games g
		JOIN "user" u ON g.gm_id = u.id
		WHERE g.gm_id = $1
		ORDER BY g.created_at DESC
	`
	rows, err := database.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var game model.Game
		if err := rows.Scan(&game.ID, &game.Name, &game.GmID, &game.GmName, &game.InviteCode, &game.IsActive, &game.ImageURL, &game.CreatedAt); err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

func JoinTable(inviteCode, userID string) (string, error) {
	var gameID string
	// Get game ID from invite code
	err := database.DB.QueryRow(context.Background(), "SELECT id FROM games WHERE invite_code = $1", inviteCode).Scan(&gameID)
	if err != nil {
		return "", err
	}

	// Insert into game_players
	query := `
		INSERT INTO game_players (game_id, user_id)
		VALUES ($1, $2)
	`
	_, err = database.DB.Exec(context.Background(), query, gameID, userID)
	if err != nil {
		return gameID, err
	}

	return gameID, nil
}

func DeleteTable(id, userID string) error {
	// 1. Get Game to check GM and get ImageURL
	game, err := GetTable(id)
	if err != nil {
		return err
	}

	if game.GmID != userID {
		return errors.New("unauthorized: only the GM can delete the table")
	}

	// 2. Get all Characters for the game to get AvatarURLs
	// We need to query this manually as we don't have a GetCharacters service handy here yet
	// or we can just query the avatar_urls directly
	rows, err := database.DB.Query(context.Background(), "SELECT avatar_url FROM characters WHERE game_id = $1", id)
	if err != nil {
		return err
	}
	defer rows.Close()

	var filesToDelete []string
	if game.ImageURL != "" {
		// Extract filename from URL
		// Assuming URL format: http://host/uploads/filename
		parts := strings.Split(game.ImageURL, "/uploads/")
		if len(parts) == 2 {
			filesToDelete = append(filesToDelete, filepath.Join("uploads", parts[1]))
		}
	}

	for rows.Next() {
		var avatarURL sql.NullString
		if err := rows.Scan(&avatarURL); err != nil {
			continue
		}
		if avatarURL.Valid && avatarURL.String != "" {
			parts := strings.Split(avatarURL.String, "/uploads/")
			if len(parts) == 2 {
				filesToDelete = append(filesToDelete, filepath.Join("uploads", parts[1]))
			}
		}
	}

	// 3. Delete Game from DB (Cascade will handle characters and players)
	_, err = database.DB.Exec(context.Background(), "DELETE FROM games WHERE id = $1", id)
	if err != nil {
		return err
	}

	// 4. Delete files from uploads/ if they exist
	for _, path := range filesToDelete {
		if err := os.Remove(path); err != nil {
			// Log error but continue? Or just ignore if file doesn't exist
			if !os.IsNotExist(err) {
				// In a real app, we might want to log this
				fmt.Printf("Failed to delete file %s: %v\n", path, err)
			}
		}
	}

	return nil
}
