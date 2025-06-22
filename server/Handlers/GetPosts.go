package Handlers

import (
	"encoding/json"
	"net/http"
	"real-time-app/Database"
	"real-time-app/Models"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var currentUserID int
	err = Database.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ?", cookie.Value).Scan(&currentUserID)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	query := `
SELECT p.id, p.user_id, p.username, p.title, p.content, p.image, p.created_at,
       u.visibility,
       EXISTS(SELECT 1 FROM likes l WHERE l.user_id = ? AND l.post_id = p.id) AS liked,
       EXISTS(SELECT 1 FROM dislikes d WHERE d.user_id = ? AND d.post_id = p.id) AS disliked
FROM posts p
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN follows f ON f.follower_id = ? AND f.followee_id = p.user_id AND f.status = 'accepted'
WHERE p.user_id = ?  -- Include your own posts always
   OR u.visibility = 'public'
   OR f.status = 'accepted'
ORDER BY p.created_at DESC
`

rows, err := Database.DB.Query(query, currentUserID, currentUserID, currentUserID, currentUserID)

	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Post struct {
		ID        int    `json:"id"`
		UserID    int    `json:"user_id"`
		Username  string `json:"username"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		Image     string `json:"image"`
		CreatedAt string `json:"created_at"`
		Liked     bool   `json:"liked"`
		Disliked  bool   `json:"disliked"`
	}

	var posts []Post
	for rows.Next() {
		var p Post
		var visibility string
		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.Image, &p.CreatedAt, &visibility, &p.Liked, &p.Disliked); err == nil {
			posts = append(posts, p)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetPostByID(w http.ResponseWriter, r *http.Request) {
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
	var post Models.Post
	err = Database.DB.QueryRow(
		"SELECT id, user_id, username, title, content, image, created_at, like_count FROM posts WHERE id = ?",
		postID,
	).Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content, &post.Image, &post.CreatedAt, &post.LikeCount)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Check if current user liked this post
	var userID int
	cookie, err := r.Cookie("session_id")
	if err == nil {
		_ = Database.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ?", cookie.Value).Scan(&userID)
		if userID != 0 {
			var liked bool
			likeErr := Database.DB.QueryRow(
				"SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = ?)",
				userID, post.ID,
			).Scan(&liked)
			if likeErr == nil {
				post.Liked = liked
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
