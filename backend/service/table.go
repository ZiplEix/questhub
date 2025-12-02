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
	bytes := make([]byte, 5) // 5 bytes = 10 hex chars
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
		SELECT DISTINCT g.id, g.name, g.gm_id, u.name, g.invite_code, g.is_active, COALESCE(g.image_url, ''), g.created_at
		FROM games g
		JOIN "user" u ON g.gm_id = u.id
		LEFT JOIN game_players gp ON g.id = gp.game_id
		WHERE g.gm_id = $1 OR gp.user_id = $1
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

	// Insert into game_invitations
	query := `
		INSERT INTO game_invitations (game_id, user_id)
		VALUES ($1, $2)
	`
	_, err = database.DB.Exec(context.Background(), query, gameID, userID)
	if err != nil {
		return gameID, err
	}

	return gameID, nil
}

func GetPendingInvitations(gameID string) ([]model.Invitation, error) {
	invitations := []model.Invitation{}
	query := `
		SELECT i.id, i.game_id, i.user_id, u.name, i.created_at
		FROM game_invitations i
		JOIN "user" u ON i.user_id = u.id
		WHERE i.game_id = $1
		ORDER BY i.created_at ASC
	`
	rows, err := database.DB.Query(context.Background(), query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invitation model.Invitation
		if err := rows.Scan(&invitation.ID, &invitation.GameID, &invitation.UserID, &invitation.UserName, &invitation.CreatedAt); err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func AcceptInvitation(gameID, userID string) error {
	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert into game_players
	_, err = tx.Exec(ctx, "INSERT INTO game_players (game_id, user_id) VALUES ($1, $2)", gameID, userID)
	if err != nil {
		return err
	}

	// 2. Delete from game_invitations
	_, err = tx.Exec(ctx, "DELETE FROM game_invitations WHERE game_id = $1 AND user_id = $2", gameID, userID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func DeclineInvitation(gameID, userID string) error {
	_, err := database.DB.Exec(context.Background(), "DELETE FROM game_invitations WHERE game_id = $1 AND user_id = $2", gameID, userID)
	return err
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

func RegenerateInviteCode(gameID string) (string, error) {
	newCode, err := generateInviteCode()
	if err != nil {
		return "", err
	}

	query := `
		UPDATE games
		SET invite_code = $1
		WHERE id = $2
	`
	_, err = database.DB.Exec(context.Background(), query, newCode, gameID)
	if err != nil {
		return "", err
	}

	return newCode, nil
}

func GetGamePlayers(gameID string) ([]model.Player, error) {
	// 1. Get Game to know who is GM
	game, err := GetTable(gameID)
	if err != nil {
		return nil, err
	}

	players := []model.Player{}
	query := `
		SELECT u.id, u.name, COALESCE(u.image, ''), gp.joined_at
		FROM game_players gp
		JOIN "user" u ON gp.user_id = u.id
		WHERE gp.game_id = $1
		ORDER BY gp.joined_at ASC
	`
	rows, err := database.DB.Query(context.Background(), query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var player model.Player
		if err := rows.Scan(&player.UserID, &player.Name, &player.AvatarURL, &player.JoinedAt); err != nil {
			return nil, err
		}
		if player.UserID == game.GmID {
			player.IsGM = true
		}
		players = append(players, player)
	}

	// Check if GM is in the list (usually not in game_players table if they created it but didn't "join" as player)
	// If GM is not in list, we should probably add them?
	// The requirement is "display players present in the game". Usually GM is considered present.
	// Let's check if GM is already in the list.
	gmInList := false
	for _, p := range players {
		if p.UserID == game.GmID {
			gmInList = true
			break
		}
	}

	if !gmInList {
		// Fetch GM details
		var gm model.Player
		gmQuery := `SELECT id, name, COALESCE(image, '') FROM "user" WHERE id = $1`
		err := database.DB.QueryRow(context.Background(), gmQuery, game.GmID).Scan(&gm.UserID, &gm.Name, &gm.AvatarURL)
		if err == nil {
			gm.IsGM = true
			gm.JoinedAt = game.CreatedAt // GM joined when game was created
			// Prepend GM to list
			players = append([]model.Player{gm}, players...)
		}
	}

	return players, nil
}

func RemovePlayer(gameID, userID string) error {
	_, err := database.DB.Exec(context.Background(), "DELETE FROM game_players WHERE game_id = $1 AND user_id = $2", gameID, userID)
	return err
}
