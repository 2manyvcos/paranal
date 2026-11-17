import {
  paranalReconnectInterest,
  paranalStatesInterest,
} from '@/api/interests';
import ParanalProvider from '@/api/ParanalProvider';
import { AbortSignal } from '@civet/core';
import {
  JPLExecutionError,
  type JPLFunc,
  type JPLRuntime,
  type JPLRuntimeSignal,
} from '@jplorg/jpl';
import { type DashboardStateReducer } from '@rotorjs/dashboard';
import { assertJSON } from './utils';

type StatesConfig = {
  registry?: string;
};

export default function jplStates(reducer: DashboardStateReducer) {
  const requestedInterests = new Set<string>();

  const apiURL = reducer.engine.getFact('apiURL');
  if (apiURL == null || typeof apiURL.value !== 'string') {
    throw new Error('missing API URL');
  }
  const provider = new ParanalProvider(apiURL.value);

  return (): JPLFunc => {
    reducer.removeInterest(paranalReconnectInterest());

    requestedInterests.forEach((interest) => {
      reducer.removeInterest(interest);
    });
    requestedInterests.clear();

    return async (
      runtime: JPLRuntime,
      signal: JPLRuntimeSignal,
      next: (output: unknown) => Promise<unknown[]> | unknown[],
      input: unknown,
    ): Promise<unknown[]> => {
      const config = assertJSON<StatesConfig>(
        runtime,
        input,
        ['null', 'object'],
        'expected an object as config but got %s (%*<100v)',
      );
      if (!config.registry)
        throw new JPLExecutionError('config.registry is required');

      reducer.addInterest(paranalReconnectInterest());

      const statesInterest = paranalStatesInterest(config.registry);
      requestedInterests.add(statesInterest);
      reducer.addInterest(statesInterest);

      const abortSignal = new AbortSignal();
      signal.subscribe(() => {
        abortSignal.abort();
      });

      const response = await provider.get(
        `v1/registries/${encodeURIComponent(config.registry)}/states`,
        undefined,
        undefined,
        undefined,
        abortSignal,
      );

      return next(response);
    };
  };
}
