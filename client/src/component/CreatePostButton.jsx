import React, { useState } from "react";
import CreatePostModal from "./CreatePostModel";
import "../../styles/CreatePostButton.css";

export default function CreatePostButton({ onPost }) {
  const [showModal, setShowModal] = useState(false);

  return (
    <>
      <button
        className="fab-create-post"
        onClick={() => setShowModal(true)}
        aria-label="Create Post"
      >
        <span className="fab-icon">+</span>
      </button>

      {showModal && (
        <CreatePostModal
          onClose={() => setShowModal(false)}
          onSubmit={onPost}
        />
      )}
    </>
  );
}
