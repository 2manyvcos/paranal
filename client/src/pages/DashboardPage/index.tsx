import { useDashboard } from '@/api/dataDashboards';
import Dashboard from '@/components/Dashboard';
import LoadingPage from '@/pages/LoadingPage';
import UnexpectedErrorPage from '@/pages/UnexpectedErrorPage';
import { useEffect } from 'react';

export type DashboardPageProps = {
  version: string;
  id: string;
};

export default function DashboardPage({ id }: DashboardPageProps) {
  const dashboard = useDashboard(id);
  const dashboardLoading = dashboard.isLoading && dashboard.isInitial;

  useEffect(() => {
    if (dashboard.error) console.error(dashboard.error);
  }, [dashboard.error]);

  if (dashboard.error) {
    return <UnexpectedErrorPage />;
  }

  if (dashboardLoading) {
    return <LoadingPage />;
  }

  // TODO: dashboard.data?.theme

  return (
    <Dashboard
      vars={dashboard.data?.vars}
      defaultLayout={dashboard.data?.defaultLayout}
      layout={dashboard.data?.layout}
      content={dashboard.data?.content}
      scripts={dashboard.data?.scripts}
    />
  );
}
