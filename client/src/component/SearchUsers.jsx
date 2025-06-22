import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import AOS from "aos";
import "aos/dist/aos.css";
import "../../styles/searchusers.css";

export default function SearchUsers() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState([]);
  const [searching, setSearching] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    AOS.init({ duration: 600, once: true });
  }, []);

  useEffect(() => {
    if (!query.trim()) {
      setResults([]);
      setError("");
      return;
    }

    const delayDebounce = setTimeout(() => {
      setSearching(true);
      fetch(`http://localhost:8080/users/search?query=${encodeURIComponent(query)}`, {
        credentials: "include",
      })
        .then(res => {
          if (!res.ok) throw new Error("Failed to fetch");
          return res.json();
        })
        .then(data => {
          setResults(data);
          setError("");
        })
        .catch(() => setError("Could not load users"))
        .finally(() => {
          setSearching(false);
          AOS.refresh(); // trigger animations on result render
        });
    }, 300);

    return () => clearTimeout(delayDebounce);
  }, [query]);

  return (
    <div className="main-layout">
      <div className="search-page-centered">
        <h1 className="search-title" data-aos="fade-down">Search Users</h1>
        <div className="search-bar" data-aos="fade-up">
          <input
            className="search-input"
            type="text"
            placeholder="Search by nickname or email..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>

        <div className="search-results">
          {searching && <p className="search-status">Searching...</p>}
          {error && <p className="search-error">{error}</p>}
          {!searching && !error && results.length === 0 && query && (
            <p className="search-status">No users found.</p>
          )}
          {results.map((user, idx) => (
            <Link
              to={`/users/${user.id}`}
              className="search-user-card"
              key={user.id}
              data-aos="fade-up"
              data-aos-delay={idx * 100} // staggered delay
            >
              <img
                src={`https://api.dicebear.com/8.x/initials/svg?seed=${user.nickname}`}
                alt="avatar"
                className="user-avatar"
              />
              <span className="user-nickname">{user.nickname}</span>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}
