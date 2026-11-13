import {
  parse as parseJPLProgram,
  type JPLFunc,
  type JPLProgram,
  type JPLRuntime,
  type JPLRuntimeSignal,
} from '@jplorg/jpl';
import {
  dashboardFactInterest,
  type DashboardStateReducer,
} from '@rotorjs/dashboard';
import { assertType } from './utils';

const programCache: { [src: string]: JPLProgram } = {};

export default function jplImport(reducer: DashboardStateReducer) {
  const requestedScripts = new Set<string>();

  return (): JPLFunc => {
    requestedScripts.forEach((scriptName) => {
      reducer.removeInterest(
        dashboardFactInterest(`scripts/${encodeURIComponent(scriptName)}`),
      );
    });
    requestedScripts.clear();

    return async (
      runtime: JPLRuntime,
      _signal: JPLRuntimeSignal,
      next: (output: unknown) => Promise<unknown[]> | unknown[],
      input: unknown,
    ): Promise<unknown[]> => {
      const scriptName = assertType<string>(
        runtime,
        input,
        ['string'],
        '%s (%*<100v) cannot be used as a script name',
      );
      const factName = `scripts/${encodeURIComponent(scriptName)}`;

      requestedScripts.add(scriptName);
      reducer.addInterest(dashboardFactInterest(factName));

      const script = reducer.engine.getFact(factName);

      if (script == null || typeof script.value !== 'string') {
        throw new Error(`unknown script "${scriptName}"`);
      }

      const src = script.value;
      if (!Object.hasOwn(programCache, src)) {
        programCache[src] = await parseJPLProgram(src);
      }
      const program = programCache[src];
      const vars = runtime.rootScope.vars;

      const i: JPLFunc = async (
        runtime: JPLRuntime,
        signal: JPLRuntimeSignal,
        next: (output: unknown) => Promise<unknown[]> | unknown[],
        input: unknown,
        ...args: unknown[]
      ): Promise<unknown[]> => {
        const controller = new AbortController();
        signal.subscribe(() => {
          controller.abort();
        });

        const results = await program.run([input], {
          runtime: {
            vars: { ...vars, args },
            signal: controller.signal,
          },
        });

        return runtime.muxAll([results], next);
      };

      return next(i);
    };
  };
}
