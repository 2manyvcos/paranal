import type { FetchProviderType } from '@civet/common';
import { useResource } from '@civet/core';

export type ActionQuery = {
  disabled?: boolean;
  serviceID: string;
};

export type Actions = {
  groups: ActionGroup[];
  actions: Action[];
};

export type ActionGroup = {
  name: string;
};

export type Action = {
  name: string;
  url: string;
  canRun: boolean;
  group: string;
  restrictToAdmins: boolean;
};

export function useServicesByIDActions(query: ActionQuery) {
  return useResource<FetchProviderType, Actions | undefined>({
    name: `v1/services/${encodeURIComponent(query.serviceID)}/actions`,
    query: undefined,
    events: true,
    disabled: query.disabled,
  });
}

export type ActionResult = {
  success: boolean;
  error?: string;
  results?: unknown[];
};

export function postServicesByIDActionsByNameRun(request: {
  dataProvider: FetchProviderType;
  serviceID: string;
  actionName: string;
}) {
  return request.dataProvider.request<ActionResult>(
    `v1/services/${encodeURIComponent(request.serviceID)}/actions/${encodeURIComponent(request.actionName)}/run`,
    { method: 'POST' },
  );
}
