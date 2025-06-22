package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-app/Database"
	"real-time-app/Models"
)

func ProfileLikedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProfileLikedHandler called")
	if r.Method != http.MethodGet {
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
		return // << You need to return here!
	}
	rows, err := Database.DB.Query(`
        SELECT p.id, p.user_id, p.username, p.title, p.content, p.image, p.created_at,
               (SELECT COUNT(*) FROM likes WHERE post_id = p.id) AS like_count
        FROM posts p
        JOIN likes l ON p.id = l.post_id
        JOIN users u ON p.user_id = u.id
        WHERE l.user_id = ?
        ORDER BY p.created_at DESC
    `, userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Models.Post
	for rows.Next() {
		var post Models.Post
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content,
			&post.Image, &post.CreatedAt, &post.LikeCount,
		)
		if err != nil {
			continue
		}
		post.Liked = true // Always true for this endpoint
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetMyPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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
	}
	rows, err := Database.DB.Query(`
        SELECT p.id, p.user_id, p.username, p.title, p.content, p.image, p.created_at, p.like_count,
               CASE WHEN l.id IS NULL THEN false ELSE true END AS liked
        FROM posts p
        LEFT JOIN likes l ON l.post_id = p.id AND l.user_id = ?
        WHERE p.user_id = ?
        ORDER BY p.created_at DESC
    `, userID, userID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Models.Post
	for rows.Next() {
		var p Models.Post
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content,
			&p.Image, &p.CreatedAt, &p.LikeCount, &p.Liked,
		); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	var username, visibility string
	err = Database.DB.QueryRow("SELECT nickname, visibility FROM users WHERE id = ?", userID).Scan(&username, &visibility)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"username":   username,
		"visibility": visibility,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
