package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	model "questhub/models/database"
	"questhub/models/request"
	"questhub/service"
	"questhub/websocket"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func CreateTable(c echo.Context) error {
	var req request.CreateTableRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name is required")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	gmID := claims["sub"].(string)

	game, err := service.CreateTable(req.Name, gmID, req.ImageURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create table").SetInternal(err)
	}

	return c.JSON(http.StatusCreated, game)
}

type GameResponse struct {
	*model.Game
	CurrentCharacterID *string `json:"current_character_id"`
	IsGM               bool    `json:"is_gm"`
	DebugMsg           string  `json:"debug_msg,omitempty"`
}

func GetTable(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	game, err := service.GetTable(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}

	response := GameResponse{
		Game: game,
		IsGM: game.GmID == userID,
	}

	// Fetch user character
	// If GM, ensure GM character exists
	// Fetch user character
	// If GM, ensure GM character exists
	if response.IsGM {
		gmChar, err := service.EnsureGMCharacter(id, userID)
		if err != nil {
			fmt.Printf("Failed to ensure GM character: %v\n", err)
			response.DebugMsg = fmt.Sprintf("Failed to ensure GM character: %v", err)
		} else if gmChar != nil {
			response.CurrentCharacterID = &gmChar.ID
			response.DebugMsg = fmt.Sprintf("GM Character ensured: %s", gmChar.ID)
		} else {
			response.DebugMsg = "EnsureGMCharacter returned nil char"
		}
	} else {
		// Regular player character fetch
		character, err := service.GetUserCharacter(id, userID)
		if err != nil {
			fmt.Printf("Failed to fetch user character: %v\n", err)
			response.DebugMsg = fmt.Sprintf("Failed to fetch user character: %v", err)
		} else if character != nil {
			response.CurrentCharacterID = &character.ID
			response.DebugMsg = fmt.Sprintf("User character found: %s", character.ID)
		} else {
			response.DebugMsg = "No character found for user"
		}
	}

	return c.JSON(http.StatusOK, response)
}

func GetGames(c echo.Context) error {
	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	games, err := service.GetGames(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch games").SetInternal(err)
	}

	return c.JSON(http.StatusOK, games)
}

