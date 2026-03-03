import { Menu, MenuButton, MenuItems } from '@headlessui/react';
import { useApplication } from '@/application';
import Image from '@/components/Image';
import Link from '@/components/Link';
import Text from '@/components/Text';
import type { Service } from '@/data/data-services';

export default function ServiceCard({ service }: { service: Service }) {
  const { healthStatusesByServiceID, versionsByServiceID } = useApplication();

  const healthStatuses = healthStatusesByServiceID[service.id] ?? [];
  const unhealthyHealthStatuses = healthStatuses.filter(
    (healthStatus) => healthStatus.unhealthy,
  );
  const versions = versionsByServiceID[service.id] ?? [];
  const outdatedVersions = versions.filter((version) => version.outdated);
  const vulnerableVersions = versions.filter((version) => version.vulnerable);

  return (
    <article className="service card" id={`service:${service.id}`}>
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
                text={service.config?.favorite ? 'favorite' : undefined}
              />

              <Text
                className="hidden badge"
                text={service.config?.hidden ? 'hidden' : undefined}
              />
            </div>
          </div>

          <div className="context control">
            <Menu>
              <MenuButton
                className="context-menu link"
                as="a"
                role="button"
                data-text={'\u22ee;'}
              >
                <span className="content">&#x22ee;</span>
              </MenuButton>

              <MenuItems className="service context dropdown" anchor="bottom">
                {/* TODO: <MenuItem>
                  <Link
                    className="edit dropdown-item"
                    as="a"
                    onClick={() => {}}
                    text="Edit"
                  />
                </MenuItem> */}
              </MenuItems>
            </Menu>
          </div>
        </div>
      </header>

      <div className="badges">
        <Link
          className="unhealthy badge"
          to={{
            pathname: '/healthstatuses',
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
