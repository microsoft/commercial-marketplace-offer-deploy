import { useRoutes } from 'react-router-dom';

// routes
import MainRoutes from './MainRoutes';

export default function ThemeRoutes() {
    return useRoutes([MainRoutes]);
}
