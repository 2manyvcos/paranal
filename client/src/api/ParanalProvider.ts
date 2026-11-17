import { FetchProvider } from '@civet/common';

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
        throw new Error(message || response.statusText);
      },
    });
  }
}

export type ParanalProviderType = InstanceType<typeof ParanalProvider>;
