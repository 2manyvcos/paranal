import type { FetchProviderType } from '@civet/common';
import { useResource } from '@civet/core';
import type { ThemeOptions } from '@mui/material';
import type {
  DashboardLayoutNode,
  DashboardTileNode,
  DashboardVar,
} from '@rotorjs/dashboard';
import type { ParanalEventOptions } from './ParanalReceiver';

export type Dashboard = {
  version: string;
  id: string;
  routes?: DashboardRoute[];
  theme?: ThemeOptions;
  vars?: Record<string, DashboardVar>;
  defaultLayout?: DashboardLayoutNode;
  layout?: DashboardLayoutNode;
  content?: DashboardTileNode[];
  scripts?: Record<string, string>;
};

export type DashboardRoute = {
  host?: string;
  path: string;
};

export type DashboardSummary = {
  version: string;
  id: string;
  routes?: DashboardRoute[];
};

export type DashboardTheme = Record<string, unknown>;

export function useDashboards() {
  return useResource<FetchProviderType, DashboardSummary[] | undefined>({
    name: 'v1/dashboards',
    query: undefined,
    events: {
      options: {
        pathPattern: 'v1/dashboards/*',
      } satisfies ParanalEventOptions,
    },
  });
}

export function useDashboard(id: string) {
  return useResource<FetchProviderType, Dashboard | undefined>({
    name: `v1/dashboards/${encodeURIComponent(id)}`,
    query: undefined,
    events: true,
  });
}
