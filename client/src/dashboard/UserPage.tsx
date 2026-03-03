import { useApplication } from '@/application';

export default function UserPage() {
  const { appName } = useApplication();

  return (
    <div className="user page">
      <title>{`Settings | ${appName}`}</title>
      TODO: user page
    </div>
  );
}
