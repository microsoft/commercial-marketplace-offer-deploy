import { Redeploy } from '@/views/Redeploy';
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
            path: '/',
            title: 'Dashboard',
            element: <Default />
        },
        {
            path: '/redeploy',
            title: 'Redeploy',
            element: <Redeploy />
        },
        {
            path: '/diagnostics',
            title: 'Diagnostics',
            element: <Diagnostics />
        },
    ]
};

export default MainRoutes;
