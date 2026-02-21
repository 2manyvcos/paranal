import { createContext, useContext } from 'react';
import { type User } from './data/user';

type Application = {
  appName: string;
  logo: string;

  authorized: boolean;
  user?: User;
  admin: boolean;
};

export function useApplicationState(user: User | undefined): Application {
  return {
    appName: window.paranal.appName,
    logo: window.paranal.logo,

    authorized: user != null,
    user: user,
    admin: user?.role === 'admin',
  };
}

export const ApplicationContext = createContext<Application>(
  null as unknown as Application,
);

export function useApplication(): Application {
  return useContext(ApplicationContext);
}
