package migration

import (
	"log"
	"real-time-app/Database"
)

func CreateTables() {
	createUsersTable := `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nickname TEXT NOT NULL,
    gender TEXT,
    age INTEGER,
    first_name TEXT,
    last_name TEXT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    visibility TEXT CHECK(visibility IN ('public', 'private')) DEFAULT 'public'  -- 👈 Add this
);`

	createSessionTable := `
	CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	createPostsTable := `
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    username TEXT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image TEXT,                          -- Image URL or file path
    like_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id)
);`

	createLikeTable := `
	CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(post_id) REFERENCES posts(id),
    UNIQUE(user_id, post_id) -- ensures a user can like each post only once
);
`
	createDisLikeTable := `
CREATE TABLE IF NOT EXISTS dislikes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(post_id) REFERENCES posts(id),
    UNIQUE(user_id, post_id) -- ensures a user can dislike each post only once
);
`
	createFollowsTable := `
CREATE TABLE IF NOT EXISTS follows (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id INTEGER NOT NULL,
    followee_id INTEGER NOT NULL,
    status TEXT CHECK(status IN ('pending', 'accepted', 'remove')) NOT NULL DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (follower_id, followee_id)
);
`
	createNotificationsTable := `
CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    recipient_id INTEGER NOT NULL,
    sender_id INTEGER,
    type TEXT NOT NULL,
    message TEXT NOT NULL,
    post_id INTEGER, -- <------ Add this line!
    is_read BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(recipient_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
);

`

	createCommentsTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		user_id INTEGER,
		username TEXT,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	createMessagesTable := `
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER,
    receiver_id INTEGER,
    message_content TEXT,
    image_url TEXT,  -- <--- NEW
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(sender_id) REFERENCES users(id),
    FOREIGN KEY(receiver_id) REFERENCES users(id)
	);`

	queries := []string{
		createUsersTable,
		createSessionTable,
		createPostsTable,
		createLikeTable,
		createFollowsTable,
		createDisLikeTable,
		createCommentsTable,
		createMessagesTable,
		createNotificationsTable,
	}

	for _, query := range queries {
		_, err := Database.DB.Exec(query)
		if err != nil {
			log.Fatal("Migration failed:", err)
			return
		}
	}

	log.Println("Migration successful: All tables created")
}
