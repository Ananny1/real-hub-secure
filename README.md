# 🔐 Real Hub Secure

**Real Hub Secure** is a real-time, full-stack social application designed for expressive posting, instant messaging, live notifications, and account control — all styled in a sleek cyber-dark UI with glowing accents and futuristic vibes.

---

## ✨ What’s Inside?

- 🖼 **Create and Explore Posts** — Share ideas, media, and styled content.
- 💬 **Real-Time Chat** — Message other users instantly using WebSockets.
- 🔔 **Notifications** — Get notified when you're followed or liked.
- 👥 **Follow System** — Public/private profile visibility with friend tracking.
- 🔎 **Live User Search** — Search users instantly by nickname or email.
- ⚙️ **Profile Control** — View your stats, update visibility, manage liked posts.

---

## 🚀 How to Run It

### 1. **Client (Frontend – React)**
```bash
cd client
npm install
npm run dev
```

Runs your React frontend using **Vite** on `localhost:5173`.

---

### 2. **Server (Backend – Go)**
```bash
cd server
go run .
```

Starts the backend on `localhost:8080`, providing:
- REST APIs
- WebSocket endpoints
- Session/cookie-based auth
- SQLite database access

---

## 🌐 Tech Stack

- **Frontend**: React, AOS, Tailwind/CSS, DiceBear Avatars
- **Backend**: Go (net/http), SQLite, REST + WebSockets
- **Realtime**: Chat, online status, notifications

---

## 🧪 Development Notes

- Uses cookies (`session_id`) for auth
- All API calls are CORS enabled and expect `credentials: "include"`
- SQLite database with basic schema for users, posts, likes, sessions

---

## 📄 License

MIT — open source, built for fun and productivity.