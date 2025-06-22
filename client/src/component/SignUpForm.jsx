import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import AOS from 'aos';
import 'aos/dist/aos.css';
import '../../styles/SignUpForm.css';

function SignupForm() {
  const [form, setForm] = useState({
    nickname: '',
    first_name: '',
    last_name: '',
    email: '',
    password: '',
    gender: '',
    age: ''
  });

  const [message, setMessage] = useState('');
  const navigate = useNavigate();
  const { signup } = useAuth();

  useEffect(() => {
    AOS.init({ duration: 800, once: true });
  }, []);

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const payload = {
      ...form,
      age: form.age === '' ? undefined : parseInt(form.age)
    };

    const result = await signup(payload);

    if (result.success) {
      navigate('/');
    } else {
      setMessage(result.message || 'Signup failed');
    }
  };

  return (
    <div className="signup-wrapper">
      <div className="form-container" data-aos="fade-up">
        <h2 data-aos="fade-down">Sign Up</h2>
        <form onSubmit={handleSubmit}>
          <input name="nickname" placeholder="Nickname" onChange={handleChange} required data-aos="fade-right" />
          <input name="first_name" placeholder="First Name" onChange={handleChange} required data-aos="fade-left" />
          <input name="last_name" placeholder="Last Name" onChange={handleChange} required data-aos="fade-right" />
          <input type="email" name="email" placeholder="Email" onChange={handleChange} required data-aos="fade-left" />
          <input type="password" name="password" placeholder="Password" onChange={handleChange} required data-aos="fade-up" />
          <select name="gender" onChange={handleChange} required data-aos="fade-up">
            <option value="">Select Gender</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
          </select>
          <input type="number" name="age" placeholder="Age" min="13" onChange={handleChange} required data-aos="fade-up" />
          <button type="submit" data-aos="zoom-in">Register</button>
        </form>
        {message && <p className="error-msg" data-aos="fade-in">{message}</p>}
        <p data-aos="fade-in">
          Already have an account? <Link to="/login">Log In</Link>
        </p>
      </div>
    </div>
  );
}

export default SignupForm;
