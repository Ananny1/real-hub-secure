package Handlers

import (
	"encoding/json"
	"net/http"
	"real-time-app/Database"
)

func GetProfileStats(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	var postCount, followingCount, followerCount int
	Database.DB.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&postCount)
	Database.DB.QueryRow("SELECT COUNT(*) FROM follows WHERE follower_id = ? AND status = 'accepted'", userID).Scan(&followingCount)
	Database.DB.QueryRow("SELECT COUNT(*) FROM follows WHERE followee_id = ? AND status = 'accepted'", userID).Scan(&followerCount)

	stats := map[string]int{
		"postCount":      postCount,
		"followingCount": followingCount,
		"followerCount":  followerCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ✅ شلون نخلي الصفحة تجيب هالأرقام بسرعة؟
// نسوي endpoint صغير وسريع (ما يرد معلومات زايدة).

// كل ما تفتح صفحتك أو تسوي تحديث، frontend يطلب هذا الـ endpoint.

// يردله response سهل:

// {
//   "postCount": 8,
//   "followingCount": 23,
//   "followerCount": 17
// }
// مباشرة يعرضهم فوق في الصفحة بدون تحميل بوستات أو بيانات زايدة.
