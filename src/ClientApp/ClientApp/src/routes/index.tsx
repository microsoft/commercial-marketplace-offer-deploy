// routes/index.tsx

import React from 'react';
import { useRoutes, RouteObject } from 'react-router-dom';
import MainRoutes from './MainRoutes';

const ThemeRoutes = () => {
  const routing = useRoutes(MainRoutes as RouteObject[]);
  return routing;
};

export default ThemeRoutes;
