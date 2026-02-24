import { useEffect } from 'react';
import { useLocation } from 'react-router';

export default function NavigateToLoginPage() {
  const location = useLocation();

  useEffect(() => {
    const search = new URLSearchParams();
    search.set('origin', location.pathname + location.search + location.hash);
    const searchString = search.toString();

    window.location.replace('/login?' + searchString);
  }, [location]);

  return null;
}
