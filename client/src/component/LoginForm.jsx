import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import AOS from 'aos';
import 'aos/dist/aos.css';
import '../../styles/LoginForm.css';

function LoginForm() {
  const [form, setForm] = useState({
    email: '',
    password: ''
  });

  const [message, setMessage] = useState('');
  const navigate = useNavigate();
  const { login } = useAuth();

  useEffect(() => {
    AOS.init({ duration: 800, once: true });
  }, []);

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const result = await login(form.email, form.password);

    if (result.success) {
      navigate('/');
    } else {
      setMessage(result.message || 'Login failed');
    }
  };

  return (
    <div className="login-wrapper">
      <div className="login-container" data-aos="zoom-in">
        <h2 data-aos="fade-down">Login</h2>
        <form onSubmit={handleSubmit}>
          <input
            type="email"
            name="email"
            placeholder="Email"
            onChange={handleChange}
            required
            data-aos="fade-right"
          />
          <input
            type="password"
            name="password"
            placeholder="Password"
            onChange={handleChange}
            required
            data-aos="fade-left"
          />
          <button type="submit" data-aos="zoom-in-up">Log In</button>
        </form>
        {message && <p className="error-msg" data-aos="fade-in">{message}</p>}
        <p data-aos="fade-in">
          Don't have an account? <Link to="/signup">Sign Up</Link>
        </p>
      </div>
    </div>
  );
}

export default LoginForm;
