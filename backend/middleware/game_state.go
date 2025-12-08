package middleware

import (
	"net/http"
	"questhub/service"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func CheckGameState(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. Get Game ID
		gameID := c.Param("id")
		if gameID == "" {
			// If no ID, skip check
			return next(c)
		}

		// 2. Get User ID
		user := c.Get("claims").(jwt.MapClaims)
		userID := user["sub"].(string)

		// 3. Get Game
		game, err := service.GetTable(gameID)
		if err != nil {
			// If game not found, let controller handle 404 or return here
			return echo.NewHTTPError(http.StatusNotFound, "Game not found")
		}

		// 4. Logic
		// GM Bypass
		if game.GmID == userID {
			return next(c)
		}

		// GET Method Bypass (Read-only)
		if c.Request().Method == http.MethodGet {
			return next(c)
		}

		// Check Pause State
		if game.State == "paused" {
			return echo.NewHTTPError(http.StatusForbidden, "Game is paused. Actions are restricted.")
		}

		return next(c)
	}
}
