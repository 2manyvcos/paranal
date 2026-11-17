import { DashboardEngine, DashboardEventTarget } from '@rotorjs/dashboard';
import { createJPLStateReducerConfig } from '@rotorjs/jpl';
import varsPresets from '@rotorjs/jpl/presets/vars';
import { attachWorker } from '@rotorjs/state';
import jplImport from './jpl/jplImport';
import jplLocation from './jpl/jplLocation';

const target = new DashboardEventTarget();
attachWorker(target, self);

/* const engine = */ new DashboardEngine(target, {
  script: createJPLStateReducerConfig({
    vars: {
      ...varsPresets,
      import: jplImport,
      location: jplLocation,
      // TODO: add http function
      // TODO: add state resolver function; use documentVisibilityState and connectivityState facts for event handler (or send change events via DashboardEventTarget)
    },
  }),
});
