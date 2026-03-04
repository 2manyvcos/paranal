import type { FetchProviderType } from '@civet/common';
import { useResource } from '@civet/core';

export type VersionQuery = {
  disabled?: boolean;
  name?: string;
  serviceID?: string;
  currentVersion?: string;
  hasCurrentCVEs?: boolean;
  latestVersion?: string;
  hasLatestCVEs?: boolean;
  status?: string;
  outdated?: boolean;
  vulnerable?: boolean;
  service?: {
    config?: {
      hidden?: boolean;
    };
  };
};

export type Version = {
  name: string;
  serviceID: string;
  currentVersion: string;
  currentCVEs: number;
  latestVersion: string;
  latestCVEs: number;
  status: string;
  outdated: boolean;
  vulnerable: boolean;
};

function versionQuery(query: VersionQuery): string {
  const search = new URLSearchParams();
  if (query.name != null) search.set('name', query.name);
  if (query.serviceID != null) search.set('serviceID', query.serviceID);
  if (query.currentVersion != null)
    search.set('currentVersion', query.currentVersion);
  if (query.hasCurrentCVEs != null)
    search.set('hasCurrentCVEs', query.hasCurrentCVEs.toString());
  if (query.latestVersion != null)
    search.set('latestVersion', query.latestVersion);
  if (query.hasLatestCVEs != null)
    search.set('hasLatestCVEs', query.hasLatestCVEs.toString());
  if (query.status != null) search.set('status', query.status);
  if (query.outdated != null) search.set('outdated', query.outdated.toString());
  if (query.vulnerable != null)
    search.set('vulnerable', query.vulnerable.toString());
  if (query.service?.config?.hidden != null)
    search.set('service.config.hidden', query.service.config.hidden.toString());
  return search.toString();
}

export function useVersions(query: VersionQuery) {
  return useResource<FetchProviderType, Version[] | undefined>({
    name: 'v1/versions',
    query: { search: versionQuery(query) },
    events: true,
    disabled: query.disabled,
  });
}
