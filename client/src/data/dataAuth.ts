import type { FetchProviderType } from '@civet/common';
import { HTTPError } from './errors';

export type Auth = {
  username: string;
  accessToken: string;
  expires: string;
};

export type PostAuthRequest = {
  username: string;
  password: string;
};

export function postAuth(
  dataProvider: FetchProviderType,
  request: PostAuthRequest,
): Promise<Auth> {
  return dataProvider.request<Auth>(
    'v1/auth',
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    },
    {
      async handleError(_url, _request, response, _meta) {
        throw new HTTPError(await response.text(), response.status);
      },
    },
  );
}
