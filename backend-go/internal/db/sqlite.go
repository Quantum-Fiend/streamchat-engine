package db

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

type DB struct {
	conn *sql.DB
}

type Message struct {
	Sender    string `json:"sender"`
	Payload   string `json:"payload"`
	RoomID    string `json:"room_id"`
	Timestamp int64  `json:"timestamp"`
}

var GlobalDB *DB

func InitDB() {
	db, err := sql.Open("sqlite", "clustertalk.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create Table
	query := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        room_id TEXT,
        sender TEXT,
        payload TEXT,
        timestamp INTEGER
    );
    `
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("SQLite Database Initialized")
	GlobalDB = &DB{conn: db}
}

func (d *DB) SaveMessage(sender, payload, roomID string, timestamp int64) {
	_, err := d.conn.Exec("INSERT INTO messages (room_id, sender, payload, timestamp) VALUES (?, ?, ?, ?)",
		roomID, sender, payload, timestamp)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
	}
}

func (d *DB) GetHistory(roomID string) ([]Message, error) {
	rows, err := d.conn.Query("SELECT sender, payload, room_id, timestamp FROM messages WHERE room_id = ? ORDER BY id DESC LIMIT 50", roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.Sender, &m.Payload, &m.RoomID, &m.Timestamp); err != nil {
			continue
		}
		// Prepend to maintain chronological order if needed, but client handles sort usually.
		// Actually, let's reverse them or let client handle.
		// We fetched DESC (newest first).
		messages = append(messages, m)
	}

	// Reverse to show oldest first in chat list
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
