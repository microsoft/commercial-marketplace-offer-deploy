// ProtectedRoute.tsx

import React from 'react';
import { Navigate } from 'react-router-dom';
import { validateToken } from './securityutils';

interface ProtectedRouteProps {
  component: React.ComponentType<any>;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ component: Component }) => {
  const isAuthenticated = validateToken();

  return isAuthenticated ? <Component /> : <Navigate to="/login" replace />;
};

export default ProtectedRoute;