func JoinTable(c echo.Context) error {
	var req struct {
		InviteCode string `json:"invite_code"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.InviteCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invite code is required")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	gameID, err := service.JoinTable(req.InviteCode, userID)
	if err != nil {
		// Check for unique constraint violation (user already in game or invited)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.JSON(http.StatusConflict, map[string]string{
				"message": "User already in game or invitation pending",
				"id":      gameID,
			})
		}
		return echo.NewHTTPError(http.StatusNotFound, "Game not found or failed to join")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Invitation sent", "id": gameID})
}

func GetPendingInvitations(c echo.Context) error {
	gameID := c.Param("id")
	if gameID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	// Verify GM - Handled by middleware

	invitations, err := service.GetPendingInvitations(gameID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch invitations").SetInternal(err)
	}

	return c.JSON(http.StatusOK, invitations)
}

func AcceptInvitation(c echo.Context) error {
	gameID := c.Param("id")
	targetUserID := c.Param("userId")

	if gameID == "" || targetUserID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or user ID")
	}

	// Verify GM - Handled by middleware

	if err := service.AcceptInvitation(gameID, targetUserID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to accept invitation").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Invitation accepted"})
}

func DeclineInvitation(c echo.Context) error {
	gameID := c.Param("id")
	targetUserID := c.Param("userId")

	if gameID == "" || targetUserID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or user ID")
	}

	// Verify GM - Handled by middleware

	if err := service.DeclineInvitation(gameID, targetUserID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decline invitation").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Invitation declined"})
}

func DeleteTable(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	if err := service.DeleteTable(id, userID); err != nil {
		if err.Error() == "unauthorized: only the GM can delete the table" {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete table").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Table deleted successfully"})
}

func RegenerateInviteCode(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	// Verify GM - Handled by middleware

	newCode, err := service.RegenerateInviteCode(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to regenerate invite code").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"invite_code": newCode})
}

func GetGamePlayers(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	// In a real app, we might want to check if the user is part of the game or is the GM
	// For now, we'll allow any authenticated user to see the players of a game they know the ID of
	// Or we can restrict it to GM only if this is strictly for the settings page.
	// The requirement says "GM settings page", so let's restrict to GM for now to be safe,
	// but usually players can see other players. Let's stick to GM check as it's for the settings page.

	// Verify GM - Handled by middleware

	players, err := service.GetGamePlayers(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch players").SetInternal(err)
	}

	return c.JSON(http.StatusOK, players)
}

func RemovePlayer(c echo.Context) error {
	gameID := c.Param("id")
	playerID := c.Param("userId")

	if gameID == "" || playerID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or player ID")
	}

	// Verify GM - Handled by middleware
	// We still need game object for GM ID check in RemovePlayer logic below?
	// Actually RemovePlayer service doesn't check GM ID, it just removes.
	// But we need to know who the GM is to prevent removing the GM.
	game, err := service.GetTable(gameID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}

	// Prevent removing GM
	if playerID == game.GmID {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot remove the Game Master")
	}

	err = service.RemovePlayer(gameID, playerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove player").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Player removed successfully"})
}

func GetGameCharacters(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	// Verify GM - Handled by middleware

	characters, err := service.GetGameCharacters(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch characters").SetInternal(err)
	}

	return c.JSON(http.StatusOK, characters)
}

func GetCharacter(c echo.Context) error {
	gameID := c.Param("id")
	charID := c.Param("charId")
	if gameID == "" || charID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or character ID")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	// Verify GM or Owner
	game, err := service.GetTable(gameID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}

	// Check if user is GM
	isGM := game.GmID == userID

	character, err := service.GetCharacter(gameID, charID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch character").SetInternal(err)
	}
	if character == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Character not found")
	}

	// If not GM, check if user owns the character
	if !isGM {
		if character.UserID == nil || *character.UserID != userID {
			return echo.NewHTTPError(http.StatusForbidden, "You can only view your own character or you must be the GM")
		}
	}

	return c.JSON(http.StatusOK, character)
}

func GetGameMonsters(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	// Verify GM - Handled by middleware

	characters, err := service.GetGameMonsters(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch monsters").SetInternal(err)
	}

	return c.JSON(http.StatusOK, characters)
}

func CreateCharacter(c echo.Context) error {
	gameID := c.Param("id")
	if gameID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	// Verify GM - Handled by middleware

	name := c.FormValue("name")
	race := c.FormValue("race")
	isNPC := c.FormValue("is_npc") == "true"
	maxHPStr := c.FormValue("max_hp")
	stats := []byte(c.FormValue("stats"))
	inventory := []byte(c.FormValue("inventory"))
	money, _ := strconv.Atoi(c.FormValue("money"))
	initiative, _ := strconv.Atoi(c.FormValue("initiative"))
	age := c.FormValue("age")
	height := c.FormValue("height")
	weight := c.FormValue("weight")
	maxSpells, _ := strconv.Atoi(c.FormValue("max_spells"))
	spells := []byte(c.FormValue("spells"))
	abilities := c.FormValue("abilities")
	experience, _ := strconv.Atoi(c.FormValue("experience"))
	charType := c.FormValue("type")
	subRace := c.FormValue("sub_race")
	armorClass, _ := strconv.Atoi(c.FormValue("armor_class"))
	speed, _ := strconv.Atoi(c.FormValue("speed"))

	// Default type logic if not provided
	if charType == "" {
		if isNPC {
			charType = "NPC"
		} else {
			charType = "PLAYER"
		}
	}

	var maxHP int
	if _, err := fmt.Sscanf(maxHPStr, "%d", &maxHP); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Max HP")
	}

	// Handle avatar upload
	avatarURL := c.FormValue("avatar_url")
	file, err := c.FormFile("avatar")
	if err == nil {
		// Upload file
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer func() {
			_ = src.Close()
		}()

		// Create uploads directory if not exists
		if err := os.MkdirAll("uploads", 0755); err != nil {
			return err
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dstPath := filepath.Join("uploads", filename)

		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer func() {
			_ = dst.Close()
		}()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// Construct full URL
		scheme := c.Scheme()
		host := c.Request().Host
		avatarURL = fmt.Sprintf("%s://%s/uploads/%s", scheme, host, filename)
	}

	// Handle inventory images
	var inventoryItems []map[string]interface{}
	if len(inventory) > 0 {
		if err := json.Unmarshal(inventory, &inventoryItems); err != nil {
			fmt.Println("Error unmarshalling inventory:", err)
		} else {
			updatedInventory := false
			for i, item := range inventoryItems {
				formKey := fmt.Sprintf("inventory_image_%d", i)
				invFile, err := c.FormFile(formKey)
				if err == nil {
					src, err := invFile.Open()
					if err != nil {
						continue
					}
					defer func() {
						_ = src.Close()
					}()

					ext := filepath.Ext(invFile.Filename)
					filename := fmt.Sprintf("inv_%d_%d%s", i, time.Now().UnixNano(), ext)
					dstPath := filepath.Join("uploads", filename)

					dst, err := os.Create(dstPath)
					if err != nil {
						continue
					}
					defer func() {
						_ = dst.Close()
					}()

					if _, err = io.Copy(dst, src); err != nil {
						continue
					}

					scheme := c.Scheme()
					host := c.Request().Host
					itemURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, host, filename)
					item["image_url"] = itemURL
					updatedInventory = true
				}
			}
			if updatedInventory {
				inventory, _ = json.Marshal(inventoryItems)
			}
		}
	}

	char, err := service.CreateCharacter(gameID, "", name, race, maxHP, isNPC, avatarURL, stats, inventory, money, initiative, age, height, weight, maxSpells, spells, abilities, experience, charType, subRace, armorClass, speed)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create character").SetInternal(err)
	}

	return c.JSON(http.StatusCreated, char)
}

func UpdateCharacter(c echo.Context) error {
	gameID := c.Param("id")
	charID := c.Param("charId")
	if gameID == "" || charID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or character ID")
	}

	// Verify GM - Handled by middleware

	name := c.FormValue("name")
	race := c.FormValue("race")
	maxHPStr := c.FormValue("max_hp")
	isNPC, _ := strconv.ParseBool(c.FormValue("is_npc"))
	stats := []byte(c.FormValue("stats"))
	inventory := []byte(c.FormValue("inventory"))
	money, _ := strconv.Atoi(c.FormValue("money"))
	initiative, _ := strconv.Atoi(c.FormValue("initiative"))
	age := c.FormValue("age")
	height := c.FormValue("height")
	weight := c.FormValue("weight")
	maxSpells, _ := strconv.Atoi(c.FormValue("max_spells"))
	spells := []byte(c.FormValue("spells"))
	abilities := c.FormValue("abilities")
	experience, _ := strconv.Atoi(c.FormValue("experience"))
	charType := c.FormValue("type")
	subRace := c.FormValue("sub_race")
	armorClass, _ := strconv.Atoi(c.FormValue("armor_class"))
	speed, _ := strconv.Atoi(c.FormValue("speed"))

	// Default type logic if not provided
	if charType == "" {
		if isNPC {
			charType = "NPC"
		} else {
			charType = "PLAYER"
		}
	}

	var maxHP int
	if _, err := fmt.Sscanf(maxHPStr, "%d", &maxHP); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Max HP")
	}

	// Handle avatar upload
	avatarURL := c.FormValue("avatar_url")
	file, err := c.FormFile("avatar")
	if err == nil {
		// Upload file
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer func() {
			_ = src.Close()
		}()

		// Create uploads directory if not exists
		if err := os.MkdirAll("uploads", 0755); err != nil {
			return err
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dstPath := filepath.Join("uploads", filename)

		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer func() {
			_ = dst.Close()
		}()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// Construct full URL
		scheme := c.Scheme()
		host := c.Request().Host
		avatarURL = fmt.Sprintf("%s://%s/uploads/%s", scheme, host, filename)
	}

	// Handle inventory images
	var inventoryItems []map[string]interface{}
	if len(inventory) > 0 {
		if err := json.Unmarshal(inventory, &inventoryItems); err != nil {
			fmt.Println("Error unmarshalling inventory:", err)
		} else {
			updatedInventory := false
			for i, item := range inventoryItems {
				formKey := fmt.Sprintf("inventory_image_%d", i)
				invFile, err := c.FormFile(formKey)
				if err == nil {
					src, err := invFile.Open()
					if err != nil {
						continue
					}
					defer func() {
						_ = src.Close()
					}()

					ext := filepath.Ext(invFile.Filename)
					filename := fmt.Sprintf("inv_%d_%d%s", i, time.Now().UnixNano(), ext)
					dstPath := filepath.Join("uploads", filename)

					dst, err := os.Create(dstPath)
					if err != nil {
						continue
					}
					defer func() {
						_ = dst.Close()
					}()

					if _, err = io.Copy(dst, src); err != nil {
						continue
					}

					scheme := c.Scheme()
					host := c.Request().Host
					itemURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, host, filename)
					item["image_url"] = itemURL
					updatedInventory = true
				}
			}
			if updatedInventory {
				inventory, _ = json.Marshal(inventoryItems)
			}
		}
	}

	char, err := service.UpdateCharacter(charID, gameID, name, race, maxHP, isNPC, avatarURL, stats, inventory, money, initiative, age, height, weight, maxSpells, spells, abilities, experience, charType, subRace, armorClass, speed)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update character").SetInternal(err)
	}

	// Broadcast update to the character owner
	if char.UserID != nil && websocket.GlobalHub != nil {
		msg := map[string]any{
			"type":    "CHARACTER_UPDATE",
			"payload": char,
		}
		msgBytes, _ := json.Marshal(msg)
		websocket.GlobalHub.BroadcastToUser(*char.UserID, msgBytes)
	}

	return c.JSON(http.StatusOK, char)
}

func GetChatHistory(c echo.Context) error {
	gameID := c.Param("id")
	claims := c.Get("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)
	messages, err := service.GetGameMessages(gameID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch chat history").SetInternal(err)
	}
	return c.JSON(http.StatusOK, messages)
}

func DeleteCharacter(c echo.Context) error {
	gameID := c.Param("id")
	charID := c.Param("charId")
	if gameID == "" || charID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or character ID")
	}

	// Verify GM - Handled by middleware

	if err := service.DeleteCharacter(charID, gameID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete character").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Character deleted successfully"})
}

func AssignCharacter(c echo.Context) error {
	gameID := c.Param("id")
	charID := c.Param("charId")
	if gameID == "" || charID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or character ID")
	}

	var req struct {
		PlayerID string `json:"player_id"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Verify GM - Handled by middleware

	if err := service.AssignCharacterToPlayer(gameID, charID, req.PlayerID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to assign character").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Character assigned successfully"})
}

