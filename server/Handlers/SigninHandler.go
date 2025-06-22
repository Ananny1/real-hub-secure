package Handlers

import (
	"encoding/json"
	"net/http"
	"real-time-app/Database"
	"real-time-app/Helpers"
	"real-time-app/Models"
	"time"

	"github.com/google/uuid"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var credential Models.Credential

	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if credential.Email == "" || credential.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := Helpers.GetUserByEmail(credential.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !Helpers.CheckPassword(user.Password, credential.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	sessionID := uuid.NewString()
	expires := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // Lax works for most local dev
		Secure:   false,                // false for localhost
		Path:     "/",
	})
	_, err = Database.DB.Exec(`
	INSERT INTO sessions (id, user_id, expires_at)
	VALUES (?, ?, ?)`, sessionID, user.ID, expires)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}
