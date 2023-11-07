// AuthContext.tsx

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { validateToken, AuthToken, login } from './securityutils';

interface AuthContextType {
  isAuthenticated: boolean;
  userToken: AuthToken | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType>(null!); // The `!` is for non-null assertion.

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [userToken, setUserToken] = useState<AuthToken | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const tokenString = localStorage.getItem('jwtToken');
    if (tokenString) {
      const token = JSON.parse(tokenString);
      // this is where the token refresh should be performned if expired - call the api to refresh the token
      if (validateToken()) {
        setUserToken(token);
        setIsAuthenticated(true);
      } else {
        localStorage.removeItem('jwtToken'); 
        setUserToken(null);
        setIsAuthenticated(false);
      }
    }
  }, []);

  const handleLogin = async (username: string, password: string) => {
    try {
      const token = await login(username, password);
      localStorage.setItem('jwtToken', JSON.stringify(token));
      setUserToken(token);
      setIsAuthenticated(true);
    } catch (error) {
      // Handle login error
      throw error;
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('jwtToken');
    setUserToken(null);
    setIsAuthenticated(false);
  };

  const value = {
    isAuthenticated,
    userToken,
    login: handleLogin,
    logout: handleLogout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
