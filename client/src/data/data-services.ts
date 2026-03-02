import type { FetchProviderType } from '@civet/common';
import { useResource } from '@civet/core';

export type ServiceQuery = {
  disabled?: boolean;
  id?: string;
  name?: string;
  description?: string;
  logo?: string;
  url?: string;
  config?: {
    favorite?: boolean;
    hidden?: boolean;
    uptimeAlerts?: boolean;
    versionAlerts?: boolean;
  };
};

export type Service = {
  id: string;
  name: string;
  description: string;
  logo: string;
  url: string;
  config: {
    favorite: boolean;
    hidden: boolean;
    uptimeAlerts: boolean;
    versionAlerts: boolean;
  };
};

function serviceQuery(query?: ServiceQuery): string {
  const search = new URLSearchParams();
  if (query?.id != null) search.set('id', query.id);
  if (query?.name != null) search.set('name', query.name);
  if (query?.description != null) search.set('description', query.description);
  if (query?.logo != null) search.set('logo', query.logo);
  if (query?.url != null) search.set('url', query.url);
  if (query?.config?.favorite != null)
    search.set('config.favorite', query.config.favorite.toString());
  if (query?.config?.hidden != null)
    search.set('config.hidden', query.config.hidden.toString());
  if (query?.config?.uptimeAlerts != null)
    search.set('config.uptimeAlerts', query.config.uptimeAlerts.toString());
  if (query?.config?.versionAlerts != null)
    search.set('config.versionAlerts', query.config.versionAlerts.toString());
  return search.toString();
}

export function useServices(query?: ServiceQuery) {
  return useResource<FetchProviderType, Service[] | undefined>({
    name: 'v1/services',
    query: { search: serviceQuery(query) },
    events: true,
    disabled: query?.disabled,
  });
}
