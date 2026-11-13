import { DashboardEngine, DashboardEventTarget } from '@rotorjs/dashboard';
import { createJPLStateReducerConfig } from '@rotorjs/jpl';
import varsPresets from '@rotorjs/jpl/presets/vars';
import { attachWorker } from '@rotorjs/state';
import jplImport from './jpl/jplImport';

const target = new DashboardEventTarget();
attachWorker(target, self);
// TODO: cleanup?
const _engine = new DashboardEngine(target, {
  script: createJPLStateReducerConfig({
    vars: {
      ...varsPresets,
      import: jplImport,
    },
  }),
});
