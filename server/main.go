package main

import (
	"fmt"
	"net/http"
	"real-time-app/Database"
	"real-time-app/Database/migration"
	"real-time-app/Handlers"

	"github.com/gorilla/mux"
)

// func loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("%s %s from %s\n", r.Method, r.RequestURI, r.RemoteAddr)
// 		next.ServeHTTP(w, r)
// 	})
// }

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,DELETE,PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	Database.ConnectDatabase()
	migration.CreateTables()

	r := mux.NewRouter()

	r.Use(corsMiddleware)
	// r.Use(loggingMiddleware)
	r.HandleFunc("/chat/upload", Handlers.UploadChatImageHandler).Methods("POST")
	r.HandleFunc("/chat/history", Handlers.GetChatHistory).Methods("GET")
	r.HandleFunc("/chat/users", Handlers.UserListHandler).Methods("GET")
	r.HandleFunc("/ws", Handlers.Ws).Methods("GET")
	r.HandleFunc("/notifications", Handlers.GetNotifications).Methods("GET")
	r.HandleFunc("/profile/visibility", Handlers.UpdateVisibilityHandler).Methods("POST", "PATCH")
	r.HandleFunc("/profile/stats", Handlers.GetProfileStats).Methods("GET")
	r.HandleFunc("/follow/accept/{follower_id}", Handlers.AcceptFollowRequest).Methods("POST")
	r.HandleFunc("/follow/{id}", Handlers.SendFollowRequest).Methods("POST")
	r.HandleFunc("/follow/status/{id}", Handlers.GetFollowStatus).Methods("GET")
	r.HandleFunc("/users/search", Handlers.SearchUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", Handlers.GetUserPublicProfile).Methods("GET")
	r.HandleFunc("/validate-session", Handlers.ValidateSessionHandler).Methods("GET")
	r.HandleFunc("/profile/user", Handlers.GetUserProfileHandler).Methods("GET")
	r.HandleFunc("/profile/myposts", Handlers.GetMyPostsHandler).Methods("GET")
	r.HandleFunc("/profile/liked", Handlers.ProfileLikedHandler).Methods("GET")
	r.HandleFunc("/dislike", Handlers.DisLikePostHandler).Methods("POST")
	r.HandleFunc("/like", Handlers.LikePostHandler).Methods("POST")
	r.HandleFunc("/posts/{id}/comments", Handlers.AddComment).Methods("POST")
	r.HandleFunc("/posts/{id}/comments", Handlers.GetComments).Methods("GET")
	r.HandleFunc("/posts/{id}", Handlers.GetPostByID).Methods("GET")
	r.HandleFunc("/posts", Handlers.GetPostsHandler).Methods("GET")
	r.HandleFunc("/posts", Handlers.CreatePostHandler).Methods("POST")
	r.HandleFunc("/logout", Handlers.LogoutHandler).Methods("POST")
	r.HandleFunc("/signup", Handlers.SignUpHandler).Methods("POST")
	r.HandleFunc("/login", Handlers.SignInHandler).Methods("POST")
	r.HandleFunc("/follow/{id}", Handlers.UnfollowUser).Methods("DELETE")

	r.HandleFunc("/", Handlers.HomeHandler).Methods("GET")

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// OPTIONS preflight for CORS (optional if corsMiddleware handles all)
	r.PathPrefix("/").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	fmt.Println("✅ Server is running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("❌ Failed to start server:", err)
	}
}
