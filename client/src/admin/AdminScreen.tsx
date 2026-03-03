import { Outlet } from 'react-router';
import { useApplication } from '@/application';

export default function AdminScreen() {
  const { appName } = useApplication();

  return (
    <div className="admin screen">
      <title>{`Admin Panel | ${appName}`}</title>
      TODO: Admin
      <Outlet />
    </div>
  );
}
