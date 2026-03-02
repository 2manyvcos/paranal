import { createContext, useContext, useMemo } from 'react';
import { useLayout, type Layout } from './data/data-layout';
import { useServices, type Service } from './data/data-services';
import {
  useUptimeStatuses,
  type UptimeStatus,
} from './data/data-uptimestatuses';
import type { User } from './data/data-user';
import { useVersions, type Version } from './data/data-versions';
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
  servicesByID: Partial<{ [serviceID: string]: Service }>;
  servicesWithFavorite: Service[];
  uptimeStatuses: UptimeStatus[];
  uptimeStatusesByServiceID: Partial<{
    [serviceID: string]: UptimeStatus[];
  }>;
  versions: Version[];
  versionsByServiceID: Partial<{ [serviceID: string]: Version[] }>;
};

export function useApplicationContextState(
  user: User | undefined,
): Application {
  const authorized = user != null;
  const admin = user?.role === 'admin';

  const unhide = true;

  const layout = useLayout({ disabled: !authorized });
  useErrorNotification(layout);

  const services = useServices({
    disabled: !authorized,
    config: { hidden: unhide ? undefined : false },
  });
  useErrorNotification(services);
  const servicesByID = useMemo(
    () =>
      Object.fromEntries(
        services.data?.map((service) => [service.id, service]) ?? [],
      ),
    [services.data],
  );

  const uptimeStatuses = useUptimeStatuses({
    disabled: !authorized,
    service: { config: { hidden: unhide ? undefined : false } },
  });
  useErrorNotification(uptimeStatuses);
  const uptimeStatusesByServiceID = useMemo(
    () =>
      Object.groupBy(
        uptimeStatuses.data ?? [],
        (uptimeStatus) => uptimeStatus.serviceID,
      ),
    [uptimeStatuses.data],
  );

  const versions = useVersions({
    disabled: !authorized,
    service: { config: { hidden: unhide ? undefined : false } },
  });
  useErrorNotification(versions);
  const versionsByServiceID = useMemo(
    () => Object.groupBy(versions.data ?? [], (version) => version.serviceID),
    [versions.data],
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
      servicesWithFavorite:
        services.data?.filter((service) => service.config.favorite) ?? [],
      uptimeStatuses: uptimeStatuses.data ?? [],
      uptimeStatusesByServiceID,
      versions: versions.data ?? [],
      versionsByServiceID,
    }),
    [
      authorized,
      user,
      admin,
      unhide,
      layout.data,
      services.data,
      servicesByID,
      uptimeStatuses,
      uptimeStatusesByServiceID,
      versions,
      versionsByServiceID,
    ],
  );
}

export const ApplicationContext = createContext<Application>(
  null as unknown as Application,
);

export function useApplication(): Application {
  return useContext(ApplicationContext);
}
