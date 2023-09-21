import MainLayout from '../layout/MainLayout';
import Default from '../views/Default';
import Diagnostics from '../views/Diagnostics';

const MainRoutes = {
    path: '/',
    element: <MainLayout />,
    children: [
        {
            path: '/',
            element: <Default />
        },
        {
            path: '/diagnostics',
            element: <Diagnostics />
        },
    ]
};

export default MainRoutes;