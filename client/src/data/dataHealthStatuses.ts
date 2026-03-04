import type { FetchProviderType } from '@civet/common';
import { useResource } from '@civet/core';

export type HealthStatusQuery = {
  disabled?: boolean;
  name?: string;
  serviceID?: string;
  status?: string;
  unhealthy?: boolean;
  service?: {
    config?: {
      hidden?: boolean;
    };
  };
};

export type HealthStatus = {
  name: string;
  serviceID: string;
  status: string;
  unhealthy: boolean;
};

function healthStatusQuery(query: HealthStatusQuery): string {
  const search = new URLSearchParams();
  if (query.name != null) search.set('name', query.name);
  if (query.serviceID != null) search.set('serviceID', query.serviceID);
  if (query.status != null) search.set('status', query.status);
  if (query.unhealthy != null)
    search.set('unhealthy', query.unhealthy.toString());
  if (query.service?.config?.hidden != null)
    search.set('service.config.hidden', query.service.config.hidden.toString());
  return search.toString();
}

export function useHealthStatuses(query: HealthStatusQuery) {
  return useResource<FetchProviderType, HealthStatus[] | undefined>({
    name: 'v1/health-statuses',
    query: { search: healthStatusQuery(query) },
    events: true,
    disabled: query.disabled,
  });
}
