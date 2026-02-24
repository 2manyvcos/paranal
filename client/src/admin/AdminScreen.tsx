import { Outlet } from 'react-router';

export default function AdminScreen() {
  return (
    <div className="admin screen">
      Admin
      <Outlet />
    </div>
  );
}
