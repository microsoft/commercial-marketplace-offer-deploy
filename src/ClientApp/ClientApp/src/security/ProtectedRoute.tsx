import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from './AuthContext';

interface ProtectedRouteProps {
  component: React.ComponentType<any>;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ component: Component}) => {
    console.log('ProtectedRoute: component: ', Component);
  const { isAuthenticated } = useAuth();
  console.log('ProtectedRoute: isAuthenticated: ', isAuthenticated);
  const location = useLocation();
  console.log('ProtectedRoute: location: ', location);
  return isAuthenticated ? <Component /> : <Navigate to="/login" replace />;
};

export default ProtectedRoute;
