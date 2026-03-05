import type { FetchProviderType } from '@civet/common';
import { HTTPError } from './errors';

export type Auth = {
  userName: string;
  accessToken: string;
  expires: string;
};

export function postAuth(request: {
  dataProvider: FetchProviderType;
  data: {
    userName: string;
    password: string;
  };
}): Promise<Auth> {
  return request.dataProvider.request<Auth>(
    'v1/auth',
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request.data),
    },
    {
      async handleError(_url, _request, response, _meta) {
        throw new HTTPError(await response.text(), response.status);
      },
    },
  );
}
