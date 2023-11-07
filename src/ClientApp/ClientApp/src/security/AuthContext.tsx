// AuthContext.tsx

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { validateToken, AuthToken, login, refreshToken } from './securityutils'; // Ensure refreshToken is defined and exported

interface AuthContextType {
  isAuthenticated: boolean;
  userToken: AuthToken | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType>(null!); // Non-null assertion

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [userToken, setUserToken] = useState<AuthToken | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    console.log('AuthProvider: useEffect');
    const initializeAuth = async () => {
        console.log('AuthProvider: initializeAuth');
      const tokenString = localStorage.getItem('jwtToken');
      console.log(`AuthProvider: tokenString: ${tokenString}`);
      if (tokenString) {
        console.log('AuthProvider: tokenString is not null');   
        const token: AuthToken = JSON.parse(tokenString);
        console.log(`AuthProvider: token: ${token}`);
        if (validateToken()) {
            console.log('AuthProvider: token is valid');
          setIsAuthenticated(true);
          console.log(`setIsAuthenticated(true) called`);
          setUserToken(token);
        } else {
          // Token is expired or invalid, try to refresh it
          try {
            console.log('validatedToken: token is invalid');
            const newToken = await refreshToken(token.id); // Assume refreshToken is a function that you will implement
            console.log('validatedToken: newToken: ', newToken);
            localStorage.setItem('jwtToken', JSON.stringify(newToken));
            console.log('validatedToken: localStorage.setItem');
            setIsAuthenticated(true);
            console.log('validatedToken: setIsAuthenticated(true)');
            setUserToken(newToken);
            console.log('validatedToken: setUserToken(newToken)');
          } catch (error) {
            // Refresh token failed
            console.log('validatedToken: refreshToken failed');
            localStorage.removeItem('jwtToken');
            console.log('validatedToken: localStorage.removeItem');
            setIsAuthenticated(false);
            console.log('validatedToken: setIsAuthenticated(false)');
            setUserToken(null);
            console.log('validatedToken: setUserToken(null)');
          }
        }
      }
    };

    initializeAuth();
  }, []);

  const handleLogin = async (username: string, password: string) => {
    console.log('AuthProvider: handleLogin');
    const token = await login(username, password);
    console.log('AuthProvider: handleLogin: token: ', token);
    localStorage.setItem('jwtToken', JSON.stringify(token));
    console.log('AuthProvider: handleLogin: localStorage.setItem');
    setIsAuthenticated(true);
    console.log('AuthProvider: handleLogin: setIsAuthenticated(true)');
    setUserToken(token);
    console.log('AuthProvider: handleLogin: setUserToken(token)');
  };

  const handleLogout = () => {
    console.log('AuthProvider: handleLogout');
    localStorage.removeItem('jwtToken');
    console.log ('AuthProvider: handleLogout: localStorage.removeItem');
    setIsAuthenticated(false);
    console.log('AuthProvider: handleLogout: setIsAuthenticated(false)');
    setUserToken(null);
    console.log('AuthProvider: handleLogout: setUserToken(null)');
  };

  // Value provided to context consumers
  const value = {
    isAuthenticated,
    userToken,
    login: handleLogin,
    logout: handleLogout,
  };

  console.log('AuthProvider: value: ', value);
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
