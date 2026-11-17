import {
  paranalReconnectInterest,
  paranalStateInterest,
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

type StateConfig = {
  registry?: string;
  id?: string;
};

export default function jplState(reducer: DashboardStateReducer) {
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
      const config = assertJSON<StateConfig>(
        runtime,
        input,
        ['null', 'object'],
        'expected an object as config but got %s (%*<100v)',
      );
      if (!config.registry)
        throw new JPLExecutionError('config.registry is required');
      if (!config.id) throw new JPLExecutionError('config.id is required');

      reducer.addInterest(paranalReconnectInterest());

      const stateInterest = paranalStateInterest(config.registry, config.id);
      requestedInterests.add(stateInterest);
      reducer.addInterest(stateInterest);

      const abortSignal = new AbortSignal();
      signal.subscribe(() => {
        abortSignal.abort();
      });

      const response = await provider.get(
        `v1/registries/${encodeURIComponent(config.registry)}/states/${encodeURIComponent(config.id)}`,
        undefined,
        undefined,
        undefined,
        abortSignal,
      );

      return next(response);
    };
  };
}
