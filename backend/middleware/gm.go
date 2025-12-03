package middleware

import (
	"net/http"
	"questhub/service"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func RequireGM(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		gameID := c.Param("id")
		if gameID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
		}

		claims := c.Get("claims").(jwt.MapClaims)
		userID := claims["sub"].(string)

		game, err := service.GetTable(gameID)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Game not found")
		}

		if game.GmID != userID {
			return echo.NewHTTPError(http.StatusForbidden, "Only the GM can perform this action")
		}

		return next(c)
	}
}
