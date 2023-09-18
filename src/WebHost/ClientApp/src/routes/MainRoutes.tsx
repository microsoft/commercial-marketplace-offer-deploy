import MainLayout from '../layout/MainLayout';
import Default from '../views/Default';

const MainRoutes = {
    path: '/',
    element: <MainLayout />,
    children: [
        {
            path: '/',
            element: <Default />
        },
    ]
};

export default MainRoutes;