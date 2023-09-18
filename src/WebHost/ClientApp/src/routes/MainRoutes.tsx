import MainLayout from '../layout/MainLayout';
import Default from '../views/Default/Index';
import { Status } from '../views/Status';

const MainRoutes = {
    path: '/',
    element: <MainLayout />,
    children: [
        {
            path: '/',
            element: <Default />
        },
        {
            path: 'status', // Define the new route path
            element: <Status /> // Use the new component
          },
    ],
};

export default MainRoutes;