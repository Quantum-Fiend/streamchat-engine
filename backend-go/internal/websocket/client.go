package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	// Define AI Service URL
	aiServiceURL = "http://localhost:8000/moderate"
	// Define Analytics Service URL (Rust)
	analyticsServiceURL = "http://localhost:9000/ingest"

	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for now to simplify local dev with Vite
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	// We send objects now, but for writing to WS we still serialize to byte array or just send JSON text
	send chan *WSMessage

	// Current Room
	roomID string
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- &ClientSubscription{Client: c, RoomID: c.roomID}
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg WSMessage
		err := c.conn.ReadJSON(&msg) // Automatically unmarshal JSON
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Enforce server-side validation or overriding
		msg.Sender = "User-" + c.conn.RemoteAddr().String() // Simple ID for now
		msg.Timestamp = time.Now().Unix()

		// If the msg.RoomID is empty, default to current room or general
		if msg.RoomID == "" {
			msg.RoomID = "general"
		}

		// Handle "join" type to switch rooms?
		// For now, assume every message is just broadcast to the room stated in the message
		// But strictly, the Client is subscribed to c.roomID.
		// Let's enforce that the Client publishes to their room.
		msg.RoomID = c.roomID

		// ----------------------------------------------------------------
		// AI MODERATION STEP
		// ----------------------------------------------------------------
		if msg.Payload != "" {
			// Prepare payload
			aiPayload := map[string]string{
				"text":    msg.Payload,
				"user_id": msg.Sender,
			}
			jsonData, _ := json.Marshal(aiPayload)

			// Call Python Service
			resp, err := http.Post(aiServiceURL, "application/json", bytes.NewBuffer(jsonData))
			if err == nil {
				defer resp.Body.Close()
				var moderationRes struct {
					IsToxic      bool   `json:"is_toxic"`
					FilteredText string `json:"filtered_text"`
				}
				if json.NewDecoder(resp.Body).Decode(&moderationRes) == nil {
					// If toxic, we can either block or replace
					if moderationRes.IsToxic {
						msg.Payload = "[CENSORED BY AI]: " + moderationRes.FilteredText
					}
				}
			} else {
				// Fail open or log error
				log.Printf("AI Service Unreachable: %v", err)
			}
		}
		// ----------------------------------------------------------------

		c.hub.broadcast <- &msg

		// ----------------------------------------------------------------
		// ANALYTICS STEP (Async Fire & Forget)
		// ----------------------------------------------------------------
		go func(m *WSMessage) {
			payload := map[string]interface{}{
				"event_type": "message",
				"room_id":    m.RoomID,
				"timestamp":  m.Timestamp,
			}
			jsonData, _ := json.Marshal(payload)
			// We ignore errors here for fire & forget
			http.Post(analyticsServiceURL, "application/json", bytes.NewBuffer(jsonData))
		}(&msg)
		// ----------------------------------------------------------------

	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}

			// Note: We removed the "queuing" logic loop here for simplicity with WriteJSON
			// In high perf scenarios we might bundle messages.

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		roomID = "general"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan *WSMessage, 256), roomID: roomID}
	client.hub.register <- &ClientSubscription{Client: client, RoomID: roomID}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
