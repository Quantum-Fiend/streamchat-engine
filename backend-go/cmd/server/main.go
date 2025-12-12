package main

import (
	"log"
	"net/http"

	"net/http"

	"cluster-talk-backend/internal/db"
	"cluster-talk-backend/internal/websocket"
)

func main() {
	// Init Database
	db.InitDB()

	hub := websocket.NewHub()

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	// Add a health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("ClusterTalk Backend started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
