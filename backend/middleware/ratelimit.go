package middleware

import (
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// DiceRollRateLimitConfig defines the configuration for dice roll rate limiting
var DiceRollRateLimitConfig = echoMiddleware.RateLimiterConfig{
	Skipper: echoMiddleware.DefaultSkipper,
	Store: echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
		echoMiddleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Every(3 * time.Second),
			Burst:     1,
			ExpiresIn: 3 * time.Minute,
		},
	),
	IdentifierExtractor: func(ctx echo.Context) (string, error) {
		// Identify by User ID if available (JWT), otherwise IP
		// Dice roll route is protected by JWT, so we should always have claims.
		// BUT fallback to IP is safe.
		// However, user provided example used RealIP().
		// If I use UserID, it's better for authenticated users.
		// User requirement "3sec par utilisateur".
		// Let's try to get User ID from context first.
		claims, ok := ctx.Get("claims").(jwt.MapClaims)
		if ok {
			if sub, ok := claims["sub"].(string); ok {
				return sub, nil
			}
		}
		// Fallback to IP
		return ctx.RealIP(), nil
	},
	ErrorHandler: func(context echo.Context, err error) error {
		return context.JSON(http.StatusTooManyRequests, map[string]string{
			"error": "Veuillez attendre 3 secondes avant de réessayer.",
		})
	},
	DenyHandler: func(context echo.Context, identifier string, err error) error {
		return context.JSON(http.StatusTooManyRequests, map[string]string{
			"error": "Veuillez attendre 3 secondes avant de réessayer.",
		})
	},
}
