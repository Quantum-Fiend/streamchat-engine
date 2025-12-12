package websocket

// WSMessage defines the JSON structure for WebSocket messages
type WSMessage struct {
	Type      string `json:"type"`      // "join", "message", "leave", "system"
	Payload   string `json:"payload"`   // Content
	Sender    string `json:"sender"`    // UserID/Name
	RoomID    string `json:"room_id"`   // Target Room
	Timestamp int64  `json:"timestamp"` // Unix Epoch
}

// ClientSubscription bundles a client and their target room
type ClientSubscription struct {
	Client *Client
	RoomID string
}
