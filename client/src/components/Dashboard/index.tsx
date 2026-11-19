import { getAPIURL, matchEvent, normalizePath } from '@/api/api';
import {
  paranalReconnectInterest,
  paranalStateInterest,
  paranalStatesInterest,
} from '@/api/interests';
import type { ParanalProviderType } from '@/api/ParanalProvider';
import type { ParanalReceiverType } from '@/api/ParanalReceiver';
import { useConfigContext } from '@civet/core';
import { useConfigContext as useEventConfigContext } from '@civet/events';
import { CssBaseline } from '@mui/material';
import {
  createTheme,
  ThemeProvider,
  type ThemeOptions,
} from '@mui/material/styles';
import {
  DashboardEventTarget,
  type DashboardLayoutNode,
  type DashboardTileNode,
  type DashboardVar,
} from '@rotorjs/dashboard';
import { jplLayouts, jplTiles } from '@rotorjs/jpl-react';
import {
  Dashboard as RotorDashboard,
  type DashboardLayoutMap,
  type DashboardTileMap,
} from '@rotorjs/react';
import commonLayouts from '@rotorjs/react-tiles/presets/layouts';
import commonTiles from '@rotorjs/react-tiles/presets/tiles';
import { attachWorker } from '@rotorjs/state';
import { useEffect, useMemo, useState } from 'react';
import { useLocation, useParams, useSearchParams } from 'react-router';
import ExposeVars from './ExposeVars';

// eslint-disable-next-line import-x/default
import DashboardWorker from '@/workers/DashboardWorker?worker';

const defaultTheme: ThemeOptions = {
  colorSchemes: {
    dark: true,
  },
  shape: {
    borderRadius: 32,
  },
};

const layouts: DashboardLayoutMap = {
  ...commonLayouts,
  ...jplLayouts,
};

const defaultLayout: DashboardLayoutNode = {
  type: 'grid',
  style: { margin: 2 },
};

const tiles: DashboardTileMap = {
  ...commonTiles,
  ...jplTiles,
};

const defaultContent: DashboardTileNode[] = [];

function approveNavigation(): boolean {
  return true;
}

export type DashboardProps = {
  theme?: ThemeOptions;
  vars?: Record<string, DashboardVar>;
  defaultLayout?: DashboardLayoutNode;
  layout?: DashboardLayoutNode;
  content?: DashboardTileNode[];
  scripts?: Record<string, string>;
};

export default function Dashboard({
  theme: themeProp,
  vars,
  defaultLayout: defaultLayoutProp,
  layout,
  content = defaultContent,
  scripts,
}: DashboardProps) {
  const location = useLocation();
  const pathParams = useParams();
  const [searchParams] = useSearchParams();
  const [target] = useState(() => new DashboardEventTarget());
  const { dataProvider } = useConfigContext<ParanalProviderType>();
  const { eventReceiver } = useEventConfigContext<ParanalReceiverType>();

  const theme = useMemo(() => {
    return createTheme(themeProp ?? defaultTheme);
  }, [themeProp]);

  useEffect(() => {
    const controller = new AbortController();

    const worker = new DashboardWorker();
    attachWorker(target, worker, { signal: controller.signal });

    return () => {
      controller.abort();
      worker.terminate();
    };
  }, [target]);

  const initialVars = useMemo(() => {
    return {
      ...vars,
      ...Object.fromEntries(
        searchParams
          .entries()
          .filter(([name]) => vars?.[name]?.exposed)
          .map(([name, param]) => {
            let value;
            try {
              value = JSON.parse(param || 'true');
            } catch {
              value = param;
            }

            return [name, { value, exposed: true }];
          }),
      ),
    };
  }, [vars, searchParams]);

  const facts = useMemo(() => {
    return {
      apiURL: { value: getAPIURL().toString() },
      location: { value: { pathname: location.pathname, params: pathParams } },
      ...Object.fromEntries(
        Object.entries(scripts ?? {}).map(([name, src]) => [
          `scripts/${encodeURIComponent(name)}`,
          { value: src },
        ]),
      ),
    };
  }, [location, pathParams, scripts]);

  useEffect(() => {
    Object.entries(facts).forEach(([name, { value }]) => {
      target.dispatchFact(name, value);
    });
  }, [target, facts]);

  useEffect(() => {
    const unsub = dataProvider!.subscribe('', () => {
      target.dispatchInterest(paranalReconnectInterest());
    });

    return unsub;
  }, [target, dataProvider]);

  useEffect(() => {
    const unsub = eventReceiver!.subscribe(
      undefined,
      undefined,
      (events: MessageEvent[]) => {
        events.forEach((event) => {
          switch (event.type) {
            case 'change': {
              const eventPath = normalizePath(event.data);

              if (matchEvent('v1/registries/*/states/*', eventPath)) {
                const [, , registry, , id] = eventPath
                  .split('/')
                  .map(decodeURIComponent);
                target.dispatchInterest(paranalStatesInterest(registry));
                target.dispatchInterest(paranalStateInterest(registry, id));
              }
            }
          }
        });
      },
    );

    return unsub;
  }, [target, eventReceiver]);

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline enableColorScheme />

      <RotorDashboard
        target={target}
        initialVars={initialVars}
        initialFacts={facts}
        layouts={layouts}
        defaultLayout={defaultLayoutProp ?? defaultLayout}
        tiles={tiles}
        layout={layout}
        content={content}
        // approveUserAction={}
        approveNavigation={approveNavigation}
      >
        <ExposeVars />
      </RotorDashboard>
    </ThemeProvider>
  );
}
