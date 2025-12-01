package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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

func CreateTable(name, gmID string) (*model.Game, error) {
	inviteCode, err := generateInviteCode()
	if err != nil {
		return nil, err
	}

	game := &model.Game{
		Name:       name,
		GmID:       gmID,
		InviteCode: inviteCode,
		IsActive:   true,
		CreatedAt:  time.Now(),
	}

	query := `
		INSERT INTO games (name, gm_id, invite_code, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err = database.DB.QueryRow(context.Background(), query, game.Name, game.GmID, game.InviteCode, game.IsActive, game.CreatedAt).Scan(&game.ID)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func GetTable(id string) (*model.Game, error) {
	game := &model.Game{}
	query := `
		SELECT id, name, gm_id, invite_code, is_active, created_at
		FROM games
		WHERE id = $1
	`
	err := database.DB.QueryRow(context.Background(), query, id).Scan(
		&game.ID, &game.Name, &game.GmID, &game.InviteCode, &game.IsActive, &game.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func GetGames(userID string) ([]model.Game, error) {
	games := []model.Game{}
	query := `
		SELECT g.id, g.name, g.gm_id, u.name, g.invite_code, g.is_active, g.created_at
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
		if err := rows.Scan(&game.ID, &game.Name, &game.GmID, &game.GmName, &game.InviteCode, &game.IsActive, &game.CreatedAt); err != nil {
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
