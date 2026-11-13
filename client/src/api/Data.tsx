import { FetchProvider } from '@civet/common';
import { ConfigProvider } from '@civet/core';
import { ConfigProvider as EventConfigProvider } from '@civet/events';
import type { ReactNode } from 'react';
import { ParanalReceiver } from './ParanalReceiver';

const apiURL = new URL(
  window.paranal.api.replace(/\/*$/, '/'),
  window.location.href,
);

const dataProvider = new FetchProvider({
  baseURL: apiURL,
  handleError: async (_requestURL, _request, response) => {
    let message: string | undefined;
    try {
      message = await response.text();
    } catch {
      // ignore
    }
    throw new Error(message || response.statusText);
  },
});

let sourceController: AbortController | undefined;
const allocateEventSource = () => {
  sourceController?.abort();
  sourceController = new AbortController();
  const eventSource = new EventSource(new URL(`v1/events`, apiURL));
  eventSource.addEventListener(
    'error',
    () => {
      console.warn('connecting to the event stream: unexpected error');
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

function normalizePath(path: string): string {
  return new URL(path, apiURL).pathname
    .substring(apiURL.pathname.length)
    .replace(/\/*$/, '')
    .split('/')
    .map((segment) => encodeURIComponent(decodeURIComponent(segment)))
    .join('/');
}

function matchEvent(pattern: string, path: string): boolean {
  const p = pattern.split('/');
  const e = path.split('/');

  if (p.length !== e.length) return false;

  const differs = p.some((ps, i) => {
    const es = e[i];
    return ps != '*' && ps != es;
  });

  return !differs;
}

const eventReceiver = new ParanalReceiver(allocateEventSource(), {
  events: ['change'],
  getEvents(resource, type, event, options) {
    if (!resource) return [event];

    switch (type) {
      case 'change': {
        const resourcePath = normalizePath(resource.name);
        const eventPath = normalizePath(event.data);

        if (options?.pathPattern) {
          if (matchEvent(options.pathPattern, eventPath)) return [event];

          return [];
        }

        if (eventPath === resourcePath) return [event];

        return [];
      }

      default:
        throw new Error(`unsupported event type "${type}"`);
    }
  },
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
