package routes

import (
	"questhub/controller"
	"questhub/middleware"

	"github.com/labstack/echo/v4"
)

func initTableRoutes(e *echo.Echo) {
	g := e.Group("/table", middleware.JWTMiddleware)

	g.POST("", controller.CreateTable)
	g.GET("", controller.GetGames)
	g.POST("/join", controller.JoinTable)
	g.GET("/:id", controller.GetTable)
	g.DELETE("/:id", controller.DeleteTable)

	// Invitation routes
	g.GET("/:id/invitations", controller.GetPendingInvitations)
	g.POST("/:id/invitations/:userId/accept", controller.AcceptInvitation)
	g.POST("/:id/invitations/:userId/decline", controller.DeclineInvitation)
	g.POST("/:id/invite-code", controller.RegenerateInviteCode)
	g.GET("/:id/players", controller.GetGamePlayers)
	g.DELETE("/:id/players/:userId", controller.RemovePlayer)
}
