package controller

import (
	"fmt"
	"net/http"
	model "questhub/models/database"
	"questhub/models/request"
	"questhub/service"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func CreateTable(c echo.Context) error {
	var req request.CreateTableRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name is required")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	gmID := claims["sub"].(string)

	game, err := service.CreateTable(req.Name, gmID, req.ImageURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create table").SetInternal(err)
	}

	return c.JSON(http.StatusCreated, game)
}

type GameResponse struct {
	*model.Game
	CurrentCharacterID *string `json:"current_character_id"`
	IsGM               bool    `json:"is_gm"`
}

func GetTable(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	game, err := service.GetTable(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}

	response := GameResponse{
		Game: game,
		IsGM: game.GmID == userID,
	}

	// Fetch user character
	character, err := service.GetUserCharacter(id, userID)
	if err != nil {
		// Log error but don't fail the request?
		// Or fail? Let's just log and continue for now, assuming no character
		fmt.Printf("Failed to fetch user character: %v\n", err)
	} else if character != nil {
		response.CurrentCharacterID = &character.ID
	}

	return c.JSON(http.StatusOK, response)
}

func GetGames(c echo.Context) error {
	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	games, err := service.GetGames(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch games").SetInternal(err)
	}

	return c.JSON(http.StatusOK, games)
}

func JoinTable(c echo.Context) error {
	var req struct {
		InviteCode string `json:"invite_code"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.InviteCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invite code is required")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	gameID, err := service.JoinTable(req.InviteCode, userID)
	if err != nil {
		// Check for unique constraint violation (user already in game or invited)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.JSON(http.StatusConflict, map[string]string{
				"message": "User already in game or invitation pending",
				"id":      gameID,
			})
		}
		return echo.NewHTTPError(http.StatusNotFound, "Game not found or failed to join")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Invitation sent", "id": gameID})
}

func GetPendingInvitations(c echo.Context) error {
	gameID := c.Param("id")
	if gameID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	// Verify GM
	game, err := service.GetTable(gameID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}
	if game.GmID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Only the GM can view invitations")
	}

	invitations, err := service.GetPendingInvitations(gameID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch invitations").SetInternal(err)
	}

	return c.JSON(http.StatusOK, invitations)
}

func AcceptInvitation(c echo.Context) error {
	gameID := c.Param("id")
	targetUserID := c.Param("userId")

	if gameID == "" || targetUserID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or user ID")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	// Verify GM
	game, err := service.GetTable(gameID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}
	if game.GmID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Only the GM can accept invitations")
	}

	if err := service.AcceptInvitation(gameID, targetUserID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to accept invitation").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Invitation accepted"})
}

func DeclineInvitation(c echo.Context) error {
	gameID := c.Param("id")
	targetUserID := c.Param("userId")

	if gameID == "" || targetUserID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or user ID")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	// Verify GM
	game, err := service.GetTable(gameID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}
	if game.GmID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Only the GM can decline invitations")
	}

	if err := service.DeclineInvitation(gameID, targetUserID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decline invitation").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Invitation declined"})
}

func DeleteTable(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	if err := service.DeleteTable(id, userID); err != nil {
		if err.Error() == "unauthorized: only the GM can delete the table" {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete table").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Table deleted successfully"})
}

func RegenerateInviteCode(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	// Verify GM
	game, err := service.GetTable(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}
	if game.GmID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Only the GM can regenerate the invite code")
	}

	newCode, err := service.RegenerateInviteCode(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to regenerate invite code").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"invite_code": newCode})
}

func GetGamePlayers(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	// In a real app, we might want to check if the user is part of the game or is the GM
	// For now, we'll allow any authenticated user to see the players of a game they know the ID of
	// Or we can restrict it to GM only if this is strictly for the settings page.
	// The requirement says "GM settings page", so let's restrict to GM for now to be safe,
	// but usually players can see other players. Let's stick to GM check as it's for the settings page.

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	game, err := service.GetTable(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}

	if game.GmID != userID {
		// Check if user is a player? For now let's just restrict to GM as requested context is GM settings
		return echo.NewHTTPError(http.StatusForbidden, "Only the GM can view players in settings")
	}

	players, err := service.GetGamePlayers(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch players").SetInternal(err)
	}

	return c.JSON(http.StatusOK, players)
}

func RemovePlayer(c echo.Context) error {
	gameID := c.Param("id")
	playerID := c.Param("userId")

	if gameID == "" || playerID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or player ID")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	// Verify GM
	game, err := service.GetTable(gameID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}
	if game.GmID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Only the GM can remove players")
	}

	// Prevent removing GM
	if playerID == game.GmID {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot remove the Game Master")
	}

	err = service.RemovePlayer(gameID, playerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove player").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Player removed successfully"})
}
