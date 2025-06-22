// AuthContext.js
import React, { createContext, useContext, useState, useEffect } from 'react';

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    useEffect(() => {
        validateSession();
    }, []);

    const validateSession = async () => {
        try {
            const response = await fetch('http://localhost:8080/validate-session', {
                credentials: 'include'
            });
            
            if (response.ok) {
                const data = await response.json();
                setUser(data.user);
                setIsAuthenticated(true);
                localStorage.setItem('auth', 'true');
            } else {
                setUser(null);
                setIsAuthenticated(false);
                localStorage.removeItem('auth');
            }
        } catch (error) {
            console.error('Session validation error:', error);
            setUser(null);
            setIsAuthenticated(false);
            localStorage.removeItem('auth');
        } finally {
            setLoading(false);
        }
    };

    const login = async (email, password) => {
        try {
            const response = await fetch('http://localhost:8080/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password }),
                credentials: 'include'
            });

            if (response.ok) {
                await validateSession(); // Refresh user data
                return { success: true };
            } else {
                const data = await response.json();
                return { success: false, message: data.message };
            }
        } catch (error) {
            return { success: false, message: 'Login failed' };
        }
    };

    const signup = async (userData) => {
        try {
            const response = await fetch('http://localhost:8080/signup', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(userData),
                credentials: 'include'
            });

            if (response.ok) {
                await validateSession(); // Refresh user data
                return { success: true };
            } else {
                const data = await response.json();
                return { success: false, message: data.message };
            }
        } catch (error) {
            return { success: false, message: 'Signup failed' };
        }
    };

    const logout = async () => {
        try {
            await fetch('http://localhost:8080/logout', {
                method: 'POST',
                credentials: 'include'
            });
        } catch (error) {
            console.error('Logout error:', error);
        } finally {
            setUser(null);
            setIsAuthenticated(false);
            localStorage.removeItem('auth');
        }
    };

    return (
        <AuthContext.Provider value={{ 
            user, 
            isAuthenticated,
            loading,
            login, 
            signup,
            logout,
            validateSession 
        }}>
            {children}
        </AuthContext.Provider>
    );
};