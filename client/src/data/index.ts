import { FetchProvider, SSEReceiver } from '@civet/common';

export class HTTPError extends Error {
  readonly status: number;

  constructor(msg: string, status: number) {
    super(msg);
    Object.setPrototypeOf(this, HTTPError.prototype);
    this.status = status;
  }
}

export const HTTP_UNAUTHORIZED = 401;

const accessTokenKey = 'paranal-access-token';
const apiURL = new URL(
  (import.meta.env.PARANAL_API || '/api/').replace(/\/*$/g, '/'),
  window.location.href,
);

export const dataProvider = new FetchProvider({
  baseURL: apiURL,
  modifyRequest(_url, request, _meta): void {
    const headers = (request.headers = new Headers(request.headers));

    const accessToken = localStorage.getItem(accessTokenKey);
    if (accessToken) headers.set('Authorization', `Bearer ${accessToken}`);
  },
  async handleError(_url, _request, response, _meta): Promise<never> {
    if (response.status === HTTP_UNAUTHORIZED) {
      dataProvider.notify('v1/user');
    }

    throw new HTTPError(await response.text(), response.status);
  },
});

const eventSource = new EventSource(
  new URL(
    `v1/events?accessToken=${localStorage.getItem(accessTokenKey) ?? ''}`,
    apiURL,
  ),
);
eventSource.addEventListener('error', () => {
  console.warn('SSE client failed to connect');
});
window.addEventListener('beforeunload', () => {
  eventSource.close();
});

export const eventReceiver = new SSEReceiver(eventSource, {
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
