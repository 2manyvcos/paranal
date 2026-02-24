import { useEffect } from 'react';
import { useLocation } from 'react-router';

export default function NavigateToLoginPage() {
  const location = useLocation();

  useEffect(() => {
    const currentLocation = location.pathname + location.search + location.hash;
    const url = new URL('/login', window.location.href);
    url.searchParams.set('origin', currentLocation);
    window.location.replace(url);
  }, [location]);

  return null;
}
