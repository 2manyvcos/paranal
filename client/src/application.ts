import { createContext, useContext, useMemo } from 'react';
import { useLayout, type Layout } from './data/data-layout';
import { useServices, type Service } from './data/data-services';
import type { User } from './data/data-user';
import { useErrorNotification } from './data/errors';

export type Application = {
  appName: string;
  tagline: string;
  logo: string;
  footer: string;
  logoutRedirectURL: string;

  authorized: boolean;
  user?: User;
  userName?: string;
  admin: boolean;

  unhide: boolean;
  layout?: Layout;
  services: Service[];
  servicesByID: { [serviceID: string]: Service };
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
  const servicesByID = useMemo(
    () =>
      Object.fromEntries(
        services.data?.map((service) => [service.id, service]) ?? [],
      ),
    [services.data],
  );

  return useMemo(
    () => ({
      appName: window.paranal.appName,
      tagline: window.paranal.tagline,
      logo: window.paranal.logo,
      footer: window.paranal.footer,
      logoutRedirectURL: window.paranal.logoutRedirectURL,

      authorized,
      user,
      userName: user?.displayName || user?.name,
      admin,

      unhide,
      layout: layout.data,
      services: services.data ?? [],
      servicesByID,
      favorites:
        services.data?.filter((service) => service.config.favorite) ?? [],
    }),
    [authorized, user, admin, unhide, layout.data, services.data, servicesByID],
  );
}

export const ApplicationContext = createContext<Application>(
  null as unknown as Application,
);

export function useApplication(): Application {
  return useContext(ApplicationContext);
}
