package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"questhub/models/database"
	"questhub/service"
	"questhub/websocket"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func RollDice(c echo.Context) error {
	gameID := c.Param("id")
	sidesStr := c.QueryParam("sides")
	sides, err := strconv.Atoi(sidesStr)
	if err != nil || sides <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sides parameter"})
	}

	// 1. Generate Result
	rand.Seed(time.Now().UnixNano())
	result := rand.Intn(sides) + 1

	// 2. Identify Sender
	// Get claims from context
	claims, ok := c.Get("claims").(jwt.MapClaims)
	senderID := ""
	senderName := "System" // Fallback
	if ok {
		if sub, ok := claims["sub"].(string); ok {
			senderID = sub
		}
		if name, ok := claims["name"].(string); ok {
			senderName = name
		}
	} else {
		fmt.Println("Error: could not cast claims to jwt.MapClaims")
	}

	// Fetch Game to check if sender is GM
	game, err := service.GetTable(gameID)
	if err != nil {
		// Log error but proceed? Or fail? proceeding might be safer, default to system/user name
		fmt.Printf("Error getting table for roll: %v\n", err)
	}

	isGM := false
	if game != nil && game.GmID == senderID {
		isGM = true
		senderName = "GM"
	}

	// Try to resolve Character Name for the user in this game UNLESS they are GM
	if senderID != "" && !isGM {
		char, err := service.GetUserCharacter(gameID, senderID)
		if err == nil && char != nil {
			// If it's a player character, use that name
			if char.Type != "GM_HIDDEN" {
				senderName = char.Name
			}
		}
	}

	// Check for secret roll
	isSecret := c.QueryParam("secret") == "true"
	if isSecret && !isGM {
		// Only GM can roll secretly
		isSecret = false
	}

	// 3. Create Chat Message Content
	// content := fmt.Sprintf("🎲 d%d : %d", sides, result)
	// If secret, we might just want the roll content, as UI adds "Secret" label for CHAT_PRIVATE
	content := fmt.Sprintf("🎲 d%d : %d", sides, result)

	msgType := "EVENT"
	var targetID *string

	if isSecret {
		msgType = "CHAT_PRIVATE"
		targetID = &senderID
	}

	// 4. Construct Message
	msg := database.ChatMessage{
		GameID:     gameID,
		SenderID:   senderID,
		SenderName: senderName,
		Content:    content,
		Type:       msgType,
		TargetID:   targetID,
		CreatedAt:  time.Now(),
	}

	// 5. Persist Message
	if err := service.SaveMessage(msg); err != nil {
		fmt.Printf("Error saving roll message: %v\n", err)
	}

	// 6. Broadcast via WebSocket
	// Hub handles CHAT_PRIVATE routing automatically if Type is CHAT_PRIVATE and TargetID is set
	websocket.GlobalHub.BroadcastToGame(gameID, msg)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"result": result,
		"sides":  sides,
	})
}
