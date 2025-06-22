import { useState, useRef } from "react";
import imageIcon from "../../../static/icons/image.svg";
import './chat.css';

const MessageInput = ({ onSend, onSendImage }) => {
    const [message, setMessage] = useState("");
    const [image, setImage] = useState(null);
    const fileInputRef = useRef(null);

    const handleSubmit = (e) => {
        e.preventDefault();
        if (message.trim()) {
            onSend(message);
            setMessage("");
        }
        if (image) {
            handleImageUpload(image);
            setImage(null);
            fileInputRef.current.value = "";
        }
    };

    const handleImageChange = (e) => {
        const file = e.target.files[0];
        if (file) setImage(file);
    };

    const handleImageUpload = async (file) => {
        const formData = new FormData();
        formData.append("image", file);
        try {
            const res = await fetch("http://localhost:8080/chat/upload", {
                method: "POST",
                body: formData,
                credentials: "include",
            });
            const data = await res.json();
            if (data.url && onSendImage) onSendImage(data.url);
        } catch (err) {
            alert("Image upload failed");
        }
    };

    return (
        <form
            onSubmit={handleSubmit}
            className="message-input-form"
            encType="multipart/form-data"
            autoComplete="off"
        >
            <input
                type="text"
                value={message}
                onChange={e => setMessage(e.target.value)}
                placeholder="Type a message..."
                className="message-text-input"
            />
            {/* SVG as file picker */}
            <button
                type="button"
                className={`input-btn${image ? " has-image" : ""}`}
                onClick={() => fileInputRef.current && fileInputRef.current.click()}
                tabIndex={-1}
                aria-label="Attach image"
                style={{ padding: 0, marginLeft: 10, background: "none", border: "none" }}
            >
                <img
                    src={imageIcon}
                    alt="Attach"
                    width={26}
                    height={26}
                    style={{
                        filter: image
                            ? "brightness(0) invert(1) drop-shadow(0 0 4px #0095f6) saturate(2)"
                            : "brightness(0) invert(1)",
                        transition: "filter 0.18s",
                        display: "block",
                        cursor: "pointer"
                    }}
                />
            </button>
            <input
                type="file"
                accept="image/*"
                ref={fileInputRef}
                onChange={handleImageChange}
                style={{ display: "none" }}
            />
            <button type="submit" className="send-btn">
                Send
            </button>
        </form>
    );
};

export default MessageInput;
