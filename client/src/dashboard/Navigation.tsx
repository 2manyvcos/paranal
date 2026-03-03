import type { FetchProviderType } from '@civet/common';
import { useConfigContext } from '@civet/core';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { useState } from 'react';
import { Link as RouterLink } from 'react-router';
import { useApplication } from '@/application';
import Image from '@/components/Image';
import Link from '@/components/Link';
import NavItemLink from '@/components/NavItemLink';
import Text from '@/components/Text';
import { unsetAccessToken } from '@/data/accessTokens';

export default function Navigation() {
  const { dataProvider } = useConfigContext<FetchProviderType>();
  const {
    appName,
    tagline,
    logo,
    userName,
    unhide,
    setUnhide,
    layout,
    logoutRedirectURL,
    healthStatuses,
    versions,
  } = useApplication();

  const [active, setActive] = useState(false);

  return (
    <nav className="navigation" data-active={active ? '' : undefined}>
      <header className="header">
        <RouterLink className="header-link" to="/">
          <Image className="logo" image={logo} />

          <div className="title">
            <Text className="name" as="h1" text={appName} />

            <Text className="tagline" text={tagline} />
          </div>
        </RouterLink>
      </header>

      <div className="burger">
        <a
          className="button"
          role="button"
          tabIndex={0}
          onClick={() => {
            setActive((prev) => !prev);
          }}
        >
          <span className="content">&#x2630;</span>
        </a>
      </div>

      <ul className="menu">
        <li className="start menu-item">
          <ul className="sub menu">
            <li className="home menu-item">
              <NavItemLink to="/?no-redirect" text="Home" />
            </li>

            {layout?.pages?.map((page) =>
              !page.name ? null : (
                <li key={page.name} className="page menu-item">
                  <NavItemLink
                    to={`/${encodeURIComponent(page.name)}`}
                    text={page.displayName || page.name}
                  />
                </li>
              ),
            )}
          </ul>
        </li>

        <li className="end menu-item">
          <ul className="sub menu">
            <li className="health-statuses menu-item">
              <NavItemLink
                to="/health-statuses"
                text="Health"
                data-total={healthStatuses.length || undefined}
                data-combined={
                  healthStatuses.filter(
                    (healthStatus) => healthStatus.unhealthy,
                  ).length || undefined
                }
                data-unhealthy={
                  healthStatuses.filter(
                    (healthStatus) => healthStatus.unhealthy,
                  ).length || undefined
                }
              />
            </li>

            <li className="versions menu-item">
              <NavItemLink
                to="/versions"
                text="Versions"
                data-total={versions.length || undefined}
                data-combined={
                  versions.filter(
                    (version) => version.outdated || version.vulnerable,
                  ).length || undefined
                }
                data-outdated={
                  versions.filter((version) => version.outdated).length ||
                  undefined
                }
                data-vulnerable={
                  versions.filter((version) => version.vulnerable).length ||
                  undefined
                }
              />
            </li>

            <li className="user menu-item">
              <Menu>
                <MenuButton
                  className="nav-item link"
                  as="a"
                  role="button"
                  tabIndex={0}
                  data-text={userName || undefined}
                >
                  <span className="content">{userName || undefined}</span>
                </MenuButton>

                <MenuItems className="user dropdown" anchor="bottom">
                  <MenuItem>
                    <Link
                      className="unhide dropdown-item"
                      as="a"
                      tabIndex={0}
                      onClick={() => {
                        setUnhide((prev) => !prev);
                      }}
                      text={
                        unhide
                          ? 'Stop showing hidden services'
                          : 'Show hidden services'
                      }
                      data-active={unhide ? '' : undefined}
                    />
                  </MenuItem>

                  <div className="dropdown-separator" />

                  <MenuItem>
                    <Link
                      className="logout dropdown-item"
                      as="a"
                      tabIndex={0}
                      href={logoutRedirectURL || undefined}
                      onClick={() => {
                        unsetAccessToken();
                        dataProvider!.notify('v1/user');
                      }}
                      text="Logout"
                    />
                  </MenuItem>
                </MenuItems>
              </Menu>
            </li>
          </ul>
        </li>
      </ul>
    </nav>
  );
}
