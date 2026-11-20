import {
  paranalReconnectInterest,
  paranalStatesInterest,
} from '@/api/interests';
import ParanalProvider from '@/api/ParanalProvider';
import { AbortSignal } from '@civet/core';
import {
  JPLRuntimeError,
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
      _input: unknown,
      ...args: unknown[]
    ): Promise<unknown[]> => {
      if (args.length < 1) throw new JPLRuntimeError('not enough arguments');
      if (args.length > 1) throw new JPLRuntimeError('too many arguments');

      const config = assertJSON<StatesConfig>(
        runtime,
        args[0],
        ['object'],
        'invalid argument 0 (%s %*<100v): expected object',
      );
      if (!config.registry)
        throw new JPLRuntimeError(
          'invalid argument 0 (%s %*<100v): missing "registry"',
        );

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
