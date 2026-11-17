import { SSEReceiver, type SSEOptions } from '@civet/common';
import type { GenericDataProvider, ResourceContextValue } from '@civet/core';
import { matchEvent, normalizePath } from './api';

export type ParanalEventOptions = SSEOptions<
  ResourceContextValue<GenericDataProvider>,
  MessageEvent
> & {
  pathPatterns?: string[];
};

export default class ParanalReceiver extends SSEReceiver<
  ResourceContextValue<GenericDataProvider>,
  MessageEvent,
  ParanalEventOptions
> {
  constructor(eventSource: EventSource | undefined) {
    super(eventSource, {
      events: ['change'],

      getEvents(resource, type, event, options) {
        switch (type) {
          case 'change': {
            const eventPath = normalizePath(event.data);

            if (options?.pathPatterns != null) {
              if (
                options.pathPatterns.some((pathPattern) =>
                  matchEvent(pathPattern, eventPath),
                )
              ) {
                return [event];
              }
              return [];
            }

            if (!resource) return [event];

            const resourcePath = normalizePath(resource.name);
            if (eventPath === resourcePath) return [event];
            return [];
          }

          default:
            throw new Error(`unsupported event type "${type}"`);
        }
      },
    });
  }
}

export type ParanalReceiverType = InstanceType<typeof ParanalReceiver>;
