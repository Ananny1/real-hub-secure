import { Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from '../contexts/AuthContext';
import SignupForm from './component/SignUpForm';
import LoginForm from './component/LoginForm';
import Home from './component/Home';
import PostPage from './component/PostPage';
import NavBar from './component/NavBar';
import Profile from './component/ProfilePage';
import SearchUsers from './component/SearchUsers';
import PublicProfile from './component/PublicProfile';
import Chat from './component/chat/Chat';
import ProtectedRoute from './component/ProtectedRoute';
import { useAuth } from '../contexts/AuthContext';
import NotificationsDashboard from './component/Notifications';
import About from './component/AboutUs';


// Create a separate component for the app content to use the auth context
function AppContent() {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return <div className="flex items-center justify-center min-h-screen">Loading...</div>;
  }

  return (
    <>
      {isAuthenticated && <NavBar />}

      <div className={isAuthenticated ? 'pt-20' : ''}>
        <Routes>
          <Route
            path="/signup"
            element={
              isAuthenticated
                ? <Navigate to="/" replace />
                : <SignupForm />
            }
          />
          <Route
            path="/login"
            element={
              isAuthenticated
                ? <Navigate to="/" replace />
                : <LoginForm />
            }
          />
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <Home />
              </ProtectedRoute>
            }
          />
          <Route
            path="/posts/:id"
            element={
              <ProtectedRoute>
                <PostPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/notifications"
            element={
              <ProtectedRoute>
                <NotificationsDashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <Profile />
              </ProtectedRoute>
            }
          />
          <Route
            path="/searchforusers"
            element={
              <ProtectedRoute>
                <SearchUsers />
              </ProtectedRoute>
            }
          />
          <Route
            path="/chat"
            element={
              <ProtectedRoute>
                <Chat />
              </ProtectedRoute>
            }
          />

          <Route
            path="/logout"
            element={
              <ProtectedRoute>
                <SignupForm />
              </ProtectedRoute>
            }
          />
          <Route path="/about" element={<About />} />



          <Route
            path="/users/:id"
            element={
              <ProtectedRoute>
                <PublicProfile />
              </ProtectedRoute>
            }
          />

        </Routes>
      </div>
    </>
  );
}

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;

// ✅ Window 1: User logs in → localStorage = 'true' → "I'm logged in"
// ❌ Window 2: Opens new tab → Reads localStorage = 'true' → "I'm also logged in!" (BUT NEVER ASKED SERVER!)