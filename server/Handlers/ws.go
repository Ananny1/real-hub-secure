package Handlers

import (
	"fmt"
	"net/http"
	"real-time-app/Database"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var userConnections = make(map[int]*websocket.Conn)

func RegisterChatConn(userID int, conn *websocket.Conn) {
	userConnections[userID] = conn
}

func UnregisterChatConn(userID int) {
	delete(userConnections, userID)
}
func chatReader(conn *websocket.Conn, userID int, nickname string) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}
		fmt.Printf("Received message from user %d: %s\n", userID, message)
		HandleChatMessage(userID, message) // <-- this does everything: DB + forward!
	}
}

func Ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Get userID and nickname from session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized: No session", http.StatusUnauthorized)
		return
	}
	var userID int
	var nickname string
	err = Database.DB.QueryRow(
		"SELECT id, nickname FROM users WHERE id = (SELECT user_id FROM sessions WHERE id = ?)",
		cookie.Value).Scan(&userID, &nickname)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return
	}

	RegisterChatConn(userID, conn)
	SetUserOnline(userID, nickname)
	defer func() {
		UnregisterChatConn(userID)
		SetUserOffline(userID)
	}()

	fmt.Printf("WebSocket connection established for user %d (%s)\n", userID, nickname)

	chatReader(conn, userID, nickname)
}
