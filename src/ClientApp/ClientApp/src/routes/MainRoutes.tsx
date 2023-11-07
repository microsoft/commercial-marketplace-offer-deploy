// MainRoutes.ts

import React from 'react';
import { RouteObject } from 'react-router-dom';
import MainLayout from '../layout/MainLayout';
import Default from '../views/Default';
import Diagnostics from '../views/Diagnostics';
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
        path: 'diagnostics',
        element: <ProtectedRoute component={Diagnostics} />,
      },
      {
        path: 'login',
        element: <LoginPage />,
      },
      // ...other routes here...
    ],
  },
];

export default MainRoutes;
