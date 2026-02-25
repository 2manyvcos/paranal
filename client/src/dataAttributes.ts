import { useEffect } from 'react';

export function useDataAttribute(
  element: HTMLElement,
  name: string,
  value: string | undefined,
) {
  useEffect(() => {
    if (element == null || value == null) return;
    const attributeName = `data-${name}`;
    element.setAttribute(attributeName, value);
    return () => {
      element.removeAttribute(attributeName);
    };
  }, [element, name, value]);
}

export function useBooleanDataAttribute(
  element: HTMLElement,
  name: string,
  value: boolean,
) {
  useDataAttribute(element, name, value ? '' : undefined);
}
