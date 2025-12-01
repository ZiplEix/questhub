package controller

import (
	"net/http"
	"questhub/models/request"
	"questhub/service"

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
