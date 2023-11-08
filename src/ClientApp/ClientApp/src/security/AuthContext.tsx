// AuthContext.tsx

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { validateToken, AuthToken, login, refreshToken } from './securityUtils'; // Ensure refreshToken is defined and exported

interface AuthContextType {
  isAuthenticated: boolean;
  userToken: AuthToken | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  checkAuth: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>(null!); // Non-null assertion

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [userToken, setUserToken] = useState<AuthToken | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const checkAuth = async () => {
    const tokenString = localStorage.getItem('jwtToken');
    if (tokenString) {
      const token: AuthToken = JSON.parse(tokenString);
      if (validateToken()) {
        setIsAuthenticated(true);
        setUserToken(token);
      } else {
        try {
          const newToken = await refreshToken(token.id);
          localStorage.setItem('jwtToken', JSON.stringify(newToken));
          setIsAuthenticated(true);
          setUserToken(newToken);
        } catch (error) {
          localStorage.removeItem('jwtToken');
          setIsAuthenticated(false);
          setUserToken(null);
        }
      } 
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);

  const handleLogin = async (username: string, password: string) => {
    const token = await login(username, password);
    localStorage.setItem('jwtToken', JSON.stringify(token));
    setIsAuthenticated(true);
    setUserToken(token);
  };

  const handleLogout = () => {
    localStorage.removeItem('jwtToken');
    setIsAuthenticated(false);
    setUserToken(null);
  };

  // Value provided to context consumers
  const value = {
    isAuthenticated,
    userToken,
    login: handleLogin,
    logout: handleLogout,
    checkAuth,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
