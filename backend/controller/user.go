package controller

import (
	"net/http"
	"questhub/service"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func GetUserStats(c echo.Context) error {
	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	stats, err := service.GetUserStats(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch stats").SetInternal(err)
	}

	return c.JSON(http.StatusOK, stats)
}

func GetUserCampaigns(c echo.Context) error {
	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	campaigns, err := service.GetUserCampaigns(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch campaigns").SetInternal(err)
	}

	return c.JSON(http.StatusOK, campaigns)
}
