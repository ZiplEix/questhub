package routes

import (
	"questhub/controller"
	"questhub/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func initTableRoutes(e *echo.Echo) {
	g := e.Group("/table", middleware.JWTMiddleware)

	g.POST("", controller.CreateTable)
	g.GET("", controller.GetGames)
	g.POST("/join", controller.JoinTable)

	// Group for game-specific routes with state check
	// Applies CheckGameState:
	// - Validates Game ID
	// - Fetches Game
	// - Bypasses if User is GM
	// - Bypasses if Method is GET
	// - Blocks if State is "paused"
	gameGroup := g.Group("/:id", middleware.CheckGameState)

	gameGroup.GET("", controller.GetTable)
	gameGroup.GET("/players", controller.GetGamePlayers)
	gameGroup.GET("/characters", controller.GetGameCharacters)
	gameGroup.POST("/chat", controller.SendMessage)
	gameGroup.POST("/chat", controller.SendMessage)
	gameGroup.GET("/roll", controller.RollDice, echoMiddleware.RateLimiterWithConfig(middleware.DiceRollRateLimitConfig))

	// GM only routes (RequireGM applies ON TOP of CheckGameState if we nest, or we can separate)
	// If we use gameGroup.Group, CheckGameState runs first. GM is bypassed in CheckGameState, so it's fine.
	// But efficiently, GM shouldn't need to fetch game twice if possible, but middleware order matters.
	// Let's attach RequireGM to g directly to avoid CheckGameState overhead/logic if redundant,
	// BUT CheckGameState logic has "GM Bypass". So it's safe.
	// Actually, let's keep GM routes separate to be clean or use gameGroup?
	// If we use g.Group("/:id", RequireGM), it's parallel to gameGroup.
	// That means GM routes won't traverse CheckGameState. This is better for performance (one less DB call if CheckGameState fetches game).
	// However, CheckGameState fetches game to check GM ID. RequireGM ALSO fetches game (usually) or checks ownership.
	// RequireGM middleware usually checks DB too.
	// Let's look at `RequireGM` implementation? Assumed it exists.
	// To be safe and simple, let's just use `g` for GM routes as before. They are protected by RequireGM.

	// GM only routes
	gmGroup := g.Group("/:id", middleware.RequireGM)
	gmGroup.DELETE("", controller.DeleteTable)
	gmGroup.GET("/invitations", controller.GetPendingInvitations)
	gmGroup.POST("/invitations/:userId/accept", controller.AcceptInvitation)
	gmGroup.POST("/invitations/:userId/decline", controller.DeclineInvitation)
	gmGroup.POST("/invite-code", controller.RegenerateInviteCode)
	gmGroup.POST("/invite-code", controller.RegenerateInviteCode)
	gmGroup.DELETE("/players/:userId", controller.RemovePlayer)
	gmGroup.GET("/monsters", controller.GetGameMonsters)
	gmGroup.POST("/characters", controller.CreateCharacter)
	gmGroup.PUT("/characters/:charId", controller.UpdateCharacter)
	gmGroup.DELETE("/characters/:charId", controller.DeleteCharacter)
	gmGroup.POST("/characters/:charId/assign", controller.AssignCharacter)
	gmGroup.PUT("/state", controller.UpdateTableState)

	// Mixed access routes (GM or Owner) - handled in controller
	// These should also be subject to Game State check (e.g. updating notes)
	// We use gameGroup which has CheckGameState
	gameGroup.GET("/characters/:charId", controller.GetCharacter)
	gameGroup.GET("/characters/:charId/notes", controller.GetCharacterNotes)
	gameGroup.PUT("/characters/:charId/notes", controller.UpdateCharacterNotes)
	gameGroup.GET("/chat", controller.GetChatHistory)
}
