import { useEffect } from 'react';
import { useLocation } from 'react-router';

export default function NavigateToLoginPage() {
  const location = useLocation();

  const search = new URLSearchParams();
  search.set('origin', location.pathname + location.search + location.hash);
  const searchString = search.toString();

  useEffect(() => {
    window.location.replace('/login?' + searchString);
  }, [searchString]);

  return null;
}
