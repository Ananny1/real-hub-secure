package Handlers

import (
	"encoding/json"
	"net/http"
	"real-time-app/Database"
	"strconv"

	"github.com/gorilla/mux"
)

// Struct for returning public posts
type PublicPost struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
}

func GetUserPublicProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch user info including visibility
	var nickname, firstName, lastName, visibility string
	err = Database.DB.QueryRow(`
        SELECT nickname, first_name, last_name, visibility
        FROM users WHERE id = ?
    `, userID).Scan(&nickname, &firstName, &lastName, &visibility)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get counts
	var postCount, followerCount, followingCount int
	Database.DB.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&postCount)
	Database.DB.QueryRow("SELECT COUNT(*) FROM follows WHERE followee_id = ? AND status = 'accepted'", userID).Scan(&followerCount)
	Database.DB.QueryRow("SELECT COUNT(*) FROM follows WHERE follower_id = ? AND status = 'accepted'", userID).Scan(&followingCount)

	// Determine if viewer can see posts
	showPosts := visibility == "public"

	// Check session and follow status if profile is private
	if !showPosts {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			var currentUserID int
			err := Database.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ?", cookie.Value).Scan(&currentUserID)
			if err == nil {
				var followStatus string
				err := Database.DB.QueryRow(`
					SELECT status FROM follows
					WHERE follower_id = ? AND followee_id = ?
				`, currentUserID, userID).Scan(&followStatus)
				if err == nil && followStatus == "accepted" {
					showPosts = true
				}
			}
		}
	}

	// Fetch posts only if allowed
	var posts []PublicPost
	if showPosts {
		rows, err := Database.DB.Query(`
			SELECT id, title, image, created_at
			FROM posts
			WHERE user_id = ?
			ORDER BY created_at DESC
			LIMIT 10
		`, userID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var post PublicPost
				if err := rows.Scan(&post.ID, &post.Title, &post.Image, &post.CreatedAt); err == nil {
					posts = append(posts, post)
				}
			}
		}
	}

	// Build response
	data := map[string]interface{}{
		"nickname":        nickname,
		"first_name":      firstName,
		"last_name":       lastName,
		"visibility":      visibility,
		"post_count":      postCount,
		"follower_count":  followerCount,
		"following_count": followingCount,
		"can_view_posts":  showPosts,
		"posts":           posts, // Empty array if not allowed
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
