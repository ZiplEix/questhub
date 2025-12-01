package controller

import (
	"net/http"
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

	game, err := service.CreateTable(req.Name, gmID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create table").SetInternal(err)
	}

	return c.JSON(http.StatusCreated, game)
}

func GetTable(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	game, err := service.GetTable(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}

	return c.JSON(http.StatusOK, game)
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
		// Check for unique constraint violation (user already in game)
		// This is a simple string check, ideally use pgconn.PgError
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			// Return 409 with game ID so frontend can redirect
			return c.JSON(http.StatusConflict, map[string]string{
				"message": "User already in game",
				"id":      gameID,
			})
		}
		return echo.NewHTTPError(http.StatusNotFound, "Game not found or failed to join")
	}

	return c.JSON(http.StatusOK, map[string]string{"id": gameID})
}
