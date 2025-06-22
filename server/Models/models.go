package Models

import "time"

type User struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	ID        string    `json:"id"`         // UUID or secure random string
	UserID    int       `json:"user_id"`    // Foreign key to User table
	ExpiresAt time.Time `json:"expires_at"` // When the session expires
}

type Post struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Image     string `json:"image"`        // <-- add this line!
	CreatedAt string `json:"created_at"`
	LikeCount int    `json:"like_count"`
	Liked     bool   `json:"liked"`        // true if the current user liked it
}

type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type Message struct {
	ID          int    `json:"id"`
	SenderID    int    `json:"sender_id"`
	ReceiverID  int    `json:"receiver_id"`
	MessageText string `json:"message_text"`
	CreatedAt   string `json:"created_at"`
}
