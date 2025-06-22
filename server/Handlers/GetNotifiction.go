package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-app/Database"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello1")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var userID int
	err = Database.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ?", cookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	rows, err := Database.DB.Query(`
SELECT n.id, n.type, n.sender_id, u.nickname, n.post_id, n.created_at
FROM notifications n
JOIN users u ON n.sender_id = u.id
WHERE n.recipient_id = ?
ORDER BY n.created_at DESC
LIMIT 30

	`, userID)

	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer rows.Close()

	type Notification struct {
		ID             int    `json:"id"`
		Type           string `json:"type"`
		SenderID       int    `json:"sender_id"`
		SenderNickname string `json:"sender_nickname"`
		PostID         *int   `json:"post_id,omitempty"`
		CreatedAt      string `json:"created_at"`
	}

	var notifications []Notification
	for rows.Next() {
		var n Notification
		var postID sql.NullInt64
		if err := rows.Scan(&n.ID, &n.Type, &n.SenderID, &n.SenderNickname, &postID, &n.CreatedAt); err == nil {
			if postID.Valid {
				pid := int(postID.Int64)
				n.PostID = &pid
			}
			notifications = append(notifications, n)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}
