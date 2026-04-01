import {
  createContext,
  useContext,
  useMemo,
  useState,
  type Dispatch,
  type SetStateAction,
} from 'react';
import {
  useHealthStatuses,
  type HealthStatus,
} from '@/data/dataHealthStatuses';
import { useLayout, type Layout } from '@/data/dataLayout';
import { useServices, type Service } from '@/data/dataServices';
import type { User } from '@/data/dataUser';
import { useVersions, type Version } from '@/data/dataVersions';
import { useErrorNotification } from '@/data/errors';

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
  setUnhide: Dispatch<SetStateAction<boolean>>;
  layout?: Layout;
  services: Service[];
  servicesByID: Partial<{ [serviceID: string]: Service }>;
  healthStatuses: HealthStatus[];
  healthStatusesByServiceID: Partial<{
    [serviceID: string]: HealthStatus[];
  }>;
  versions: Version[];
  versionsByServiceID: Partial<{ [serviceID: string]: Version[] }>;
};

export function useApplicationContextState(
  user: User | undefined,
): Application {
  const authorized = user != null;
  const admin = user?.role === 'admin';

  const [unhide, setUnhide] = useState(false);

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

  const healthStatuses = useHealthStatuses({
    disabled: !authorized,
    service: { config: { hidden: unhide ? undefined : false } },
  });
  useErrorNotification(healthStatuses);
  const healthStatusesByServiceID = useMemo(
    () =>
      Object.groupBy(
        healthStatuses.data ?? [],
        (healthStatus) => healthStatus.serviceID,
      ),
    [healthStatuses.data],
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
      setUnhide,
      layout: layout.data,
      services: services.data ?? [],
      servicesByID,
      healthStatuses: healthStatuses.data ?? [],
      healthStatusesByServiceID,
      versions: versions.data ?? [],
      versionsByServiceID,
    }),
    [
      authorized,
      user,
      admin,
      unhide,
      setUnhide,
      layout.data,
      services.data,
      servicesByID,
      healthStatuses,
      healthStatusesByServiceID,
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
