import { Transition } from '@headlessui/react';
import { useCallback, useEffect, useState, useSyncExternalStore } from 'react';
import Button from '@/components/Button';
import Text from '@/components/Text';
import { type Notification as NotificationType } from './notifications';

export default function Notification({
  notification,
}: {
  notification: NotificationType;
}) {
  const [shouldShow, setShouldShow] = useState(false);
  const [visible, setVisible] = useState(false);
  const data = useSyncExternalStore(
    notification.notifier.subscribe,
    useCallback(() => notification.data, [notification]),
  );

  useEffect(() => {
    if (notification.controller.signal.aborted) {
      notification.remove();
      return;
    }
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setShouldShow(true);
    const controller = new AbortController();
    notification.controller.signal.addEventListener(
      'abort',
      () => {
        setShouldShow(false);
      },
      { signal: controller.signal },
    );
    return () => {
      controller.abort();
    };
  }, [notification]);

  useEffect(() => {
    if (!visible || data.loading) return;
    const timeout = setTimeout(
      () => {
        setShouldShow(false);
      },
      data.error ? 5000 : 3000,
    );
    return () => {
      clearTimeout(timeout);
    };
  }, [visible, data]);

  return (
    <Transition
      show={shouldShow}
      afterEnter={() => {
        setVisible(true);
      }}
      beforeLeave={() => {
        setVisible(false);
      }}
      afterLeave={() => {
        notification.remove();
      }}
    >
      <div
        className="notification"
        data-loading={data.loading ? '' : undefined}
        data-success={data.success ? '' : undefined}
        data-error={data.error ? '' : undefined}
      >
        <Text className="message" text={data.message} />

        <Button
          className="close"
          onClick={() => {
            notification.controller.abort();
          }}
          text="Close"
        />
      </div>
    </Transition>
  );
}
