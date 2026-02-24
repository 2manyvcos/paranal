import { Notifier } from '@civet/core';
import { v4 as uuid } from 'uuid';

export type Notification = {
  id?: string;
  message: string;
};

let previousNotification: Notification | undefined;
const notifications: Notification[] = [];
const notifier = new Notifier<[]>();

export function subscribeNotifications(cb: () => void): () => void {
  return notifier.subscribe(cb);
}

export function getCurrentNotification(): Notification {
  return notifications[0];
}

export function getPreviousNotification(): Notification | undefined {
  return previousNotification;
}

export function removeCurrentNotification() {
  previousNotification = notifications.shift();
  notifier.trigger();
}

export function notify(message: string) {
  notifications.push({ id: uuid(), message });
  notifier.trigger();
}
