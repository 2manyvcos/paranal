import { DashboardEngine, DashboardEventTarget } from '@rotorjs/dashboard';
import { createJPLStateReducerConfig } from '@rotorjs/jpl';
import varsPresets from '@rotorjs/jpl/presets/vars';
import { attachWorker } from '@rotorjs/state';
import jplImport from './vars/jplImport';
import jplLocation from './vars/jplLocation';
import jplState from './vars/jplState';
import jplStates from './vars/jplStates';

const target = new DashboardEventTarget();
attachWorker(target, self);

/* const engine = */ new DashboardEngine(target, {
  script: createJPLStateReducerConfig({
    vars: {
      ...varsPresets,
      import: jplImport,
      location: jplLocation,
      states: jplStates,
      state: jplState,
      // TODO: add http function
    },
  }),
});
