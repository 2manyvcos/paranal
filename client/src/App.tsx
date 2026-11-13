import { useEffect } from 'react';
import { Route, Routes } from 'react-router';
import { useDashboards } from './api/dataDashboards';
import DashboardPage from './pages/DashboardPage';
import LoadingPage from './pages/LoadingPage';
import NotFoundPage from './pages/NotFoundPage';
import UnexpectedErrorPage from './pages/UnexpectedErrorPage';

export default function App() {
  const host = window.location.host;

  const dashboards = useDashboards();
  const dashboardsLoading = dashboards.isLoading && dashboards.isInitial;

  useEffect(() => {
    if (dashboards.error) console.error(dashboards.error);
  }, [dashboards.error]);

  if (dashboards.error) {
    return <UnexpectedErrorPage />;
  }

  if (dashboardsLoading) {
    return <LoadingPage />;
  }

  return (
    <Routes>
      {dashboards.data?.map((dashboard) => {
        return dashboard.routes?.map((route) => {
          if (route.host && route.host !== host) return undefined;

          return (
            <Route
              path={route.path}
              element={
                <DashboardPage version={dashboard.version} id={dashboard.id} />
              }
            />
          );
        });
      })}

      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
