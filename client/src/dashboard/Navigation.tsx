import type { FetchProviderType } from '@civet/common';
import { useConfigContext } from '@civet/core';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { useState } from 'react';
import { Link } from 'react-router';
import { useApplication } from '@/application';
import { unsetAccessToken } from '@/data/accessTokens';

export default function Navigation() {
  const { dataProvider } = useConfigContext<FetchProviderType>();
  const { appName, logo, userName, favorites, logoutRedirectURL } =
    useApplication();

  const [active, setActive] = useState(false);

  return (
    <nav className="navigation" data-active={active ? '' : undefined}>
      <Link className="title" to="/?no-redirect">
        <div className="logo">
          <img className="image" src={logo} />
        </div>

        <h1 className="app-name">
          <span className="text">{appName}</span>
        </h1>
      </Link>

      <div className="burger-container">
        <a
          className="burger"
          role="button"
          onClick={() => {
            setActive((prev) => !prev);
          }}
        />
      </div>

      <ul className="menu">
        <li className="start menu-item">
          <ul className="sub menu">
            <li className="favorites menu-item" data-count={favorites.length}>
              <Link className="item" to="/favorites">
                <span className="text">Favorites</span>
              </Link>
            </li>

            <li className="page menu-item">
              <Link className="item" to="/services">
                <span className="text">Services</span>
              </Link>
            </li>
          </ul>
        </li>

        <li className="end menu-item">
          <ul className="sub menu">
            <li className="user menu-item">
              <Menu>
                <MenuButton className="item" as="a" role="button">
                  <span className="text">{userName}</span>
                </MenuButton>

                <MenuItems className="user dropdown" anchor="bottom">
                  <MenuItem>
                    <a
                      className="logout dropdown-item"
                      href={logoutRedirectURL || undefined}
                      onClick={() => {
                        unsetAccessToken();
                        dataProvider!.notify('v1/user');
                      }}
                    >
                      <span className="text">Logout</span>
                    </a>
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
