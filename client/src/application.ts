import { createContext, useContext } from 'react';

type Application = {
  appName: string;
  logo: string;

  authorized: boolean;
  user?: {
    startPage: string;
    admin: boolean;
  };
};

export function useApplicationState(): Application {
  return {
    appName: window.paranal.appName,
    logo: window.paranal.logo,

    authorized: true,
    user: {
      startPage: '',
      admin: false,
    },
  };
}

export const ApplicationContext = createContext<Application>(
  null as unknown as Application,
);

export function useApplication(): Application {
  return useContext(ApplicationContext);
}
