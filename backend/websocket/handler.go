package websocket

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true // Allow all origins for now, configure as needed
	}

	// Extract UserID from JWT via context (set by middleware)
	var userID string
	claims, ok := c.Get("claims").(jwt.MapClaims)
	if ok {
		if sub, ok := claims["sub"].(string); ok {
			userID = sub
		}
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), UserID: userID}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	return nil
}
