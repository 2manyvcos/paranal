import type { FetchProviderType } from '@civet/common';
import { useConfigContext } from '@civet/core';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/react';
import { Link } from 'react-router';
import { useApplication } from '@/application';
import { unsetAccessToken } from '@/data/accessTokens';

export default function Header() {
  const { dataProvider } = useConfigContext<FetchProviderType>();
  const application = useApplication();

  return (
    <header className="header">
      <div className="title-bar">
        <Link to="/?no-redirect" className="title">
          <div className="logo" />

          <div className="app-name" />
        </Link>
      </div>

      <nav className="navigation">
        <div className="start"></div>

        <div className="end">
          <Menu>
            <MenuButton className="user item">
              {application.user?.displayName || application.user?.name || ''}
            </MenuButton>

            <MenuItems anchor="bottom" className="user menu">
              <MenuItem>
                <a
                  className="logout item"
                  href={application.logoutRedirectURL || undefined}
                  onClick={() => {
                    unsetAccessToken();
                    dataProvider!.notify('v1/user');
                  }}
                />
              </MenuItem>
            </MenuItems>
          </Menu>
        </div>
      </nav>
    </header>
  );
}
