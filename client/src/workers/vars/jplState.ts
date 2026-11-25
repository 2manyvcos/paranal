import { HTTP_NOT_FOUND, type HTTPError } from '@/api/errors';
import {
  paranalReconnectInterest,
  paranalStateInterest,
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
      _input: unknown,
      ...args: unknown[]
    ): Promise<unknown[]> => {
      if (args.length < 1) throw new JPLRuntimeError('not enough arguments');
      if (args.length > 1) throw new JPLRuntimeError('too many arguments');

      const config = assertJSON<StateConfig>(
        runtime,
        args[0],
        ['object'],
        'invalid argument 0 (%s %*<100v): expected object',
      );
      if (!config.registry)
        throw new JPLRuntimeError(
          'invalid argument 0 (%s %*<100v): missing "registry"',
        );
      if (!config.id)
        throw new JPLRuntimeError(
          'invalid argument 0 (%s %*<100v): missing "id"',
        );

      reducer.addInterest(paranalReconnectInterest());

      const stateInterest = paranalStateInterest(config.registry, config.id);
      requestedInterests.add(stateInterest);
      reducer.addInterest(stateInterest);

      const abortSignal = new AbortSignal();
      signal.subscribe(() => {
        abortSignal.abort();
      });

      try {
        const response = await provider.get(
          `v1/registries/${encodeURIComponent(config.registry)}/states/${encodeURIComponent(config.id)}`,
          undefined,
          undefined,
          undefined,
          abortSignal,
        );

        return next(response);
      } catch (error) {
        if ((error as HTTPError).status === HTTP_NOT_FOUND) {
          return next(null);
        }

        throw error;
      }
    };
  };
}
