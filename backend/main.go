package main

import (
	"log"
	"os"

	"questhub/config"
	"questhub/database"
	mdw "questhub/middleware"

	"context"
	"time"

	"github.com/ZiplEix/better-logs/httpmw"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.InitEnv()

	database.InitDB()

	logger, cleanup := config.InitLogger()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := cleanup(ctx); err != nil {
			log.Printf("logger cleanup error: %v", err)
		}
	}()
	// Prevent unused variable error if logger is not used directly in main
	_ = logger

	authURL := os.Getenv("BETTER_AUTH_URL")
	if err := mdw.InitJWKS(authURL); err != nil {
		log.Fatalf("Failed to initialize JWKS: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	// Add Better Logs middleware
	e.Use(echo.WrapMiddleware(httpmw.Middleware))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	log.Println("Server running on http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))

}
