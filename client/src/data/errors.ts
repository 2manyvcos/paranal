import type { FetchProviderType } from '@civet/common';
import type { ResourceContextValue } from '@civet/core';
import { useEffect } from 'react';
import { notify } from '@/notifications/notification';

export class HTTPError extends Error {
  readonly status: number;

  constructor(msg: string, status: number) {
    super(msg);
    Object.setPrototypeOf(this, HTTPError.prototype);
    this.status = status;
  }
}

export const HTTP_BAD_REQUEST = 400;
export const HTTP_UNAUTHORIZED = 401;
export const HTTP_FORBIDDEN = 403;
export const HTTP_METHOD_NOT_ALLOWED = 405;
export const HTTP_CONFLICT = 409;
export const HTTP_UNSUPPORTED_MEDIA_TYPE = 415;
export const HTTP_NOT_FOUND = 404;
export const HTTP_INTERNAL_SERVER_ERROR = 500;

export function useErrorNotification(
  resource: ResourceContextValue<FetchProviderType>,
): void {
  useEffect(() => {
    if (resource.error) {
      if (resource.error instanceof HTTPError) {
        switch (resource.error.status) {
          case HTTP_UNAUTHORIZED:
          case HTTP_NOT_FOUND:
            return;
          case HTTP_BAD_REQUEST:
          case HTTP_FORBIDDEN:
          case HTTP_CONFLICT:
            console.error(resource.error);
            notify(resource.error.message);
            return;
        }
      }

      console.error(resource.error);
      notify('Something went wrong.');
    }
  }, [resource.error]);
}
