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

	// GM only routes
	gmGroup := g.Group("/:id", middleware.RequireGM)
	gmGroup.DELETE("", controller.DeleteTable)
	gmGroup.GET("/invitations", controller.GetPendingInvitations)
	gmGroup.POST("/invitations/:userId/accept", controller.AcceptInvitation)
	gmGroup.POST("/invitations/:userId/decline", controller.DeclineInvitation)
	gmGroup.POST("/invite-code", controller.RegenerateInviteCode)
	gmGroup.GET("/players", controller.GetGamePlayers)
	gmGroup.DELETE("/players/:userId", controller.RemovePlayer)
	gmGroup.GET("/characters", controller.GetGameCharacters)
	gmGroup.GET("/monsters", controller.GetGameMonsters)
	gmGroup.POST("/characters", controller.CreateCharacter)
	gmGroup.PUT("/characters/:charId", controller.UpdateCharacter)
	gmGroup.DELETE("/characters/:charId", controller.DeleteCharacter)
	gmGroup.POST("/characters/:charId/assign", controller.AssignCharacter)

	// Mixed access routes (GM or Owner) - handled in controller
	g.GET("/:id/characters/:charId", controller.GetCharacter)
}
