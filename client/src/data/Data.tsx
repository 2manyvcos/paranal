import { FetchProvider, SSEReceiver } from '@civet/common';
import { ConfigProvider } from '@civet/core';
import { ConfigProvider as EventConfigProvider } from '@civet/events';
import { type ReactNode } from 'react';
import { getAccessToken } from './accessTokens';
import { HTTP_UNAUTHORIZED, HTTPError } from './errors';

const apiURL = new URL(
  (import.meta.env.PARANAL_API || '/api/').replace(/\/*$/g, '/'),
  window.location.href,
);

const dataProvider = new FetchProvider({
  baseURL: apiURL,
  modifyRequest(_url, request, _meta): void {
    const headers = (request.headers = new Headers(request.headers));

    const accessToken = getAccessToken();
    if (accessToken) headers.set('Authorization', `Bearer ${accessToken}`);
  },
  async handleError(_url, _request, response, _meta): Promise<never> {
    if (response.status === HTTP_UNAUTHORIZED) {
      dataProvider.notify('v1/user');
    }

    throw new HTTPError(await response.text(), response.status);
  },
});

let sourceController: AbortController | undefined;
const nextEventSource = () => {
  sourceController?.abort();
  sourceController = new AbortController();
  const eventSource = new EventSource(
    new URL(`v1/events?accessToken=${getAccessToken()}`, apiURL),
  );
  eventSource.addEventListener(
    'error',
    () => {
      console.warn('Error connecting to the event stream');
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

const eventReceiver = new SSEReceiver(nextEventSource(), {
  events: ['update'],
  getEvents(resource, type, event) {
    if (!resource) return [event];

    switch (type) {
      case 'update': {
        const resourcePath = new URL(resource.name, apiURL).pathname
          .substring(0, apiURL.toString().length - 1)
          .replace(/\/*$/g, '/');
        const eventPath = event.data.replace(/\/*$/g, '/');

        if (eventPath === resourcePath) return [event];

        const resourceSegmentCount = resourcePath.split('/').length;
        const eventSegmentCount = eventPath.split('/').length;
        // the event may target a record of the current resource, but must not target a sub endpoint
        if (
          eventSegmentCount === resourceSegmentCount + 1 &&
          resourceSegmentCount % 2 === 0
        )
          return [event];

        return [];
      }

      default:
        throw new Error(`unsupported event type: ${type}`);
    }
  },
});

window.addEventListener('online', () => {
  eventReceiver.setEventSource(nextEventSource());
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
