import type { FetchProviderType } from '@civet/common';
import { useResource } from '@civet/core';

export type LayoutQuery = {
  disabled?: boolean;
};

export type Layout = {
  pages?: Page[];
};

export type Page = {
  name?: string;
  displayName?: string;
  header?: string;

  sections?: Section[];
};

export type Section = {
  id?: string;
  name?: string;
  icon?: string;
  serviceIDs?: string[];
};

export function useLayout(query: LayoutQuery) {
  return useResource<FetchProviderType, Layout | undefined>({
    name: 'v1/layout',
    query: undefined,
    events: true,
    disabled: query.disabled,
  });
}
