package websocket

import (
	"cluster-talk-backend/internal/db"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the

// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Rooms map: RoomID -> Set of Clients
	rooms map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *WSMessage

	// Register requests from the clients.
	register chan *ClientSubscription

	// Unregister requests from clients.
	unregister chan *ClientSubscription
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *WSMessage),
		register:   make(chan *ClientSubscription),
		unregister: make(chan *ClientSubscription),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.register:
			// Register client to global list
			h.clients[sub.Client] = true

			// Register to specific room
			roomID := sub.RoomID
			if _, ok := h.rooms[roomID]; !ok {
				h.rooms[roomID] = make(map[*Client]bool)
			}
			h.rooms[roomID][sub.Client] = true

			// Send History (Async)
			go func(c *Client, rid string) {
				msgs, err := db.GlobalDB.GetHistory(rid)
				if err == nil {
					for _, m := range msgs {
						// Convert DB msg to WS msg
						wsMsg := &WSMessage{
							Type:      "message",
							Sender:    m.Sender,
							Payload:   m.Payload,
							RoomID:    m.RoomID,
							Timestamp: m.Timestamp,
						}
						c.send <- wsMsg
					}
				} else {
					log.Printf("Error fetching history: %v", err)
				}
			}(sub.Client, roomID)

		case sub := <-h.unregister:
			if _, ok := h.clients[sub.Client]; ok {
				delete(h.clients, sub.Client)
				close(sub.Client.send)

				// Remove from room as well
				if room, ok := h.rooms[sub.RoomID]; ok {
					delete(room, sub.Client)
					if len(room) == 0 {
						delete(h.rooms, sub.RoomID)
					}
				}
			}
		case message := <-h.broadcast:
			// Save to DB (Async)
			go db.GlobalDB.SaveMessage(message.Sender, message.Payload, message.RoomID, message.Timestamp)

			// Broadcast only to clients in the specific room
			if room, ok := h.rooms[message.RoomID]; ok {

				for client := range room {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
						delete(room, client)
					}
				}
			}
		}
	}
}
