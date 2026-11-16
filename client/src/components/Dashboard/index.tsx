// eslint-disable-next-line import-x/default
import DashboardWorker from '@/workers/DashboardWorker?worker';
import {
  DashboardEventTarget,
  type DashboardLayoutNode,
  type DashboardTileNode,
  type DashboardVar,
  type ErrorDashboardNode,
} from '@rotorjs/dashboard';
import { jplLayouts, jplTiles } from '@rotorjs/jpl-react';
import {
  Dashboard as RotorDashboard,
  type DashboardLayoutMap,
  type DashboardTileMap,
} from '@rotorjs/react';
import { attachWorker } from '@rotorjs/state';
import {
  useEffect,
  useMemo,
  useState,
  type ComponentType,
  type PropsWithChildren,
} from 'react';
import { useLocation, useParams, useSearchParams } from 'react-router';
import ExposeVars from './ExposeVars';

function CustomError({ error }: ErrorDashboardNode) {
  // TODO:

  return (
    <div
      style={{
        minWidth: 100,
        minHeight: 100,
        color: 'red',
        textAlign: 'start',
        border: '1px dashed red',
        boxSizing: 'border-box',
        borderRadius: 10,
        padding: 10,
      }}
    >
      Error:
      <pre>{(error as Error)?.message ?? error?.toString()}</pre>
    </div>
  );
}

const layouts: DashboardLayoutMap = {
  ...jplLayouts,
  error: CustomError as ComponentType<PropsWithChildren<DashboardLayoutNode>>,
};

const defaultLayout: DashboardLayoutNode = { type: 'grid' };

const tiles: DashboardTileMap = {
  ...jplTiles,
  error: CustomError as ComponentType<DashboardTileNode>,
};

const defaultContent: DashboardTileNode[] = [];

function approveNavigation(): boolean {
  return true;
}

export type DashboardProps = {
  vars?: Record<string, DashboardVar>;
  defaultLayout?: DashboardLayoutNode;
  layout?: DashboardLayoutNode;
  content?: DashboardTileNode[];
  scripts?: Record<string, string>;
};

export default function Dashboard({
  vars,
  defaultLayout: defaultLayoutProp = defaultLayout,
  layout,
  content = defaultContent,
  scripts,
}: DashboardProps) {
  const location = useLocation();
  const pathParams = useParams();
  const [searchParams] = useSearchParams();
  const [target] = useState(() => new DashboardEventTarget());

  useEffect(() => {
    const controller = new AbortController();

    const worker = new DashboardWorker();
    attachWorker(target, worker, { signal: controller.signal });

    // TODO: use in worker event handler

    document.addEventListener(
      'visibilitychange',
      () => {
        target.dispatchFact(
          'documentVisibilityState',
          document.visibilityState,
        );
      },
      { signal: controller.signal },
    );

    window.addEventListener(
      'online',
      () => {
        target.dispatchFact('connectivityState', 'online');
      },
      { signal: controller.signal },
    );

    window.addEventListener(
      'offline',
      () => {
        target.dispatchFact('connectivityState', 'offline');
      },
      { signal: controller.signal },
    );

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
      connectivityState: {
        value: window.navigator.onLine ? 'online' : 'offline',
      },
      documentVisibilityState: { value: document.visibilityState },
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

  return (
    <RotorDashboard
      target={target}
      initialVars={initialVars}
      initialFacts={facts}
      layouts={layouts}
      defaultLayout={defaultLayoutProp}
      tiles={tiles}
      layout={layout}
      content={content}
      // approveUserAction={}
      approveNavigation={approveNavigation}
    >
      <ExposeVars />
    </RotorDashboard>
  );
}
