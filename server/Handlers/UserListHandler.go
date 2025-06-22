package Handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type OnlineUser struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
}

var (
	onlineUsers   = make(map[int]*OnlineUser)
	onlineUsersMu sync.Mutex
)

func SetUserOnline(id int, nickname string) {
	onlineUsersMu.Lock()
	onlineUsers[id] = &OnlineUser{ID: id, Nickname: nickname}
	fmt.Printf("Now online: %#v\n", onlineUsers)
	onlineUsersMu.Unlock()
}

func SetUserOffline(id int) {
	onlineUsersMu.Lock()
	delete(onlineUsers, id)
	onlineUsersMu.Unlock()
}

func GetOnlineUsers() []*OnlineUser {
	onlineUsersMu.Lock()
	defer onlineUsersMu.Unlock()
	fmt.Printf("Serving user list: %#v\n", onlineUsers)
	list := make([]*OnlineUser, 0, len(onlineUsers))
	for _, u := range onlineUsers {
		list = append(list, u)
	}
	return list
}

func UserListHandler(w http.ResponseWriter, r *http.Request) {
	users := GetOnlineUsers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
