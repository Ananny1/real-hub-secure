import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import ConfirmPopup from "./ConfirmPopup";
import "../../styles/profile.css";

export default function PublicProfile() {
  const { id } = useParams();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [followStatus, setFollowStatus] = useState("none");
  const [showUnfollowConfirm, setShowUnfollowConfirm] = useState(false);
  const [canViewPosts, setCanViewPosts] = useState(false);

  useEffect(() => {
    // Fetch user public info
    fetch(`http://localhost:8080/users/${id}`, { credentials: "include" })
      .then(res => {
        if (!res.ok) throw new Error("User not found");
        return res.json();
      })
      .then(data => {
        setUser(data);
        setLoading(false);
        setCanViewPosts(data.visibility === "public"); // default view
      })
      .catch(() => {
        setUser(null);
        setLoading(false);
      });

    // Fetch current user's follow status with this user
    fetch(`http://localhost:8080/follow/status/${id}`, { credentials: "include" })
      .then(res => res.json())
      .then(data => {
        setFollowStatus(data.status || "none");
        if (data.status === "accepted") setCanViewPosts(true);
      })
      .catch(() => setFollowStatus("none"));
  }, [id]);

  const sendFollowRequest = () => {
    fetch(`http://localhost:8080/follow/${id}`, {
      method: "POST",
      credentials: "include",
    })
      .then(res => res.json())
      .then(data => {
        alert(data.message);
        // If user is public, follow auto-accepted
        if (data.message === "Followed successfully" || data.status === "accepted") {
          setFollowStatus("accepted");
          setCanViewPosts(true);
        } else {
          setFollowStatus("pending");
        }
      })
      .catch(() => alert("Failed to send follow request"));
  };

  const handleUnfollow = () => {
    fetch(`http://localhost:8080/follow/${id}`, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json"
      }
    })
      .then(res => res.json())
      .then(data => {
        alert(data.message);
        setFollowStatus("none");
        setCanViewPosts(user?.visibility === "public");
        setShowUnfollowConfirm(false);
      })
      .catch(() => alert("Failed to unfollow"));
  };

  if (loading) return <div className="main-layout"><div className="profile-page-container">Loading...</div></div>;
  if (!user) return <div className="main-layout"><div className="profile-page-container">User not found.</div></div>;

  return (
    <div className="main-layout">
      <div className="profile-page-container">
        <div className="profile-upper">
          <img
            className="profile-avatar"
            src={`https://api.dicebear.com/8.x/initials/svg?seed=${user.nickname}`}
            alt="Avatar"
          />
          <div>
            <h1 className="profile-title">{user.nickname}'s Profile</h1>
            <p className="profile-title">{user.first_name} {user.last_name}</p>
            <p style={{ fontSize: "14px", opacity: 0.6 }}>Visibility: {user.visibility}</p>

            <div className="profile-stats">
              <div className="profile-stat">
                <span className="profile-stat-number">{user.post_count}</span>
                <span className="profile-stat-label">Posts</span>
              </div>
              <div className="profile-stat">
                <span className="profile-stat-number">{user.following_count}</span>
                <span className="profile-stat-label">Following</span>
              </div>
              <div className="profile-stat">
                <span className="profile-stat-number">{user.follower_count}</span>
                <span className="profile-stat-label">Followers</span>
              </div>
            </div>

{/* Follow Actions */}
{followStatus === "none" && (
  <button onClick={sendFollowRequest} className="follow-button">
    {user.visibility === "public" ? "Follow" : "Send Follow Request"}
  </button>
)}
{followStatus === "pending" && (
  <button className="follow-button disabled" disabled>Request Pending</button>
)}
{followStatus === "accepted" && (
  <button className="follow-button" onClick={() => setShowUnfollowConfirm(true)}>Following</button>
)}

          </div>
        </div>

        {/* Unfollow Confirmation */}
        {showUnfollowConfirm && (
          <ConfirmPopup
            message={`Are you sure you want to unfollow ${user.nickname}?`}
            onConfirm={handleUnfollow}
            onCancel={() => setShowUnfollowConfirm(false)}
          />
        )}

        {/* Posts */}
        <div className="profile-section">
          <h2 style={{ paddingLeft: "24px", paddingTop: "16px" }}>{user.nickname}'s Recent Posts</h2>
          <div className="profile-section-scrollable">
            <div className="profile-post-list">
              {!canViewPosts ? (
                <p className="profile-empty-msg">This account is private. Follow to see their posts.</p>
              ) : !Array.isArray(user.posts) || user.posts.length === 0 ? (
                <p className="profile-empty-msg">No posts yet.</p>
              ) : (
                user.posts.map(post => (
                  <div className="profile-post-card" key={post.id}>
                    <img
                      src={
                        post.image?.startsWith("http")
                          ? post.image
                          : `http://localhost:8080/uploads/${post.image}`
                      }
                      alt="Post"
                      className="profile-post-img"
                    />
                    <div className="profile-post-info">
                      <div className="profile-post-title">{post.title}</div>
                      <div className="profile-post-date">{post.created_at}</div>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
