import { Outlet } from 'react-router';

export default function ErrorScreen() {
  return (
    <div className="error screen">
      <Outlet />
    </div>
  );
}
