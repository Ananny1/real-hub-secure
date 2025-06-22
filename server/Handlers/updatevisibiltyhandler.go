package Handlers

import (
	"encoding/json"
	"net/http"
	"real-time-app/Database"
)

func UpdateVisibilityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

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

	type Req struct {
		Visibility string `json:"visibility"`
	}
	var body Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if body.Visibility != "public" && body.Visibility != "private" {
		http.Error(w, "Invalid visibility", http.StatusBadRequest)
		return
	}

	result, err := Database.DB.Exec("UPDATE users SET visibility = ? WHERE id = ?", body.Visibility, userID)
	if err != nil {
		http.Error(w, "Failed to update DB", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No user updated", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"visibility": body.Visibility,
	})
}
