import {
  JPLRuntimeError,
  type JPLFunc,
  type JPLRuntime,
  type JPLRuntimeSignal,
} from '@jplorg/jpl';
import {
  dashboardFactInterest,
  type DashboardStateReducer,
} from '@rotorjs/dashboard';

export default function jplTheme(reducer: DashboardStateReducer) {
  return (): JPLFunc => {
    reducer.removeInterest(dashboardFactInterest('theme'));

    return async (
      _runtime: JPLRuntime,
      _signal: JPLRuntimeSignal,
      next: (output: unknown) => Promise<unknown[]> | unknown[],
      _input: unknown,
      ...args: unknown[]
    ): Promise<unknown[]> => {
      if (args.length > 0) throw new JPLRuntimeError('too many arguments');

      reducer.addInterest(dashboardFactInterest('theme'));

      return next(reducer.engine.getFact('theme')?.value);
    };
  };
}
