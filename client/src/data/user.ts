import type { FetchProviderType } from '@civet/common';
import { Meta, useResource } from '@civet/core';
import { useCallback } from 'react';
import { HTTPError } from '@/data';

export type User = {
  name: string;
  role: string;
  hasPassword: boolean;
  displayName: string;
  startPage: string;
  errorAlerts: boolean;
  uptimeAlerts: boolean;
  versionAlerts: boolean;
};

export function useUser() {
  return useResource<FetchProviderType, User | undefined>({
    name: 'v1/user',
    query: undefined,
    options: {
      handleError: useCallback(
        async (
          _url: URL,
          _request: RequestInit,
          response: Response,
          _meta: Meta,
        ): Promise<never> => {
          throw new HTTPError(await response.text(), response.status);
        },
        [],
      ),
    },
    events: true,
  });
}
