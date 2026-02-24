import { Navigate, Outlet, useLocation, useSearchParams } from 'react-router';
import { useApplication } from '@/application';
import Footer from './Footer';
import Header from './Header';

export default function DashboardScreen() {
  const application = useApplication();
  const location = useLocation();
  const [searchParams] = useSearchParams();

  return (
    <div className="dashboard screen">
      {location.pathname === '/' &&
      !searchParams.has('no-redirect') &&
      application.user?.startPage &&
      application.user.startPage !== '/' ? (
        <Navigate to={application.user?.startPage} replace />
      ) : null}

      <Header />

      <Outlet />

      <Footer />
    </div>
  );
}
