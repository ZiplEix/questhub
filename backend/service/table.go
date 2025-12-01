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
	bytes := make([]byte, 3) // 3 bytes = 6 hex chars
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
