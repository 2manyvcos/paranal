import { ConfigProvider } from '@civet/core';
import { ConfigProvider as EventConfigProvider } from '@civet/events';
import type { ReactNode } from 'react';
import { getAPIURL } from './api';
import ParanalProvider from './ParanalProvider';
import ParanalReceiver from './ParanalReceiver';

let sourceController: AbortController | undefined;
const allocateEventSource = () => {
  sourceController?.abort();
  sourceController = new AbortController();
  const eventSource = new EventSource(new URL(`v1/events`, getAPIURL()));
  eventSource.addEventListener(
    'error',
    () => {
      console.warn('connecting to the event stream: unexpected error');

      // TODO: auto reconnect after timeout
    },
    { signal: sourceController.signal },
  );
  document.addEventListener(
    'visibilitychange',
    () => {
      if (document.visibilityState !== 'hidden') return;
      eventSource.close();
    },
    { signal: sourceController.signal },
  );
  window.addEventListener(
    'beforeunload',
    () => {
      eventSource.close();
    },
    { signal: sourceController.signal },
  );
  return eventSource;
};

const dataProvider = new ParanalProvider(getAPIURL());
const eventReceiver = new ParanalReceiver(allocateEventSource());

document.addEventListener('visibilitychange', () => {
  if (document.visibilityState === 'hidden') return;
  eventReceiver.setEventSource(allocateEventSource());
  dataProvider.notify(undefined);
});

window.addEventListener('online', () => {
  eventReceiver.setEventSource(allocateEventSource());
  dataProvider.notify(undefined);
});

export default function Data({ children }: { children: ReactNode }) {
  return (
    <ConfigProvider dataProvider={dataProvider}>
      <EventConfigProvider eventReceiver={eventReceiver}>
        {children}
      </EventConfigProvider>
    </ConfigProvider>
  );
}
