import { Notifier } from '@civet/core';
import { v4 as uuid } from 'uuid';

export type NotificationMessage = {
  message: string;
  loading?: boolean;
  success?: boolean;
  error?: boolean;
};

export type Notification = {
  id: string;
  controller: AbortController;
  notifier: Notifier<[]>;
  data: NotificationMessage;
  remove: () => void;
};

let notifications: Notification[] = [];
const notificationsNotifier = new Notifier<[]>();

export function subscribeNotifications(cb: () => void): () => void {
  return notificationsNotifier.subscribe(cb);
}

export function getNotifications(): Notification[] {
  return notifications;
}

function registerNotification(
  data: NotificationMessage,
  signal?: AbortSignal,
  notifier?: Notifier<[NotificationMessage]>,
) {
  if (signal?.aborted) return;

  const notification: Notification = {
    id: uuid(),
    controller: new AbortController(),
    notifier: new Notifier(),
    data,
    remove: () => {
      notifications = notifications.filter(
        (item) => item.id !== notification.id,
      );
      notificationsNotifier.trigger();
    },
  };
  signal?.addEventListener('abort', () => {
    notification.controller.abort();
  });
  notifier?.subscribe((data) => {
    notification.data = data;
    notification.notifier.trigger();
  });

  notifications = [...notifications, notification];
  notificationsNotifier.trigger();
}

export function notify(message: string, signal?: AbortSignal) {
  registerNotification({ message }, signal);
}

notify.loading = (message: string, signal: AbortSignal) => {
  registerNotification({ message, loading: true }, signal);
};

notify.success = (message: string, signal?: AbortSignal) => {
  registerNotification({ message, success: true }, signal);
};

notify.error = (message: string, signal?: AbortSignal) => {
  registerNotification({ message, error: true }, signal);
};

notify.promise = (
  promise: Promise<unknown>,
  messages: { loading: string; success?: string; error?: string },
  signal?: AbortSignal,
) => {
  const notifier = new Notifier<[NotificationMessage]>();
  registerNotification(
    { message: messages.loading, loading: true },
    signal,
    notifier,
  );
  promise
    .then((result) => {
      notifier.trigger({
        message: messages.success ?? result!.toString(),
        success: true,
      });
    })
    .catch((error) => {
      notifier.trigger({
        message: messages.error ?? error.toString(),
        error: true,
      });
    });
};
