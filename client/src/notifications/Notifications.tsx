import { useSyncExternalStore } from 'react';
import { createPortal } from 'react-dom';
import Notification from './Notification';
import { getNotifications, subscribeNotifications } from './notifications';

export default function Notifications() {
  const notifications = useSyncExternalStore(
    subscribeNotifications,
    getNotifications,
  );

  return createPortal(
    <div className="notifications" role="dialog" tabIndex={-1}>
      <div className="panel">
        {notifications.map((notification) => (
          <Notification key={notification.id} notification={notification} />
        ))}
      </div>
    </div>,
    document.body,
  );
}
