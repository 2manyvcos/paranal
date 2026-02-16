import { Navigate, Outlet, useLocation, useSearchParams } from 'react-router';
import { useApplication } from '@/application';

export default function Dashboard() {
  const application = useApplication();
  const location = useLocation();
  const [searchParams] = useSearchParams();

  if (location.pathname === '/' && !searchParams.has('no-redirect')) {
    if (application.user?.startPage && application.user?.startPage !== '/') {
      return <Navigate to={application.user?.startPage} replace />;
    }
  }

  return (
    <div>
      Dashboard
      <Outlet />
    </div>
  );
}
