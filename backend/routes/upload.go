package routes

import (
	"questhub/controller"
	"questhub/middleware"

	"github.com/labstack/echo/v4"
)

func initUploadRoutes(e *echo.Echo) {
	g := e.Group("/upload", middleware.JWTMiddleware)

	g.POST("/image", controller.UploadImage)
}
