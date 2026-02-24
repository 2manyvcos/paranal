import { useEffect } from 'react';
import type { Application } from '@/application';

function useBooleanDataAttribute(name: string, value: boolean) {
  useEffect(() => {
    if (!value) return;
    const root = document.getElementById('root')!;
    root.setAttribute(name, '');
    return () => {
      root.removeAttribute(name);
    };
  }, [name, value]);
}

export function useCustomization(
  userLoading: boolean,
  application: Application,
) {
  useBooleanDataAttribute('data-loading', userLoading);

  useEffect(() => {
    const rootStyle = document.documentElement.style;
    // const root = document.getElementById('root')!;

    // root.setAttribute(name, value);

    rootStyle.setProperty(
      '--t-user-name',
      JSON.stringify(
        application.user?.displayName || application.user?.name || '',
      ),
    );
  }, [userLoading, application]);
}
