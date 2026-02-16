import { Navigate, useLocation } from 'react-router';

export default function NavigateToLoginPage() {
  const location = useLocation();

  const search = new URLSearchParams();
  search.set('origin', location.pathname + location.search + location.hash);

  return <Navigate to={{ pathname: '/login', search: search.toString() }} />;
}
