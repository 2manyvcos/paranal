import type { FetchProviderType } from '@civet/common';
import { useConfigContext } from '@civet/core';
import { Button, Field, Input, Label } from '@headlessui/react';
import { useCallback, useState, type SubmitEvent } from 'react';
import { setAccessToken } from '@/data/accessTokens';
import type { Auth } from '@/data/auth';
import { HTTP_UNAUTHORIZED, HTTPError } from '@/data/errors';
import { notify } from '@/notifications/notification';

export default function LoginPage() {
  const { dataProvider } = useConfigContext<FetchProviderType>();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [rememberLogin, setRememberLogin] = useState(true);
  const [loading, setLoading] = useState(false);

  const login = useCallback(
    async (e: SubmitEvent<HTMLFormElement>) => {
      e.preventDefault();

      try {
        setLoading(true);
        const { accessToken } = await dataProvider!.request<Auth>(
          'v1/auth',
          {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
          },
          {
            async handleError(_url, _request, response, _meta) {
              throw new HTTPError(await response.text(), response.status);
            },
          },
        );
        setAccessToken(accessToken, rememberLogin);
        dataProvider!.notify('v1/user');
      } catch (error) {
        if (error instanceof HTTPError && error.status === HTTP_UNAUTHORIZED) {
          notify('Invalid credentials. Please try again.');
        } else {
          console.error(error);
          notify('Something went wrong. Please try again.');
        }
      } finally {
        setLoading(false);
      }
    },
    [dataProvider, username, password, rememberLogin],
  );

  return (
    <div className="login page">
      <form
        className="login form"
        onSubmit={login}
        data-loading={loading ? '' : undefined}
      >
        <Field className="username">
          <Label className="label" />

          <Input
            className="input"
            type="text"
            name="username"
            required
            autoFocus
            autoComplete="username"
            autoCapitalize="none"
            value={username}
            onChange={(event) => {
              setUsername(event.target.value);
            }}
          />
        </Field>

        <Field className="password">
          <Label className="label" />

          <Input
            className="input"
            type="password"
            name="password"
            required
            autoComplete="current-password"
            autoCapitalize="none"
            value={password}
            onChange={(event) => {
              setPassword(event.target.value);
            }}
          />
        </Field>

        <Field className="remember-login">
          <Label className="label" />

          <Input
            className="input"
            type="checkbox"
            name="remember-login"
            checked={rememberLogin}
            onChange={(event) => {
              setRememberLogin(event.target.checked);
            }}
          />
        </Field>

        <Button className="submit button" type="submit" />
      </form>
    </div>
  );
}
