import { SSEReceiver, type SSEOptions } from '@civet/common';
import type { GenericDataProvider, ResourceContextValue } from '@civet/core';

export type ParanalEventOptions = SSEOptions<
  ResourceContextValue<GenericDataProvider>,
  MessageEvent
> & {
  pathPattern?: string;
};

export class ParanalReceiver extends SSEReceiver<
  ResourceContextValue<GenericDataProvider>,
  MessageEvent,
  ParanalEventOptions
> {}

export type ParanalReceiverType = InstanceType<typeof ParanalReceiver>;
