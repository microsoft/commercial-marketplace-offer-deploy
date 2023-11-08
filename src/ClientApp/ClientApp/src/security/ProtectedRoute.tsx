import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from './AuthContext';

interface ProtectedRouteProps {
  component: React.ComponentType<any>;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ component: Component}) => {
  const { isAuthenticated } = useAuth();
  const location = useLocation();
  return isAuthenticated ? <Component /> : <Navigate to="/login" replace />;
};

export default ProtectedRoute;
