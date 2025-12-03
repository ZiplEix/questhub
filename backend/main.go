package main

import (
	"log"
	"os"

	"questhub/config"
	"questhub/database"
	mdw "questhub/middleware"
	"questhub/routes"

	"github.com/ZiplEix/better-logs/httpmw"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"questhub/websocket"
)

func main() {
	config.InitEnv()

	database.InitDB()

	// logger, cleanup := config.InitLogger()
	// defer func() {
	// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 	defer cancel()
	// 	if err := cleanup(ctx); err != nil {
	// 		log.Printf("logger cleanup error: %v", err)
	// 	}
	// }()
	// // Prevent unused variable error if logger is not used directly in main
	// _ = logger

	authURL := os.Getenv("BETTER_AUTH_URL")
	if err := mdw.InitJWKS(authURL); err != nil {
		log.Fatalf("Failed to initialize JWKS: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	// Add Better Logs middleware with skipper for WebSocket
	betterLogsMdw := echo.WrapMiddleware(httpmw.Middleware)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip for WebSocket route to allow hijacking
			if c.Request().URL.Path == "/ws" {
				return next(c)
			}
			return betterLogsMdw(next)(c)
		}
	})

	// WebSocket Hub
	hub := websocket.NewHub()
	websocket.GlobalHub = hub
	go hub.Run()

	// WebSocket Route (Protected)
	e.GET("/ws", func(c echo.Context) error {
		return websocket.ServeWs(hub, c)
	}, mdw.JWTMiddleware)

	// Serve static files
	e.Static("/uploads", "uploads")

	routes.SetupRoutes(e)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	log.Println("Server running on http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))

}
