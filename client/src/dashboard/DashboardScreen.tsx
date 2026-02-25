import { Navigate, Outlet, useLocation, useSearchParams } from 'react-router';
import { useApplication } from '@/application';
import Footer from './Footer';
import Navigation from './Navigation';

export default function DashboardScreen() {
  const { user } = useApplication();
  const location = useLocation();
  const [searchParams] = useSearchParams();

  return (
    <div className="dashboard screen">
      {location.pathname === '/' &&
      !searchParams.has('no-redirect') &&
      user?.startPage &&
      user.startPage !== '/' ? (
        <Navigate to={user?.startPage} replace />
      ) : null}

      <Navigation />

      <Outlet />

      <Footer />
    </div>
  );
}
