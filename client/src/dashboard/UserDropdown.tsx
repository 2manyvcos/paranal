import type { FetchProviderType } from '@civet/common';
import { useConfigContext } from '@civet/core';
import { MenuItem, MenuItems, MenuSeparator } from '@headlessui/react';
import { useApplication } from '@/application';
import Link from '@/components/Link';
import { unsetAccessToken } from '@/data/accessTokens';

export default function UserDropdown() {
  const { dataProvider } = useConfigContext<FetchProviderType>();
  const { admin, unhide, setUnhide, logoutRedirectURL } = useApplication();

  return (
    <MenuItems className="user dropdown" anchor="bottom">
      <MenuItem>
        <Link className="user dropdown-item" to="/user" text="Settings" />
      </MenuItem>

      {!admin ? null : (
        <MenuItem>
          <Link
            className="admin dropdown-item"
            to="/admin"
            text="Admin Panel"
          />
        </MenuItem>
      )}

      <MenuSeparator className="dropdown-separator" />

      <MenuItem>
        <Link
          className="unhide dropdown-item"
          as="a"
          onClick={() => {
            setUnhide((prev) => !prev);
          }}
          text={
            unhide ? 'Stop showing hidden services' : 'Show hidden services'
          }
          data-enabled={unhide ? '' : undefined}
        />
      </MenuItem>

      <MenuSeparator className="dropdown-separator" />

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
  );
}
