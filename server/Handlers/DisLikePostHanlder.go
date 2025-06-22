package Handlers

import (
	"database/sql"
	"net/http"
	"real-time-app/Database"
)

func DisLikePostHandler(w http.ResponseWriter, r *http.Request) {
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
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	postID := r.FormValue("post_id")
	if postID == "" {
		http.Error(w, "Missing post_id", http.StatusBadRequest)
		return
	}
	var dislikeID int
	err = Database.DB.QueryRow("SELECT id FROM dislikes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&dislikeID)
	if err == sql.ErrNoRows {
		_, err := Database.DB.Exec("INSERT INTO dislikes (user_id, post_id) VALUES (?, ?)", userID, postID)
		if err != nil {
			http.Error(w, "Failed to dislike", http.StatusInternalServerError)
			return
		}

		_, _ = Database.DB.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"disliked": true, "liked": false}`))
		return
	} else if err == nil {
		_, err := Database.DB.Exec("DELETE FROM dislikes WHERE id = ?", dislikeID)
		if err != nil {
			http.Error(w, "Failed to remove dislike", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"disliked": false}`))
		return
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
