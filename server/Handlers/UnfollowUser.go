package Handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-app/Database"
	"strconv"

	"github.com/gorilla/mux"
)

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followeeIDStr := vars["id"]

	// Step 1: Get session_id cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized: missing session cookie", http.StatusUnauthorized)
		return
	}
	sessionToken := cookie.Value
	fmt.Println("Session token:", sessionToken)

	// Step 2: Get user ID from session
	var userID int
	err = Database.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ?", sessionToken).Scan(&userID)
	if err != nil {
		fmt.Println("Session lookup failed:", err)
		http.Error(w, "Unauthorized: invalid session", http.StatusUnauthorized)
		return
	}

	// Step 3: Parse followee ID
	followeeID, err := strconv.Atoi(followeeIDStr)
	if err != nil {
		http.Error(w, "Invalid followee ID", http.StatusBadRequest)
		return
	}

	// Step 4: Delete the follow relationship
	result, err := Database.DB.Exec(
		"DELETE FROM follows WHERE follower_id = ? AND followee_id = ?",
		userID, followeeID,
	)
	if err != nil {
		fmt.Println("Database delete failed:", err)
		http.Error(w, "Failed to unfollow", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No such follow relationship", http.StatusNotFound)
		return
	}

	// Step 5: Return success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Unfollowed successfully",
	})
}