func GetCharacterNotes(c echo.Context) error {
	gameID := c.Param("id")
	charID := c.Param("charId")
	if gameID == "" || charID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or character ID")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	requestingUserID := claims["sub"].(string)

	// Verify GM or Owner
	game, err := service.GetTable(gameID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}

	isGM := game.GmID == requestingUserID

	character, err := service.GetCharacter(gameID, charID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch character").SetInternal(err)
	}
	if character == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Character not found")
	}

	targetUserID := requestingUserID
	if isGM {
		// If GM is viewing, they want to see the owner's notes
		if character.UserID != nil {
			targetUserID = *character.UserID
		} else {
			// Character has no user, so no player notes exist.
			// Attempting to fetch notes for "no user" is invalid or we can return empty.
			return c.JSON(http.StatusOK, map[string]string{"content": ""})
		}
	} else {
		// If not GM, must own the character
		if character.UserID == nil || *character.UserID != requestingUserID {
			return echo.NewHTTPError(http.StatusForbidden, "You can only view notes for your own character or you must be the GM")
		}
	}

	notes, err := service.GetNotes(gameID, targetUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch notes").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"content": notes})
}

func UpdateCharacterNotes(c echo.Context) error {
	gameID := c.Param("id")
	charID := c.Param("charId")
	if gameID == "" || charID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID or character ID")
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	requestingUserID := claims["sub"].(string)

	// Verify GM or Owner
	game, err := service.GetTable(gameID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}

	isGM := game.GmID == requestingUserID

	character, err := service.GetCharacter(gameID, charID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch character").SetInternal(err)
	}
	if character == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Character not found")
	}

	targetUserID := requestingUserID
	if isGM {
		// If GM is updating, they update the owner's notes
		if character.UserID != nil {
			targetUserID = *character.UserID
		} else {
			// Cannot update notes for unassigned character as notes are tied to Player+Game
			return echo.NewHTTPError(http.StatusBadRequest, "Cannot add notes to an unassigned character (notes are linked to players)")
		}
	} else {
		// If not GM, must own the character
		if character.UserID == nil || *character.UserID != requestingUserID {
			return echo.NewHTTPError(http.StatusForbidden, "You can only update notes for your own character or you must be the GM")
		}
	}

	if err := service.UpdateNotes(gameID, targetUserID, req.Content); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update notes").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Notes updated successfully"})
}

