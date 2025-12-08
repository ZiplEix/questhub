package websocket

import (
	"encoding/json"
	"log"
	"questhub/service"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

var GlobalHub *Hub

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// Parse message to get game_id and check if it's private
			var msgMap map[string]any
			if err := json.Unmarshal(message, &msgMap); err != nil {
				log.Printf("error unmarshalling message: %v", err)
				continue
			}

			// Get Game ID
			gameID, _ := msgMap["game_id"].(string)

			// Handle Private Messages first (optimization: check type before players?)
			// But private messages effectively are restricted to game usually,
			// but the logic here handles explicit targeting.
			if msgType, ok := msgMap["type"].(string); ok && msgType == "CHAT_PRIVATE" {
				if targetID, ok := msgMap["target_id"].(string); ok && targetID != "" {
					h.BroadcastToUser(targetID, message)
					// Also send back to sender so they see their own private message
					if senderID, ok := msgMap["sender_id"].(string); ok && senderID != targetID {
						h.BroadcastToUser(senderID, message)
					}
					continue
				}
			}

			// Get allowed users for this game
			allowedUsers := make(map[string]bool)
			if gameID != "" {
				players, err := service.GetGamePlayers(gameID)
				if err == nil {
					for _, p := range players {
						allowedUsers[p.UserID] = true
					}
				} else {
					log.Printf("error fetching game players for broadcast: %v", err)
				}
			}

			for client := range h.clients {
				// Filter by game
				if gameID != "" {
					if !allowedUsers[client.UserID] {
						continue
					}
				}

				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// BroadcastToUser sends a message to a specific user.
// This is a placeholder for more targeted broadcasting if needed.
func (h *Hub) BroadcastToUser(userID string, message []byte) {
	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// Broadcast sends a message to all clients.
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}
