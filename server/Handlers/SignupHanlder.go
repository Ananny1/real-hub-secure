package Handlers

import (
	"encoding/json"
	"net/http"
	"real-time-app/Database"
	Helpers "real-time-app/Helpers"
	"real-time-app/Models"
	"time"

	"github.com/google/uuid"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user Models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if user.Email == "" || user.Password == "" || user.Nickname == "" || user.Age < 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := Helpers.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	res, err := Database.DB.Exec(`
		INSERT INTO users (nickname, email, password, gender, age, first_name, last_name)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.Nickname, user.Email, hashedPassword, user.Gender, user.Age, user.FirstName, user.LastName)

	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	userID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting user ID", http.StatusInternalServerError)
		return
	}

	// Create session
	sessionID := uuid.NewString()
	expires := time.Now().Add(24 * time.Hour)

	_, err = Database.DB.Exec(`INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`,
		sessionID, userID, expires)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // Lax works for most local dev
		Secure:   false,                // false for localhost
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
}
