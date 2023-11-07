import MainLayout from '../layout/MainLayout';
import Default from '../views/Default';
import Diagnostics from '../views/Diagnostics';
import LoginPage from '../views/LoginPage';
import ProtectedRoute from '../security/ProtectedRoute';

const MainRoutes = {
    path: '/',
    element: <MainLayout />,
    children: [
        {
            path: '/',
            title: 'Dashboard',
            element: <ProtectedRoute component={Default} />
        },
        {
            path: '/diagnostics',
            title: 'Diagnostics',
            element: <ProtectedRoute component={Diagnostics} />
        },
        {
            path: '/login',
            element: <LoginPage />, 
            title: 'Login'
          }
    ]
};

export default MainRoutes;