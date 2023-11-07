import { Redeploy } from '@/views/Redeploy';
import MainLayout from '../layout/MainLayout';
import Default from '../views/Default';
import Diagnostics from '../views/Diagnostics';

const MainRoutes = {
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