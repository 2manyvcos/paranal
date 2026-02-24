declare global {
  interface Window {
    paranal: {
      appName: string;
      tagline: string;
      logo: string;
      footer: string;
      logoutRedirectURL: string;
    };
  }
}
