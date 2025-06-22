package Handlers

import (
	"encoding/json"
	"net/http"
	"real-time-app/Database"
	"real-time-app/Models"
	"strconv"

	"github.com/gorilla/mux"
)

func GetComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	rows, err := Database.DB.Query(
		"SELECT id, post_id, user_id, username, content, created_at FROM comments WHERE post_id = ? ORDER BY created_at DESC",
		postID,
	)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Models.Comment
	for rows.Next() {
		var c Models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt); err != nil {
			http.Error(w, "Error reading comments", http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
