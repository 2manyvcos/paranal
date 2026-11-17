import type { JPLFunc, JPLRuntime, JPLRuntimeSignal } from '@jplorg/jpl';
import {
  dashboardFactInterest,
  type DashboardStateReducer,
} from '@rotorjs/dashboard';

export default function jplLocation(reducer: DashboardStateReducer) {
  return (): JPLFunc => {
    reducer.removeInterest(dashboardFactInterest('location'));

    return async (
      _runtime: JPLRuntime,
      _signal: JPLRuntimeSignal,
      next: (output: unknown) => Promise<unknown[]> | unknown[],
    ): Promise<unknown[]> => {
      reducer.addInterest(dashboardFactInterest('location'));

      return next(reducer.engine.getFact('location')?.value);
    };
  };
}
