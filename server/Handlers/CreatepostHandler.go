package Handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-app/Database"
	"real-time-app/Helpers"
	"real-time-app/Models"
	"time"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "In`ternal server error", http.StatusInternalServerError)
		}
		return
	}

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

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}
	var imageFilename string
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageFilename, err = Helpers.SaveFile(file, handler)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
	} else if err != http.ErrMissingFile {
		http.Error(w, "Failed to get image file", http.StatusBadRequest)
		return
	} else {
		imageFilename = ""
	}

	result, err := Database.DB.Exec(
		"INSERT INTO posts (user_id, username, title, content, image) VALUES (?, ?, ?, ?, ?)",
		userID, username, title, content, imageFilename,
	)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}
	postID, _ := result.LastInsertId()

	post := Models.Post{
		ID:        int(postID),
		UserID:    userID,
		Username:  username,
		Title:     title,
		Content:   content,
		Image:     imageFilename,
		CreatedAt: time.Now().Format(time.RFC3339),
		LikeCount: 0,
		Liked:     false,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)

}
