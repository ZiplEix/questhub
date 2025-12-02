package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var jwks *keyfunc.JWKS

func InitJWKS(authURL string) error {
	jwksURL := authURL + "/api/auth/jwks"

	const maxRetries = 5
	const retryDelay = 5 * time.Second

	var err error
	for i := range maxRetries {
		jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{
			RefreshUnknownKID: true,
		})
		if err == nil {
			return nil
		}
		log.Printf("Failed to initialize JWKS (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryDelay)
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("failed to initialize JWKS after %d attempts: %w", maxRetries, err)
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid Authorization header")
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, jwks.Keyfunc)
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
		}

		c.Set("claims", claims)

		return next(c)
	}
}
