import { FetchProvider } from '@civet/common';
import { HTTPError } from './errors';

export default class ParanalProvider extends FetchProvider {
  constructor(apiURL: string | URL) {
    super({
      baseURL: apiURL,
      handleError: async (_requestURL, _request, response) => {
        let message: string | undefined;
        try {
          message = await response.text();
        } catch {
          // ignore
        }
        throw new HTTPError(message || response.statusText, response.status);
      },
    });
  }
}

export type ParanalProviderType = InstanceType<typeof ParanalProvider>;
