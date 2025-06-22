import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import '../../styles/PostPage.css';

export default function PostPage() {
    const { id } = useParams();
    const [post, setPost] = useState(null);
    const [comments, setComments] = useState([]);
    const [commentText, setCommentText] = useState("");
    const [loadingPost, setLoadingPost] = useState(true);
    const [loadingComments, setLoadingComments] = useState(true);
    const [submitLoading, setSubmitLoading] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetch(`http://localhost:8080/posts/${id}`, { credentials: "include" })
            .then(res => {
                if (!res.ok) throw new Error("Post not found");
                return res.json();
            })
            .then(setPost)
            .catch(err => setError(err.message))
            .finally(() => setLoadingPost(false));

        fetch(`http://localhost:8080/posts/${id}/comments`, { credentials: "include" })
            .then(res => {
                if (!res.ok) throw new Error("Failed to load comments");
                return res.json();
            })
            .then(data => setComments(Array.isArray(data) ? data : []))
            .catch(err => setError(err.message))
            .finally(() => setLoadingComments(false));
    }, [id]);

    const handleAddComment = (e) => {
        e.preventDefault();
        setSubmitLoading(true);

        fetch(`http://localhost:8080/posts/${id}/comments`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ content: commentText })
        })
            .then(res => {
                if (!res.ok) throw new Error("Failed to submit comment");
                return res.json();
            })
            .then(newComment => {
                setComments(prev => [newComment, ...prev]);
                setCommentText("");
            })
            .catch(err => setError(err.message))
            .finally(() => setSubmitLoading(false));
    };

    if (loadingPost) return <div className="loading">Loading post...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <main className="post-page-wrapper">
            <div className="post-card">
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
            </div>

            <div className="comments-section">
                <h3>Comments</h3>
                {loadingComments ? (
                    <div className="loading">Loading comments...</div>
                ) : comments.length === 0 ? (
                    <div className="empty">No comments yet.</div>
                ) : (
                    comments.map(comment => (
                        <div key={comment.id} className="comment-card">
                            <div className="comment-author">{comment.username}</div>
                            <div className="comment-content">{comment.content}</div>
                            <div className="comment-meta">{comment.created_at}</div>
                        </div>
                    ))
                )}

                <form className="add-comment-form" onSubmit={handleAddComment}>
                    <textarea
                        placeholder="Write your comment..."
                        value={commentText}
                        onChange={e => setCommentText(e.target.value)}
                        required
                        disabled={submitLoading}
                    />
                    <button type="submit" disabled={submitLoading || !commentText.trim()}>
                        {submitLoading ? "Posting..." : "Submit"}
                    </button>
                    {error && <div className="comment-error">{error}</div>}
                </form>
            </div>
        </main>
    );
}
