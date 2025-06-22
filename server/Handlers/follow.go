package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-app/Database"
	"strconv"

	"github.com/gorilla/mux"
)

// --- SEND FOLLOW REQUEST ---
func SendFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session and current user ID (follower)
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var followerID int
	err = Database.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ?", cookie.Value).Scan(&followerID)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// Get followee ID from URL param
	vars := mux.Vars(r)
	followeeID, err := strconv.Atoi(vars["id"])
	if err != nil || followeeID == followerID {
		http.Error(w, "Invalid followee ID", http.StatusBadRequest)
		return
	}

	// Check the followee's visibility
	var visibility string
	err = Database.DB.QueryRow("SELECT visibility FROM users WHERE id = ?", followeeID).Scan(&visibility)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Determine status based on visibility
	status := "pending"
	if visibility == "public" {
		status = "accepted"
	}

	// Insert or update the follow request
	_, err = Database.DB.Exec(`
        INSERT INTO follows (follower_id, followee_id, status)
        VALUES (?, ?, ?)
        ON CONFLICT(follower_id, followee_id) DO UPDATE SET status=?`,
		followerID, followeeID, status, status)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to send follow request", http.StatusInternalServerError)
		return
	}

	message := "Follow request sent"
	// after inserting notification for follow
	if status == "accepted" {
		// --- Notification for public profile: Follow is auto-accepted
		res, _ := Database.DB.Exec(`
		INSERT INTO notifications (recipient_id, sender_id, type, message)
		VALUES (?, ?, 'follow', ?)`,
			followeeID, followerID, "You have a new follower.",
		)
		notifID, _ := res.LastInsertId()
		var senderNickname string
		_ = Database.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", followerID).Scan(&senderNickname)
		notif := Notification{
			ID:             int(notifID),
			Type:           "follow",
			SenderID:       followerID,
			SenderNickname: senderNickname,
			Message:        "You have a new follower.",
		}
		SendNotification(followeeID, notif)
	} else if status == "pending" {
		// --- Notification for private profile: New follow request
		res, _ := Database.DB.Exec(`
		INSERT INTO notifications (recipient_id, sender_id, type, message)
		VALUES (?, ?, 'follow_request', ?)`,
			followeeID, followerID, "You have a new follow request.",
		)
		notifID, _ := res.LastInsertId()
		var senderNickname string
		_ = Database.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", followerID).Scan(&senderNickname)
		notif := Notification{
			ID:             int(notifID),
			Type:           "follow_request",
			SenderID:       followerID,
			SenderNickname: senderNickname,
			Message:        "You have a new follow request.",
		}
		SendNotification(followeeID, notif)
	}

	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func GetFollowStatus(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid target user ID", http.StatusBadRequest)
		return
	}

	var status string
	err = Database.DB.QueryRow(`
		SELECT status FROM follows 
		WHERE follower_id = ? AND followee_id = ?`, currentUserID, targetID).Scan(&status)

	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(map[string]string{"status": "none"})
		return
	} else if err != nil {
		http.Error(w, "Error checking follow status", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

// --- ACCEPT FOLLOW REQUEST ---
func AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authenticate followee
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var followeeID int
	err = Database.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ?", cookie.Value).Scan(&followeeID)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	followerID, err := strconv.Atoi(vars["follower_id"])
	if err != nil || followerID == followeeID {
		http.Error(w, "Invalid follower ID", http.StatusBadRequest)
		return
	}

	var currentStatus string
	err = Database.DB.QueryRow(`
		SELECT status FROM follows 
		WHERE follower_id = ? AND followee_id = ?`,
		followerID, followeeID).Scan(&currentStatus)

	if err == sql.ErrNoRows {
		http.Error(w, "Follow request not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if currentStatus != "pending" {
		http.Error(w, "Cannot accept a non-pending request", http.StatusBadRequest)
		return
	}

	_, err = Database.DB.Exec(`
		UPDATE follows SET status = 'accepted'
		WHERE follower_id = ? AND followee_id = ?`,
		followerID, followeeID)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to accept follow request", http.StatusInternalServerError)
		return
	}
	_, _ = Database.DB.Exec(`
		INSERT INTO notifications (recipient_id, sender_id, type, message)
		VALUES (?, ?, 'follow', ?)`,
		followerID, followeeID, "Your follow request was accepted.",
	)

	json.NewEncoder(w).Encode(map[string]string{"message": "Follow request accepted"})
}

func GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
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
		SELECT u.id, u.nickname, u.email
		FROM follows f
		JOIN users u ON f.follower_id = u.id
		WHERE f.followee_id = ? AND f.status = 'pending'`, userID)
	if err != nil {
		http.Error(w, "Failed to fetch requests", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pending []map[string]interface{}
	for rows.Next() {
		var id int
		var nickname, email string
		if err := rows.Scan(&id, &nickname, &email); err == nil {
			pending = append(pending, map[string]interface{}{
				"id":       id,
				"nickname": nickname,
				"email":    email,
			})
		}
	}

	json.NewEncoder(w).Encode(pending)
}
