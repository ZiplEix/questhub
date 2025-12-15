package service

import (
	"context"
	"fmt"
	"questhub/database"
	"time"
)

type UserStats struct {
	GamesPlayed       int64 `json:"games_played"`
	GamesMastered     int64 `json:"games_mastered"`
	CharactersCreated int64 `json:"characters_created"`
}

type UserCampaign struct {
	GameID          string    `json:"game_id"`
	GameName        string    `json:"game_name"`
	GameImageURL    string    `json:"game_image_url"`
	CharacterName   string    `json:"character_name"`
	CharacterAvatar string    `json:"character_avatar_url"`
	JoinedAt        time.Time `json:"joined_at"`
}

func GetUserStats(userID string) (*UserStats, error) {
	stats := &UserStats{}

	// Count games played
	err := database.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM game_players WHERE user_id = $1
	`, userID).Scan(&stats.GamesPlayed)
	if err != nil {
		return nil, fmt.Errorf("failed to count games played: %w", err)
	}

	// Count games mastered
	err = database.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM games WHERE gm_id = $1
	`, userID).Scan(&stats.GamesMastered)
	if err != nil {
		return nil, fmt.Errorf("failed to count games mastered: %w", err)
	}

	// Count characters created
	err = database.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM game_characters WHERE user_id = $1
	`, userID).Scan(&stats.CharactersCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to count characters: %w", err)
	}

	return stats, nil
}

func GetUserCampaigns(userID string) ([]UserCampaign, error) {
	rows, err := database.DB.Query(context.Background(), `
		SELECT 
			g.id, g.name, COALESCE(g.image_url, ''), 
			c.name, COALESCE(c.avatar_url, ''), gc.assigned_at
		FROM characters c
		JOIN game_characters gc ON c.id = gc.character_id
		JOIN games g ON g.id = gc.game_id
		WHERE gc.user_id = $1
		ORDER BY gc.assigned_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []UserCampaign
	for rows.Next() {
		var c UserCampaign
		// We need to handle potential NULLs if we change the query, but here inner join ensures game exists.
		// However, image_url and avatar_url can be null in DB, handled by COALESCE in SQL.
		// But in Go, we need to map to string. The SQL COALESCE handles it.
		if err := rows.Scan(&c.GameID, &c.GameName, &c.GameImageURL, &c.CharacterName, &c.CharacterAvatar, &c.JoinedAt); err != nil {
			return nil, fmt.Errorf("failed to scan campaign row: %w", err)
		}
		campaigns = append(campaigns, c)
	}

	return campaigns, nil
}
