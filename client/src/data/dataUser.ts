import type { FetchProviderType } from '@civet/common';
import { Meta, useResource } from '@civet/core';
import { HTTP_UNAUTHORIZED, HTTPError } from '@/data/errors';

export type UserQuery = {
  disabled?: boolean;
};

export type User = {
  name: string;
  role: string;
  hasPassword: boolean;
  displayName: string;
  startPage: string;
  errorAlerts: boolean;
  healthAlerts: boolean;
  versionAlerts: boolean;
};

async function handleError(
  _url: URL,
  _request: RequestInit,
  response: Response,
  _meta: Meta,
): Promise<undefined> {
  if (response.status === HTTP_UNAUTHORIZED) return undefined;

  throw new HTTPError(await response.text(), response.status);
}

export function useUser(query: UserQuery) {
  return useResource<FetchProviderType, User | undefined>({
    name: 'v1/user',
    query: undefined,
    options: { handleError },
    events: true,
    disabled: query.disabled,
  });
}
