package Handlers

import (
	"encoding/json"
	"net/http"
	"real-time-app/Database"
	"strings"
)

func SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if strings.TrimSpace(query) == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}
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

	// Search users excluding self
	rows, err := Database.DB.Query(`
        SELECT id, nickname FROM users
        WHERE (nickname LIKE ? OR email LIKE ?)
        AND id != ?
        LIMIT 20
    `, "%"+query+"%", "%"+query+"%", currentUserID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var nickname string
		if err := rows.Scan(&id, &nickname); err == nil {
			results = append(results, map[string]interface{}{
				"id":       id,
				"nickname": nickname,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
