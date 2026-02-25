import { Button, Transition } from '@headlessui/react';
import { useEffect, useState, useSyncExternalStore } from 'react';
import {
  getCurrentNotification,
  removeCurrentNotification,
  subscribeNotifications,
} from './notification';

export default function Notifications() {
  const [visible, setVisible] = useState(false);
  const [autoHide, setAutoHide] = useState<number>();
  const currentNotification = useSyncExternalStore(
    subscribeNotifications,
    getCurrentNotification,
  );

  useEffect(() => {
    if (currentNotification == null) return;
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setVisible(true);
  }, [currentNotification]);

  if (currentNotification == null) return null;

  return (
    <Transition
      key={currentNotification.id}
      show={visible}
      afterEnter={() => {
        setAutoHide(
          setTimeout(() => {
            setVisible(false);
          }, 3000),
        );
      }}
      beforeLeave={() => {
        clearTimeout(autoHide);
      }}
      afterLeave={() => {
        removeCurrentNotification();
      }}
    >
      <div className="notification popup" role="dialog" tabIndex={-1}>
        <div className="panel">
          <div className="message">{currentNotification.message}</div>

          <Button
            className="close"
            onClick={() => {
              setVisible(false);
            }}
          >
            <span className="text">Close</span>
          </Button>
        </div>
      </div>
    </Transition>
  );
}
