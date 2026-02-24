import type { FetchProviderType } from '@civet/common';
import { Meta, useResource } from '@civet/core';
import { HTTPError } from '@/data/errors';

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

async function handleError(
  _url: URL,
  _request: RequestInit,
  response: Response,
  _meta: Meta,
): Promise<never> {
  throw new HTTPError(await response.text(), response.status);
}

export function useUser() {
  return useResource<FetchProviderType, User | undefined>({
    name: 'v1/user',
    query: undefined,
    options: { handleError },
    events: true,
  });
}
