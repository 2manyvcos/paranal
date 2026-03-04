import { useSyncExternalStore } from 'react';
import Notification from './Notification';
import { getNotifications, subscribeNotifications } from './notifications';

export default function Notifications() {
  const notifications = useSyncExternalStore(
    subscribeNotifications,
    getNotifications,
  );

  return (
    <div className="notifications popup" role="dialog" tabIndex={-1}>
      <div className="panel">
        {notifications.map((notification) => (
          <Notification key={notification.id} notification={notification} />
        ))}
      </div>
    </div>
  );
}
