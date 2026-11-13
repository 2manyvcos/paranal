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
import { useSearchParams } from 'react-router';
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
  const [searchParams] = useSearchParams();
  const [target] = useState(() => new DashboardEventTarget());

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

  const initialFacts = useMemo(() => {
    return Object.fromEntries(
      Object.entries(scripts ?? {}).map(([name, src]) => [
        `scripts/${encodeURIComponent(name)}`,
        { value: src },
      ]),
    );
  }, [scripts]);

  // TODO: either update when facts / vars change in the definition, or don't subscribe to events at all

  return (
    <RotorDashboard
      target={target}
      initialVars={initialVars}
      initialFacts={initialFacts}
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
