import { Menu, MenuButton } from '@headlessui/react';
import { useLocation } from 'react-router';
import { useApplication } from '@/application';
import Image from '@/components/Image';
import Link from '@/components/Link';
import Text from '@/components/Text';
import { type Service } from '@/data/dataServices';
import ServiceMenu from './ServiceMenu';

export default function ServiceCard({
  service,
  anchor = true,
}: {
  service: Service;
  anchor?: boolean;
}) {
  const { hash } = useLocation();
  const { healthStatusesByServiceID, versionsByServiceID } = useApplication();

  const healthStatuses = healthStatusesByServiceID[service.id] ?? [];
  const unhealthyHealthStatuses = healthStatuses.filter(
    (healthStatus) => healthStatus.unhealthy,
  );
  const versions = versionsByServiceID[service.id] ?? [];
  const outdatedVersions = versions.filter((version) => version.outdated);
  const vulnerableVersions = versions.filter((version) => version.vulnerable);

  return (
    <article
      className="service card"
      id={anchor ? `service:${service.id}` : undefined}
      data-active={anchor && hash === `#service:${service.id}` ? '' : undefined}
    >
      <header className="header">
        <a className="header-link" href={service.url || undefined}>
          <Image className="logo" image={service.logo} />

          <div className="title">
            <Text className="name" as="h4" text={service.name} />

            <Text className="description" text={service.description} />
          </div>
        </a>

        <div className="controls">
          <div className="badges control">
            <div className="badges">
              <Text
                className="favorite badge"
                text={service.config.favorite ? 'favorite' : undefined}
                data-enabled={service.config.favorite || undefined}
              />

              <Text
                className="hidden badge"
                text={service.config.hidden ? 'hidden' : undefined}
                data-enabled={service.config.hidden || undefined}
              />
            </div>
          </div>

          <div className="context control">
            <Menu>
              {({ open }) => (
                <>
                  <MenuButton
                    className="context-menu link"
                    as="a"
                    role="button"
                    tabIndex={0}
                    data-text={'\u22ee;'}
                  >
                    <span className="content">&#x22ee;</span>
                  </MenuButton>

                  <ServiceMenu service={service} open={open} />
                </>
              )}
            </Menu>
          </div>
        </div>
      </header>

      <div className="badges">
        <Link
          className="unhealthy badge"
          to={{
            pathname: '/health-statuses',
            search: new URLSearchParams({
              serviceID: service.id,
              unhealthy: 'true',
            }).toString(),
            hash: `service:${service.id}`,
          }}
          text={
            healthStatuses.length
              ? unhealthyHealthStatuses.length
                ? 'unhealthy'
                : 'healthy'
              : undefined
          }
          data-total={healthStatuses.length || undefined}
          data-count={unhealthyHealthStatuses.length || undefined}
        />

        <Link
          className="outdated badge"
          to={{
            pathname: '/versions',
            search: new URLSearchParams({
              serviceID: service.id,
              outdated: 'true',
            }).toString(),
            hash: `service:${service.id}`,
          }}
          text={
            versions.length
              ? outdatedVersions.length
                ? 'outdated'
                : 'up to date'
              : undefined
          }
          data-total={versions.length || undefined}
          data-count={outdatedVersions.length || undefined}
        />

        <Link
          className="vulnerable badge"
          to={{
            pathname: '/versions',
            search: new URLSearchParams({
              serviceID: service.id,
              vulnerable: 'true',
            }).toString(),
            hash: `service:${service.id}`,
          }}
          text={vulnerableVersions.length ? 'vulnerable' : undefined}
          data-total={versions.length || undefined}
          data-count={vulnerableVersions.length || undefined}
          data-cves={
            versions.reduce((sum, version) => sum + version.currentCVEs, 0) ||
            undefined
          }
        />
      </div>
    </article>
  );
}
