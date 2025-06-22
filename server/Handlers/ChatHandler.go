package Handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-app/Database"
	"real-time-app/Helpers"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Type     string `json:"type"`
	From     int    `json:"from"`
	To       int    `json:"to"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url,omitempty"`
}

// Add this function to your Handlers file:

func UploadChatImageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form (e.g., 10 MB max memory)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	// Get the uploaded file
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uniqueFilename, err := Helpers.SaveFile(file, handler)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Build the public URL (adjust if needed)
	imageURL := "http://localhost:8080/uploads/" + uniqueFilename

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": imageURL,
	})
}

func SaveMessageToDB(from, to int, content, imageURL string) error {
	fmt.Printf("DB INSERT TRY: %d -> %d | %s | %s\n", from, to, content, imageURL)
	_, err := Database.DB.Exec(
		"INSERT INTO messages (sender_id, receiver_id, message_content, image_url) VALUES (?, ?, ?, ?)",
		from, to, content, imageURL,
	)
	if err != nil {
		fmt.Println("DB error:", err)
	}
	return err
}

func HandleChatMessage(fromUser int, raw []byte) {
	fmt.Println("Raw WS payload:", string(raw))

	var msg ChatMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		fmt.Println("Invalid chat message json:", err)
		return
	}
	fmt.Printf("Parsed chat message: %+v\n", msg)

	if msg.Type != "chat" || msg.To == 0 {
		fmt.Println("Invalid chat message structure")
		return
	}

	msg.From = fromUser

	if err := SaveMessageToDB(msg.From, msg.To, msg.Content, msg.ImageURL); err != nil {
		fmt.Println("Failed to save message:", err)
	} else {
		fmt.Println("Saved chat to DB!")
	}

	if conn, ok := userConnections[msg.To]; ok {
		data, _ := json.Marshal(msg)
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

func GetChatHistory(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var userID int
	err = Database.DB.QueryRow(
		"SELECT user_id FROM sessions WHERE id = ?", cookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// 2. Parse query param: 'with' (other user)
	otherIDStr := r.URL.Query().Get("with")
	var otherID int
	if otherIDStr == "" {
		http.Error(w, "Missing or invalid 'with' param", http.StatusBadRequest)
		return
	}
	if n, err := fmt.Sscanf(otherIDStr, "%d", &otherID); n != 1 || err != nil {
		http.Error(w, "Missing or invalid 'with' param", http.StatusBadRequest)
		return
	}

	// 3. Query messages between the two users (both directions)
	rows, err := Database.DB.Query(`
    SELECT sender_id, receiver_id, message_content, image_url, created_at
    FROM messages
    WHERE (sender_id=? AND receiver_id=?) OR (sender_id=? AND receiver_id=?)
    ORDER BY created_at ASC
`, userID, otherID, otherID, userID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Message struct {
		SenderID   int    `json:"sender_id"`
		ReceiverID int    `json:"receiver_id"`
		Content    string `json:"content"`
		ImageURL   string `json:"image_url,omitempty"`
		Timestamp  string `json:"timestamp"`
	}

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.SenderID, &m.ReceiverID, &m.Content, &m.ImageURL, &m.Timestamp); err == nil {
			messages = append(messages, m)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
