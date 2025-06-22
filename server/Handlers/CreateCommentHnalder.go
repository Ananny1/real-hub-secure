package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-app/Database"
	"real-time-app/Models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	fmt.Println(postID)

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Println(cookie)

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

	fmt.Println(userID)
	var username string
	err = Database.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	fmt.Println(username)
	var req struct {
		Content string `json:"content"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := Database.DB.Exec(
		"INSERT INTO comments (post_id, user_id, username, content) VALUES (?, ?, ?, ?)",
		postID, userID, username, req.Content)

	commentID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}
	fmt.Println(commentID)

	comment := Models.Comment{
		ID:        int(commentID),
		PostID:    postID,
		UserID:    userID,
		Username:  username,
		Content:   req.Content,
		CreatedAt: now,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}
