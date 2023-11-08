import React, { useEffect, useState } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from './AuthContext';

interface ProtectedRouteProps {
  component: React.ComponentType<any>;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ component: Component}) => {
  const { isAuthenticated, checkAuth } = useAuth();
  const [isCheckingAuth, setIsCheckingAuth] = useState(true);
  console.log('ProtectedRoute - isAuthenticated: ', isAuthenticated);
  const location = useLocation();

  useEffect(() => {
    const verifyAuth = async () => {
        await checkAuth();
        console.log('ProtectedRoute - isAuthenticated updated: ', isAuthenticated);
        setIsCheckingAuth(false);
      };
  
      verifyAuth();
  }, []);

  if (isCheckingAuth) {
    return <div>Loading...</div>;
  }

  return isAuthenticated ? <Component /> : <Navigate to="/login" replace />;
};

export default ProtectedRoute;
