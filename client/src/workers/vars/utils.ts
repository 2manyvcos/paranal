import { JPLTypeError, type JPLRuntime } from '@jplorg/jpl';

export function assertType<T>(
  runtime: JPLRuntime,
  arg: unknown,
  types: string[],
  message: ((type: string, value: unknown) => Error) | string,
): T {
  const value = runtime.unwrapValue(arg ?? null);
  const t = runtime.type(value);
  if (!types.includes(t)) {
    throw typeof message === 'function'
      ? message(t, value)
      : new JPLTypeError(message, t, value);
  }
  return value as T;
}

export function assertJSON<T>(
  runtime: JPLRuntime,
  arg: unknown,
  types: string[],
  message: ((type: string, value: unknown) => Error) | string,
): T {
  const value = runtime.stripJSON(arg ?? null);
  const t = runtime.type(value);
  if (!types.includes(t)) {
    throw typeof message === 'function'
      ? message(t, value)
      : new JPLTypeError(message, t, value);
  }
  return value as T;
}
