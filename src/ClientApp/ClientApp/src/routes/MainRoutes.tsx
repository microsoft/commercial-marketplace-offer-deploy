// MainRoutes.ts

import React from 'react';
import { RouteObject } from 'react-router-dom';
import MainLayout from '../layout/MainLayout';
import Default from '../views/Default';
import Diagnostics from '../views/Diagnostics';
import Redeploy from '@/views/Redeploy';
import LoginPage from '../views/LoginPage';
import ProtectedRoute from '../security/ProtectedRoute';

const MainRoutes: RouteObject[] = [
  {
    path: '/',
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <ProtectedRoute component={Default} />,
      },
      {
        path: 'redeploy',
        element: <ProtectedRoute component={Redeploy} />,
      },
      {
        path: 'diagnostics',
        element: <ProtectedRoute component={Diagnostics} />,
      },
    ],
  },
  {
    path: 'login',
    element: <LoginPage />,
  },
];

export default MainRoutes;
