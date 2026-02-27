import type { FetchProviderType } from '@civet/common';
import { useConfigContext } from '@civet/core';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { useState } from 'react';
import { Link as RouterLink } from 'react-router';
import { useApplication } from '@/application';
import Image from '@/components/Image';
import Link from '@/components/Link';
import NavItem from '@/components/NavItem';
import Text from '@/components/Text';
import { unsetAccessToken } from '@/data/accessTokens';

export default function Navigation() {
  const { dataProvider } = useConfigContext<FetchProviderType>();
  const { appName, tagline, logo, userName, layout, logoutRedirectURL } =
    useApplication();

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
          onClick={() => {
            setActive((prev) => !prev);
          }}
        >
          <span className="content" />
        </a>
      </div>

      <ul className="menu">
        <li className="start menu-item">
          <ul className="sub menu">
            <li className="home menu-item">
              <NavItem to="/?no-redirect" text="Home" />
            </li>

            {layout?.pages?.map((page) =>
              !page.name ? null : (
                <li key={page.name} className="page menu-item">
                  <NavItem
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
            <li className="user menu-item">
              <Menu>
                <MenuButton
                  className="nav-item"
                  as="a"
                  role="button"
                  data-text={userName || undefined}
                >
                  <span className="content">{userName || undefined}</span>
                </MenuButton>

                <MenuItems className="user dropdown" anchor="bottom">
                  <MenuItem>
                    <Link
                      className="logout dropdown-item"
                      as="a"
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
