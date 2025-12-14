package routes

import (
	"questhub/controller"
	"questhub/middleware"

	"github.com/labstack/echo/v4"
)

func initUserRoutes(e *echo.Echo) {
	g := e.Group("/user", middleware.JWTMiddleware)

	g.GET("/stats", controller.GetUserStats)
	g.GET("/campaigns", controller.GetUserCampaigns)
}
