import { useApplication } from '@/application';

export default function VersionsPage() {
  const { appName } = useApplication();

  return (
    <div className="versions page">
      <title>{`Versions | ${appName}`}</title>
      TODO: versions page
    </div>
  );
}
