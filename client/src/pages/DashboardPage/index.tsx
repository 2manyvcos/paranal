import { useDashboard } from '@/api/dataDashboards';
import Dashboard from '@/components/Dashboard';
import ErrorBoundary from '@/components/ErrorBoundary';
import LoadingPage from '@/pages/LoadingPage';
import UnexpectedErrorPage from '@/pages/UnexpectedErrorPage';
import objectHash from 'object-hash';
import { useMemo } from 'react';
import { useLocation } from 'react-router';

export type DashboardPageProps = {
  version: string;
  id: string;
};

export default function DashboardPage({ id }: DashboardPageProps) {
  const location = useLocation();

  const dashboard = useDashboard(id);
  const dashboardLoading = dashboard.isLoading && dashboard.isInitial;

  const dashboardKey = useMemo(
    () =>
      objectHash({
        pathname: location.pathname,
        id: dashboard.data?.id ?? '',
        vars: dashboard.data?.vars ?? {},
      }),
    [location.pathname, dashboard.data],
  );

  if (dashboard.error) {
    return <UnexpectedErrorPage error={dashboard.error} />;
  }

  if (dashboardLoading) {
    return <LoadingPage />;
  }

  return (
    <ErrorBoundary
      fallback={(error: unknown) => <UnexpectedErrorPage error={error} />}
      errorKey={dashboard.data}
    >
      <Dashboard
        key={dashboardKey}
        theme={dashboard.data?.theme}
        vars={dashboard.data?.vars}
        defaultLayout={dashboard.data?.defaultLayout}
        layout={dashboard.data?.layout}
        content={dashboard.data?.content}
        scripts={dashboard.data?.scripts}
      />
    </ErrorBoundary>
  );
}
