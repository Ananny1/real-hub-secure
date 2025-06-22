const ChatBox = ({ messages }) => {
    return (
        <div className="chat-box">
            {messages.map((msg, index) => {
                const hasImage = msg.imageUrl && msg.imageUrl.trim();
                const hasText = msg.text && msg.text.trim();
                const isImageOnly = hasImage && !hasText;

                return (
                    <div
                        key={index}
                        className={`chat-message ${msg.from === "" ? "sent" : "received"} ${isImageOnly ? "image-only" : ""}`}
                    >
                        <strong>{msg.from}</strong>{" "}
                        <>
                            {hasImage && (
                                <img
                                    src={msg.imageUrl}
                                    alt="Shared image"
                                    className="chat-image"
                                    onClick={() => {
                                        // Optional: Open image in new tab when clicked
                                        window.open(msg.imageUrl, '_blank');
                                    }}
                                />
                            )}
                            {hasText && (
                                <div className="message-text" style={{ marginTop: hasImage ? 6 : 0 }}>
                                    {msg.text}
                                </div>
                            )}
                        </>
                    </div>
                );
            })}
        </div>
    );
};

export default ChatBox;