import { createContext, useContext } from 'react';

type Application = {
  appName: string;
  logo: string;
  apiURL: string;

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
    apiURL: import.meta.env.PARANAL_API || '/api',

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
