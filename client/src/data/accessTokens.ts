const accessTokenKey = 'paranal-access-token';

export function getAccessToken(): string {
  return (
    sessionStorage.getItem(accessTokenKey) ??
    localStorage.getItem(accessTokenKey) ??
    ''
  );
}

export function setAccessToken(accessToken: string, rememberLogin: boolean) {
  unsetAccessToken();
  (rememberLogin ? localStorage : sessionStorage).setItem(
    accessTokenKey,
    accessToken,
  );
}

export function unsetAccessToken() {
  localStorage.removeItem(accessTokenKey);
  sessionStorage.removeItem(accessTokenKey);
}
