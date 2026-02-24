import { useEffect } from 'react';
import { Outlet, useSearchParams } from 'react-router';
import { useApplication } from '@/application';

export default function AuthScreen() {
  const { authorized } = useApplication();
  const [searchParams] = useSearchParams();

  useEffect(() => {
    if (authorized) window.location.replace(searchParams.get('origin') || '/');
  }, [authorized, searchParams]);

  if (authorized) {
    return null;
  }

  return (
    <div className="auth screen">
      <Outlet />
    </div>
  );
}
