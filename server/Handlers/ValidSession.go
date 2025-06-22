package Handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-app/Database"
	"real-time-app/Helpers"
	"real-time-app/Models"
	"time"
)

// Add this to your Handlers package
func ValidateSession(r *http.Request) (*Models.User, string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, "", fmt.Errorf("no session cookie")
	}

	var userID int
	var expiresAt time.Time
	err = Database.DB.QueryRow(`
        SELECT user_id, expires_at FROM sessions 
        WHERE id = ? AND expires_at > ?`,
		cookie.Value, time.Now()).Scan(&userID, &expiresAt)

	if err != nil {
		return nil, "", fmt.Errorf("invalid or expired session")
	}

	user, err := Helpers.GetUserByID(userID) // You'll need to create this function
	if err != nil {
		return nil, "", fmt.Errorf("user not found")
	}

	return user, cookie.Value, nil
}

func ValidateSessionHandler(w http.ResponseWriter, r *http.Request) {

	user, sessionID, err := ValidateSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Println("✅ Session validated for:", user.Email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":       user.ID,
			"email":    user.Email,
			"nickname": user.Nickname,
		},
		"session_id": sessionID,
	})
}
