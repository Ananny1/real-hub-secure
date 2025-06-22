import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import AOS from "aos";
import "aos/dist/aos.css";
import "../../styles/profile.css";

export default function Profile() {
  const [likedPosts, setLikedPosts] = useState([]);
  const [myPosts, setMyPosts] = useState([]);
  const [username, setUsername] = useState("Guest");
  const [visibility, setVisibility] = useState("private");
  const [stats, setStats] = useState({ postCount: 0, followingCount: 0, followerCount: 0 });

  useEffect(() => {
    AOS.init({ duration: 800 });
  }, []);

  useEffect(() => {
    fetch("http://localhost:8080/profile/user", { credentials: "include" })
      .then(res => res.json())
      .then(data => {
        setUsername(data.username || "Guest");
        setVisibility(data.visibility || "private");
      })
      .catch(() => {
        setUsername("Guest");
        setVisibility("private");
      });

    fetch("http://localhost:8080/profile/liked", { credentials: "include" })
      .then(res => res.json())
      .then(data => setLikedPosts(Array.isArray(data) ? data : []))
      .catch(() => setLikedPosts([]));

    fetch("http://localhost:8080/profile/myposts", { credentials: "include" })
      .then(res => res.json())
      .then(data => setMyPosts(Array.isArray(data) ? data : []))
      .catch(() => setMyPosts([]));

    fetch("http://localhost:8080/profile/stats", { credentials: "include" })
      .then(res => res.json())
      .then(data => setStats(data))
      .catch(() => setStats({ postCount: 0, followingCount: 0, followerCount: 0 }));
  }, []);

  useEffect(() => {
    AOS.refresh();
  }, [likedPosts, myPosts]);

  const updateVisibility = (newValue) => {
    fetch("http://localhost:8080/profile/visibility", {
      method: "PATCH",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ visibility: newValue })
    })
      .then(res => res.json())
      .then(() => setVisibility(newValue))
      .catch(() => alert("Failed to update visibility"));
  };

  return (
    <div className="main-layout">
      <div className="profile-page-container">
        <div className="profile-upper" data-aos="fade-down">
          <img
            className="profile-avatar"
            src="https://api.dicebear.com/8.x/initials/svg?seed=User"
            alt="Avatar"
          />
          <div>
            <h1 className="profile-title">Your Profile</h1>
            <p className="profile-title">Welcome back! {username}</p>

            <div className="profile-stats">
              <div className="profile-stat" data-aos="zoom-in" data-aos-delay="100">
                <span className="profile-stat-number">{stats.postCount}</span>
                <span className="profile-stat-label">Posts</span>
              </div>
              <div className="profile-stat" data-aos="zoom-in" data-aos-delay="200">
                <span className="profile-stat-number">{stats.followingCount}</span>
                <span className="profile-stat-label">Following</span>
              </div>
              <div className="profile-stat" data-aos="zoom-in" data-aos-delay="300">
                <span className="profile-stat-number">{stats.followerCount}</span>
                <span className="profile-stat-label">Followers</span>
              </div>
            </div>

            <div className="visibility-toggle" data-aos="fade-up">
              <p>
                Account visibility:{" "}
                <strong className={visibility === "public" ? "public" : ""}>
                  {visibility.toUpperCase()}
                </strong>
              </p>
              <select
                className="visibility-dropdown"
                value={visibility}
                onChange={(e) => updateVisibility(e.target.value)}
              >
                <option value="private">Private</option>
                <option value="public">Public</option>
              </select>
            </div>
          </div>
        </div>

        <div className="profile-lower">
          <section className="profile-section" data-aos="fade-up">
            <h2>❤️ Liked Posts</h2>
            {likedPosts.length === 0 ? (
              <p className="profile-empty-msg">You haven't liked any posts yet.</p>
            ) : (
              <div className="profile-section-scrollable">
                <div className="profile-post-list">
                  {likedPosts.map((post, idx) => (
                    <Link
                      to={`/posts/${post.id}`}
                      className="profile-post-card"
                      key={post.id}
                      data-aos="fade-up"
                      data-aos-delay={idx * 100}
                    >
                      <img
                        src={
                          post.image
                            ? post.image.startsWith("http")
                              ? post.image
                              : post.image.startsWith("/")
                                ? `http://localhost:8080${post.image}`
                                : `http://localhost:8080/uploads/${post.image}`
                            : "https://images.unsplash.com/photo-1506744038136-46273834b3fb?auto=format&fit=crop&w=800&q=80"
                        }
                        alt="Post"
                        className="profile-post-img"
                      />
                      <div className="profile-post-info">
                        <span className="liked-badge">❤️ Liked Post</span>
                        <div className="profile-post-title">{post.title}</div>
                        <div className="profile-post-date">{post.created_at}</div>
                      </div>
                    </Link>
                  ))}
                </div>
              </div>
            )}
          </section>

          <section className="profile-section" data-aos="fade-up">
            <h2>Your Posts</h2>
            {myPosts.length === 0 ? (
              <div className="profile-empty-state" data-aos="fade-in">
                <img
                  src="https://images.unsplash.com/photo-1587614382346-4ec1c9ff2c02?auto=format&fit=crop&w=800&q=80"
                  alt="No posts"
                  className="profile-empty-img"
                />
                <p className="profile-empty-msg">You haven't posted anything yet.</p>
              </div>
            ) : (

              <div className="profile-section-scrollable">
                <div className="profile-post-list">
                  {myPosts.map((post, idx) => (
                    <Link
                      to={`/posts/${post.id}`}
                      className="profile-post-card"
                      key={post.id}
                      data-aos="fade-up"
                      data-aos-delay={idx * 100}
                    >
                      <img
                        src={
                          post.image
                            ? post.image.startsWith("http")
                              ? post.image
                              : post.image.startsWith("/")
                                ? `http://localhost:8080${post.image}`
                                : `http://localhost:8080/uploads/${post.image}`
                            : "https://images.unsplash.com/photo-1506744038136-46273834b3fb?auto=format&fit=crop&w=800&q=80"
                        }
                        alt="Post"
                        className="profile-post-img"
                      />

                      <div className="profile-post-info">
                        <div className="profile-post-title">{post.title}</div>
                        <div className="profile-post-date">{post.created_at}</div>
                      </div>
                    </Link>
                  ))}
                </div>
              </div>
            )}
          </section>
        </div>
      </div>
    </div>
  );
}
