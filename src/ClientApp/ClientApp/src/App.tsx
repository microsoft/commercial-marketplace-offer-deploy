import { useRoutes } from 'react-router-dom';
import MainRoutes from './routes/MainRoutes';

export default function App() {

  const routes = useRoutes(MainRoutes);

  return (
    <>
      {routes} 
    </>
  );
}
