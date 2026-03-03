import { Menu, MenuButton, MenuItems } from '@headlessui/react';
import { useApplication } from '@/application';
import Image from '@/components/Image';
import Link from '@/components/Link';
import Markdown from '@/components/Markdown';
import Text from '@/components/Text';
import type { Page } from '@/data/data-layout';

export default function ServicePage({ page }: { page: Page }) {
  const {
    appName,
    servicesByID,
    uptimeStatusesByServiceID,
    versionsByServiceID,
  } = useApplication();

  return (
    <div className="service page">
      <title>{`${page.displayName || page.name} | ${appName}`}</title>

      <header className="header">
        <Text className="name" as="h2" text={page.displayName || page.name} />

        <Markdown className="header" markdown={page.header} />
      </header>

      <main className="sections">
        {page.sections?.map((section) =>
          !section.id ? null : (
            <section key={section.id} className="section">
              <header className="header">
                <Image className="icon" image={section.icon} />

                <div className="title">
                  <Text className="name" as="h3" text={section.name} />
                </div>
              </header>

              <div className="services">
                {section.serviceIDs?.map((serviceID) => {
                  const service = servicesByID[serviceID];
                  if (!service) return null;
                  const uptimeStatuses =
                    uptimeStatusesByServiceID[serviceID] ?? [];
                  const unhealthyUptimeStatuses = uptimeStatuses.filter(
                    (uptimeStatus) => uptimeStatus.unhealthy,
                  );
                  const versions = versionsByServiceID[serviceID] ?? [];
                  const outdatedVersions = versions.filter(
                    (version) => version.outdated,
                  );
                  const vulnerableVersions = versions.filter(
                    (version) => version.vulnerable,
                  );
                  return (
                    <article
                      key={serviceID}
                      className="service"
                      id={`service:${serviceID}`}
                    >
                      <header className="header">
                        <a
                          className="header-link"
                          href={service.url || undefined}
                        >
                          <Image className="logo" image={service.logo} />

                          <div className="title">
                            <Text
                              className="name"
                              as="h4"
                              text={service.name}
                            />

                            <Text
                              className="description"
                              text={service.description}
                            />
                          </div>
                        </a>

                        <div className="controls">
                          <div className="badges control">
                            <div className="badges">
                              <Text
                                className="favorite badge"
                                text={
                                  service.config?.favorite
                                    ? 'favorite'
                                    : undefined
                                }
                              />

                              <Text
                                className="hidden badge"
                                text={
                                  service.config?.hidden ? 'hidden' : undefined
                                }
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

                              <MenuItems
                                className="service context dropdown"
                                anchor="bottom"
                              >
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
                            pathname: '/uptimestatuses',
                            search: new URLSearchParams({
                              serviceID,
                              unhealthy: 'true',
                            }).toString(),
                            hash: `service:${serviceID}`,
                          }}
                          text={
                            uptimeStatuses.length
                              ? unhealthyUptimeStatuses.length
                                ? 'unhealthy'
                                : 'healthy'
                              : undefined
                          }
                          data-total={uptimeStatuses.length || undefined}
                          data-count={
                            unhealthyUptimeStatuses.length || undefined
                          }
                        />

                        <Link
                          className="outdated badge"
                          to={{
                            pathname: '/versions',
                            search: new URLSearchParams({
                              serviceID,
                              outdated: 'true',
                            }).toString(),
                            hash: `service:${serviceID}`,
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
                              serviceID,
                              vulnerable: 'true',
                            }).toString(),
                            hash: `service:${serviceID}`,
                          }}
                          text={
                            vulnerableVersions.length ? 'vulnerable' : undefined
                          }
                          data-total={versions.length || undefined}
                          data-count={vulnerableVersions.length || undefined}
                          data-cves={
                            versions.reduce(
                              (sum, version) => sum + version.currentCVEs,
                              0,
                            ) || undefined
                          }
                        />
                      </div>
                    </article>
                  );
                })}
              </div>
            </section>
          ),
        )}
      </main>
    </div>
  );
}
