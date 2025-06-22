// src/components/ConfirmPopup.jsx
import React from "react";
import "../../styles/confirmPopup.css";

export default function ConfirmPopup({ message, onConfirm, onCancel }) {
  return (
    <div className="popup-overlay">
      <div className="popup-box">
        <p className="popup-message">{message}</p>
        <div className="popup-buttons">
          <button className="popup-btn confirm" onClick={onConfirm}>Yes</button>
          <button className="popup-btn cancel" onClick={onCancel}>Cancel</button>
        </div>
      </div>
    </div>
  );
}
