import {
  JPLRuntimeError,
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

const programCache: Record<string, JPLProgram> = {};

export default function jplImport(reducer: DashboardStateReducer) {
  const requestedInterests = new Set<string>();

  return (): JPLFunc => {
    requestedInterests.forEach((interest) => {
      reducer.removeInterest(interest);
    });
    requestedInterests.clear();

    return async (
      runtime: JPLRuntime,
      _signal: JPLRuntimeSignal,
      next: (output: unknown) => Promise<unknown[]> | unknown[],
      input: unknown,
      ...args: unknown[]
    ): Promise<unknown[]> => {
      if (args.length > 0) throw new JPLRuntimeError('too many arguments');

      const scriptName = assertType<string>(
        runtime,
        input,
        ['string'],
        'invalid input (%s %*<100v): expected string',
      );
      const scriptFact = `scripts/${encodeURIComponent(scriptName)}`;

      const scriptInterest = dashboardFactInterest(scriptFact);
      requestedInterests.add(scriptInterest);
      reducer.addInterest(scriptInterest);

      const script = reducer.engine.getFact(scriptFact);
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
