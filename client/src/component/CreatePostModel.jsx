import React, { useState } from "react";
import "../../styles/CreatePostModal.css";

export default function CreatePostModal({ onClose, onSubmit }) {
    const [title, setTitle] = useState("");
    const [content, setContent] = useState("");
    const [image, setImage] = useState(null);
    const [preview, setPreview] = useState(null);
    const [loading, setLoading] = useState(false);

    const handleImageChange = (e) => {
        const file = e.target.files[0];
        setImage(file);
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => setPreview(reader.result);
            reader.readAsDataURL(file);
        } else {
            setPreview(null);
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            const formData = new FormData();
            formData.append('title', title);
            formData.append('content', content);
            if (image) formData.append('image', image);

            const res = await fetch('http://localhost:8080/posts', {
                method: "POST",
                body: formData,
                credentials: "include",
            });

            if (!res.ok) throw new Error("Failed to create post");
            const post = await res.json();
            onSubmit(post);
            onClose();
        } catch (err) {
            alert("Could not create post. Try again.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="modal-overlay">
            <div className="modal">
                <h2>Create New Post</h2>
                <form onSubmit={handleSubmit}>
                    <input
                        type="text"
                        placeholder="Title"
                        value={title}
                        onChange={e => setTitle(e.target.value)}
                        required
                        disabled={loading}
                    />
                    <textarea
                        placeholder="What's on your mind?"
                        value={content}
                        onChange={e => setContent(e.target.value)}
                        required
                        disabled={loading}
                    />

                    <div className="file-section">
                        <input
                            type="file"
                            accept="image/*"
                            onChange={handleImageChange}
                            disabled={loading}
                        />
                        {preview && <img src={preview} alt="Preview" className="image-preview" />}
                    </div>

                    <div className="modal-actions">
                        <button type="submit" disabled={loading}>
                            {loading ? "Posting..." : "Post"}
                        </button>
                        <button type="button" onClick={onClose} disabled={loading}>Cancel</button>
                    </div>
                </form>
            </div>
        </div>
    );
} 