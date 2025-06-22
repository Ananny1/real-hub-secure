package Handlers

import (
	"database/sql"
	"net/http"
	"real-time-app/Database"
	"strconv"
)

type Notification struct {
	ID             int    `json:"id,omitempty"`
	Type           string `json:"type"`
	SenderID       int    `json:"sender_id"`
	SenderNickname string `json:"sender_nickname,omitempty"`
	PostID         *int   `json:"post_id,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	Message        string `json:"message,omitempty"`
}

func SendNotification(userID int, notif Notification) {
	if conn, ok := userConnections[userID]; ok {
		conn.WriteJSON(notif)
	}
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

	postID := r.FormValue("post_id")
	if postID == "" {
		http.Error(w, "Post ID required", http.StatusBadRequest)
		return
	}
	postIDInt, _ := strconv.Atoi(postID)

	// 1. Find the post's owner
	var postOwnerID int
	err = Database.DB.QueryRow(
		"SELECT user_id FROM posts WHERE id = ?", postID,
	).Scan(&postOwnerID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	var likeID int
	err = Database.DB.QueryRow(
		"SELECT id FROM likes WHERE user_id = ? AND post_id = ?", userID, postID,
	).Scan(&likeID)

	if err == sql.ErrNoRows {
		_, err := Database.DB.Exec(
			"INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, postID,
		)
		if err != nil {
			http.Error(w, "Failed to like", http.StatusInternalServerError)
			return
		}

		// 3. Don't notify yourself for liking your own post
		if postOwnerID != userID {
			// Insert notification
			res, _ := Database.DB.Exec(
				`INSERT INTO notifications (recipient_id, sender_id, type, post_id, message)
                 VALUES (?, ?, 'like', ?, ?)`,
				postOwnerID, userID, postIDInt, "Someone liked your post.",
			)
			notifID, _ := res.LastInsertId()
			// Real-time notification (fetch sender nickname)
			var senderNickname string
			_ = Database.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&senderNickname)
			notif := Notification{
				ID:             int(notifID),
				Type:           "like",
				SenderID:       userID,
				SenderNickname: senderNickname,
				PostID:         &postIDInt,
				Message:        "Someone liked your post.",
			}
			SendNotification(postOwnerID, notif)
		}

		w.Write([]byte(`{"liked": true}`))
		return
	} else if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// Already liked: remove like
	_, err = Database.DB.Exec(
		"DELETE FROM likes WHERE id = ?", likeID,
	)
	if err != nil {
		http.Error(w, "Failed to unlike", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"liked": false}`))
}
