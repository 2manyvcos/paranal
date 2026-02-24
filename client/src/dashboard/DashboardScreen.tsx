import { Navigate, Outlet, useLocation, useSearchParams } from 'react-router';
import { useApplication } from '@/application';

export default function DashboardScreen() {
  const application = useApplication();
  const location = useLocation();
  const [searchParams] = useSearchParams();

  if (location.pathname === '/' && !searchParams.has('no-redirect')) {
    if (application.user?.startPage && application.user?.startPage !== '/') {
      return <Navigate to={application.user?.startPage} replace />;
    }
  }

  return (
    <div className="dashboard screen">
      Dashboard
      <Outlet />
    </div>
  );
}
