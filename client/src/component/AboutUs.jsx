import React, { useEffect } from 'react';
import AOS from 'aos';
import 'aos/dist/aos.css';
import "../../styles/About.css";

export default function About() {
  useEffect(() => {
    window.scrollTo(0, 0); // <--- add this!
    AOS.init({ duration: 1000 });
  }, []);

  return (
    <div className="about-page">
      <header className="hero">
        <h1 data-aos="fade-down">Real Hub Secure</h1>
        <p data-aos="fade-up">The ASCII-powered universe of posts, people, and pure expression.</p>
      </header>

      <section data-aos="fade-up" className="section">
        <h2 className="highlight">⚡ What Is This?</h2>
        <p>
          <strong>Real Hub Secure</strong> began as a bold ASCII art engine — a utility that turned text into visual noise (the good kind).
          But now? It's a platform. A digital space where expression meets connection.
        </p>
      </section>

      <section data-aos="fade-up" className="section">
        <h2 className="highlight">📝 Post Like a Pro</h2>
        <p>
          Whether it’s a thought, a quote, or a full-on manifesto, Real Hub Secure lets you post in style. Customize it. Color it.
          Shape it into something that speaks louder than plain text ever could. Then let others feel it.
        </p>
      </section>

      <section data-aos="fade-up" className="section">
        <h2 className="highlight">🔎 Explore. Search. Find.</h2>
        <p>
          Want to see what others are creating? Just scroll. Looking for someone specific? Use our smooth live-search to find users,
          usernames, or tagged ideas instantly. It’s fast, real-time, and always one keystroke ahead.
        </p>
      </section>

      <section data-aos="fade-up" className="section">
        <h2 className="highlight">💬 Real-Time Chat</h2>
        <p>
          It's not a platform if you can't talk. Slide into someone's messages, respond to comments, or start a conversation — all in real time.
          Built with blazing-fast WebSocket communication so you never miss a beat. Or a beep.
        </p>
      </section>

      <section data-aos="fade-up" className="section">
        <h2 className="highlight">🔔 Notifications That Matter</h2>
        <p>
          Like your post blew up? Someone mentioned you in their ASCII monologue? You’ll know instantly. Our friendly, minimal notification system 
          keeps you connected — without overwhelming your zen.
        </p>
      </section>

      <section data-aos="fade-up" className="section">
        <h2 className="highlight">🎯 Friendly. Intuitive. Yours.</h2>
        <p>
          No clutter. No chaos. Just a beautiful dark-mode interface, responsive layouts, and gentle animations.
          Real Hub Secure is designed to be your new favorite digital place — part creative tool, part social haven.
        </p>
      </section>

      <footer>
        <p>Built for the expressive. Fueled by code and caffeine. © 2025 Real Hub Secure.</p>
      </footer>
    </div>
  );
}
