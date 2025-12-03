package websocket

import "encoding/json"

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
			// Parse message to check if it's private
			var msgMap map[string]interface{}
			if err := json.Unmarshal(message, &msgMap); err == nil {
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
			}

			for client := range h.clients {
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
