import { Navigate, useSearchParams } from 'react-router';
import { useApplication } from './application';

export default function LoginPage() {
  const application = useApplication();
  const [searchParams] = useSearchParams();

  if (application.authorized) {
    return <Navigate to={searchParams.get('origin') || '/'} replace />;
  }

  return <div>Login</div>;
}
