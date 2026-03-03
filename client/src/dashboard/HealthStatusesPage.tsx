import { useApplication } from '@/application';

export default function HealthStatusesPage() {
  const { appName } = useApplication();

  return (
    <div className="health-statuses page">
      <title>{`Health | ${appName}`}</title>
      TODO: health page
    </div>
  );
}
