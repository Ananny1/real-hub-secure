import "../../styles/Notif.css";
import React, { useEffect, useState, useRef } from "react";
import AOS from "aos";
import "aos/dist/aos.css";

export default function NotificationsDashboard() {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);
  const wsRef = useRef(null);

  useEffect(() => {
    window.scrollTo(0, 0); // <--- add this!
    AOS.init({ duration: 600, once: true });

    fetch("http://localhost:8080/notifications", { credentials: "include" })
      .then(res => res.json())
      .then(data => {
        setNotifications(Array.isArray(data) ? data : []);
        setLoading(false);
      })
      .catch(() => {
        setNotifications([]);
        setLoading(false);
      });

    wsRef.current = new WebSocket("ws://localhost:8080/ws");
    wsRef.current.onopen = () => console.log("WebSocket opened");
    wsRef.current.onclose = () => console.log("WebSocket closed");
    wsRef.current.onerror = (err) => console.log("WebSocket error", err);
    wsRef.current.onmessage = (event) => {
      try {
        const notif = JSON.parse(event.data);
        setNotifications(prev => [notif, ...prev]);
      } catch {
        console.warn("Non-JSON message:", event.data);
      }
    };

    return () => {
      if (wsRef.current) wsRef.current.close();
    };
  }, []);

  useEffect(() => {
    AOS.refresh(); // rerun animation on new notifications
  }, [notifications]);

  return (
    <div className="main-layout">
      <div className="notifications-section">
        <h2 className="notifications-header">Notifications dashboard</h2>
        {loading ? (
          <p>Loading...</p>
        ) : notifications.length === 0 ? (
          <p>No notifications yet.</p>
        ) : (
          <ul className="notifications-list">
            {notifications.map(n => (
              <li key={n.id} className="notification-row" data-aos="fade-up">
                <img
                  className="notif-avatar"
                  src={`https://api.dicebear.com/8.x/initials/svg?seed=${n.sender_nickname || "User"}`}
                  alt="User avatar"
                />
                <div className="notif-body">
                  <span className="notif-username">{n.sender_nickname}</span>
                  {n.type === "like" && <> liked your post.</>}
                  {n.type === "follow" && (
                    <>
                      <span className="notif-action insta-follow">
                        started following you.
                      </span>
                    </>
                  )}
                  {n.type === "follow_request" && <> requested to follow you.</>}
                </div>
                <span className="notif-date">
                  {n.created_at
                    ? new Date(n.created_at).toLocaleString([], { dateStyle: "short", timeStyle: "short" })
                    : ""}
                </span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
