declare global {
  interface Window {
    paranal: {
      appName: string;
      tagline: string;
      logo: string;
      logoutRedirectURL: string;
    };
  }
}
