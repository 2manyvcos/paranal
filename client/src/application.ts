import { createContext, useContext, useMemo } from 'react';
import { useErrorNotification } from './data/errors';
import { useLayout, type Layout } from './data/layout';
import { useServices, type Service } from './data/services';
import type { User } from './data/user';

export type Application = {
  appName: string;
  logo: string;
  footer: string;
  logoutRedirectURL: string;

  authorized: boolean;
  user?: User;
  admin: boolean;

  unhide: boolean;
  layout?: Layout;
  services: Service[];
  favorites: Service[];
};

export function useApplicationContextState(
  user: User | undefined,
): Application {
  const authorized = user != null;
  const admin = user?.role === 'admin';

  const unhide = false;

  const layout = useLayout({ disabled: !authorized });
  useErrorNotification(layout);

  const services = useServices({
    config: { hidden: unhide ? undefined : false },
    disabled: !authorized,
  });
  useErrorNotification(services);

  return useMemo(
    () => ({
      appName: window.paranal.appName,
      logo: window.paranal.logo,
      footer: window.paranal.footer,
      logoutRedirectURL: window.paranal.logoutRedirectURL,

      authorized,
      user,
      admin,

      unhide,
      layout: layout.data,
      services: services.data ?? [],
      favorites:
        services.data?.filter((service) => service.config.favorite) ?? [],
    }),
    [authorized, user, admin, unhide, layout.data, services.data],
  );
}

export const ApplicationContext = createContext<Application>(
  null as unknown as Application,
);

export function useApplication(): Application {
  return useContext(ApplicationContext);
}