func SendMessage(c echo.Context) error {
	gameID := c.Param("id")
	if gameID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	var req struct {
		Content    string  `json:"content"`
		Type       string  `json:"type"`
		TargetID   *string `json:"target_id,omitempty"`
		SenderName string  `json:"sender_name"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	claims := c.Get("claims").(jwt.MapClaims)
	senderID := claims["sub"].(string)

	msg := model.ChatMessage{
		GameID:     gameID,
		SenderID:   senderID,
		SenderName: req.SenderName,
		Content:    req.Content,
		Type:       req.Type,
		TargetID:   req.TargetID,
		CreatedAt:  time.Now(),
	}

	if err := service.SaveMessage(msg); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save message").SetInternal(err)
	}

	// Broadcast via WebSocket
	if websocket.GlobalHub != nil {
		msgBytes, _ := json.Marshal(msg)
		websocket.GlobalHub.Broadcast(msgBytes)
	}

	return c.JSON(http.StatusOK, msg)
}

func UpdateTableState(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing game ID")
	}

	var req struct {
		State string `json:"state"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.State != "ongoing" && req.State != "paused" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid state. Must be 'ongoing' or 'paused'")
	}

	// Verify GM - Handled by middleware

	if err := service.UpdateGameState(id, req.State); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update game state").SetInternal(err)
	}

	// Optional: Broadcast state change via WS?
	if websocket.GlobalHub != nil {
		msg := map[string]string{
			"type":    "GAME_STATE_UPDATE",
			"game_id": id,
			"state":   req.State,
		}
		msgBytes, _ := json.Marshal(msg)
		websocket.GlobalHub.Broadcast(msgBytes)
	}

	return c.JSON(http.StatusOK, map[string]string{"state": req.State})
}
