import type { FetchProviderType } from '@civet/common';
import { useConfigContext } from '@civet/core';
import { Field, Input, Label } from '@headlessui/react';
import { useCallback, useState, type SubmitEvent } from 'react';
import { useApplication } from '@/application';
import Button from '@/components/Button';
import Text from '@/components/Text';
import { setAccessToken } from '@/data/accessTokens';
import { postAuth } from '@/data/dataAuth';
import { HTTP_UNAUTHORIZED, HTTPError } from '@/data/errors';
import { notify } from '@/notifications/notifications';

export default function LoginForm() {
  const { dataProvider } = useConfigContext<FetchProviderType>();
  const { appName } = useApplication();

  const [userName, setUserName] = useState('');
  const [password, setPassword] = useState('');
  const [rememberLogin, setRememberLogin] = useState(true);
  const [loading, setLoading] = useState(false);

  const submit = useCallback(
    async (e: SubmitEvent<HTMLFormElement>) => {
      e.preventDefault();

      try {
        setLoading(true);
        const { accessToken } = await postAuth({
          dataProvider: dataProvider!,
          data: { userName, password },
        });
        setAccessToken(accessToken, rememberLogin);
        dataProvider!.notify('v1/user');
      } catch (error) {
        if (error instanceof HTTPError && error.status === HTTP_UNAUTHORIZED) {
          notify.error('Invalid credentials. Please try again.');
        } else {
          console.error(error);
          notify.error('Something went wrong. Please try again.');
        }
      } finally {
        setLoading(false);
      }
    },
    [dataProvider, userName, password, rememberLogin],
  );

  return (
    <form
      className="login form"
      onSubmit={submit}
      data-loading={loading ? '' : undefined}
    >
      <Text className="title" as="h1" text={`Login to ${appName}`} />

      <Field className="user-name field">
        <Text as={Label as unknown as 'label'} className="label" text="User" />

        <Input
          className="input"
          type="text"
          name="user-name"
          required
          autoFocus
          autoComplete="username"
          autoCapitalize="none"
          value={userName}
          onChange={(event) => {
            setUserName(event.target.value);
          }}
        />
      </Field>

      <Field className="password field">
        <Text
          as={Label as unknown as 'label'}
          className="label"
          text="Password"
        />

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

      <Field className="remember-login field">
        <Text
          as={Label as unknown as 'label'}
          className="label"
          text="Remember me"
        />

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

      <div className="buttons">
        <Button className="submit" type="submit" text="Login" />
      </div>
    </form>
  );
}
