import React, { useState, useEffect } from "react";
import CreatePostButton from "./CreatePostButton";
import { Link } from "react-router-dom";
import AOS from "aos";
import "aos/dist/aos.css";
import likeIcon from "/static/icons/like.svg";
import dislikeIcon from "/static/icons/dislike.svg";
import "../../styles/Home.css";

export default function Home() {
  const [posts, setPosts] = useState([]);

  useEffect(() => {
    AOS.init({ duration: 800, once: true });

    fetch("http://localhost:8080/posts", { credentials: "include" })
      .then(res => res.json())
      .then(data => {
        console.log("Fetched posts:", data);
        if (Array.isArray(data)) {
          setPosts(data);
        } else {
          console.warn("Expected array but got:", data);
          setPosts([]);
        }
      })
      .catch(err => {
        console.error("Failed to fetch posts", err);
        setPosts([]);
      });
  }, []);

  function handleAddPost(post) {
    setPosts(prev => [post, ...prev]);
    AOS.refresh(); // Animate newly added post
  }

  function handleLike(idx, e) {
    e.preventDefault();
    e.stopPropagation();

    const post = posts[idx];
    setPosts(prev =>
      prev.map((p, i) =>
        i === idx ? { ...p, liked: !p.liked, disliked: p.liked ? p.disliked : false } : p
      )
    );

    const formData = new URLSearchParams();
    formData.append("post_id", post.id);

    fetch("http://localhost:8080/like", {
      method: "POST",
      credentials: "include",
      body: formData,
    })
      .then(res => res.json())
      .then(data => {
        setPosts(prev =>
          prev.map((p, i) =>
            i === idx ? { ...p, liked: data.liked, disliked: data.disliked } : p
          )
        );
      })
      .catch(err => {
        console.error("Failed to like/unlike", err);
        setPosts(prev =>
          prev.map((p, i) =>
            i === idx ? { ...p, liked: !p.liked } : p
          )
        );
      });
  }

  function handleDislike(idx, e) {
    e.preventDefault();
    e.stopPropagation();

    const post = posts[idx];
    setPosts(prev =>
      prev.map((p, i) =>
        i === idx ? { ...p, disliked: !p.disliked, liked: p.disliked ? p.liked : false } : p
      )
    );

    const formData = new URLSearchParams();
    formData.append("post_id", post.id);

    fetch("http://localhost:8080/dislike", {
      method: "POST",
      credentials: "include",
      body: formData,
    })
      .then(res => res.json())
      .then(data => {
        setPosts(prev =>
          prev.map((p, i) =>
            i === idx ? { ...p, disliked: data.disliked, liked: data.liked } : p
          )
        );
      })
      .catch(err => {
        console.error("Failed to dislike/undislike", err);
        setPosts(prev =>
          prev.map((p, i) =>
            i === idx ? { ...p, disliked: !p.disliked } : p
          )
        );
      });
  }

  return (
    <div className="main-layout">
      <div className="home-container">
        {Array.isArray(posts) && posts.map((post, idx) => (
          <Link
            to={`/posts/${post.id}`}
            className="post-link"
            key={post.id || idx}
          >
            <div
              className="post-card"
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
                alt="Post visual"
                className="post-image"
              />

              <div className="post-title">{post.title}</div>
              <div className="post-body">{post.content}</div>
              <div className="post-meta">
                Posted by <b>{post.username}</b> • {post.created_at}
              </div>

              <div className="like-dislike-row">
                <button
                  className={`like-btn${post.liked ? " liked" : ""}`}
                  onClick={e => handleLike(idx, e)}
                  aria-label="Like"
                  disabled={post.disliked}
                  style={{
                    opacity: post.disliked ? 0.4 : 1,
                    pointerEvents: post.disliked ? "none" : "auto",
                  }}
                >
                  <svg
                    className="like-icon"
                    width="32"
                    height="32"
                    viewBox="0 0 24 24"
                    fill={post.liked ? "#0099ff" : "#65676b"}
                    style={post.liked ? { filter: "drop-shadow(0 0 12px #b3e0ff)" } : {}}
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M2 21h2V9H2v12zM23 10c0-1.1-.9-2-2-2h-6.31l.95-4.57.03-.32
                      c0-.41-.17-.79-.44-1.06L14.17 2 7.59 8.59C7.22 8.95 7 9.45
                      7 10v9c0 1.1.9 2 2 2h9c.78 0 1.48-.45 1.81-1.13l3.02-6.03
                      C22.94 13.2 23 12.6 23 12V10z" />
                  </svg>
                </button>

                <button
                  className={`dislike-btn${post.disliked ? " disliked" : ""}`}
                  onClick={e => handleDislike(idx, e)}
                  aria-label="Dislike"
                  disabled={post.liked}
                  style={{
                    opacity: post.liked ? 0.4 : 1,
                    pointerEvents: post.liked ? "none" : "auto",
                  }}
                >
                  <svg
                    className="dislike-icon"
                    width="32"
                    height="32"
                    viewBox="0 0 24 24"
                    fill={post.disliked ? "#ff0000" : "#65676b"}
                    style={{
                      transform: "rotate(180deg)",
                      filter: post.disliked ? "drop-shadow(0 0 12px #ffb3b3)" : undefined,
                    }}
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M2 21h2V9H2v12zM23 10c0-1.1-.9-2-2-2h-6.31l.95-4.57.03-.32
                      c0-.41-.17-.79-.44-1.06L14.17 2 7.59 8.59C7.22 8.95 7 9.45
                      7 10v9c0 1.1.9 2 2 2h9c.78 0 1.48-.45 1.81-1.13l3.02-6.03
                      C22.94 13.2 23 12.6 23 12V10z" />
                  </svg>
                </button>
              </div>
            </div>
          </Link>
        ))}
        <CreatePostButton onPost={handleAddPost} />
      </div>
    </div>
  );
}
