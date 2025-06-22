import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import '../../styles/NavBar.css';

export default function NavBar() {
  const navigate = useNavigate();
  const { logout } = useAuth();

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const navItems = [
    { label: 'Home', icon: '/static/icons/Home.svg', route: '/' },
    { label: 'Profile', icon: '/static/icons/Profile.svg', route: '/profile' },
    { label: 'Chat', icon: '/static/icons/chat.svg', route: '/chat' },
    { label: 'Search', icon: '/static/icons/search.svg', route: '/searchforusers' },
    { label: 'Notifications', icon: '/static/icons/follow.svg', route: '/notifications' },
    { label: 'About', icon: '/static/icons/about.svg', route: '/about' },
    { label: 'Logout', icon: '/static/icons/logout.svg', action: handleLogout },
  ];

  return (
    <nav className="side-nav-container">
      <div className="side-nav-logo">
        <span className="brand-part-1">Real Hub</span>
        <span className="brand-part-2">/ Secure /</span>
      </div>

      <div className="side-nav-items">
        {navItems.map((item, index) => (
          <div
            key={index}
            className="side-nav-item"
            onClick={item.action ? item.action : () => navigate(item.route)}
          >
            <img src={item.icon} alt={item.label} />
            <span>{item.label}</span>
          </div>
        ))}
      </div>
    </nav>
  );
}
