let _apiURL: URL;
export function getAPIURL(): URL {
  if (_apiURL == null)
    _apiURL = new URL(
      window.paranal.api.replace(/\/*$/, '/'),
      window.location.href,
    );

  return _apiURL;
}

export function normalizePath(path: string): string {
  const apiURL = getAPIURL();

  return new URL(path, apiURL).pathname
    .substring(apiURL.pathname.length)
    .replace(/\/*$/, '')
    .split('/')
    .map((segment) => encodeURIComponent(decodeURIComponent(segment)))
    .join('/');
}

export function matchEvent(pattern: string, path: string): boolean {
  const p = pattern.split('/');
  const e = path.split('/');

  if (p.length !== e.length) return false;

  const differs = p.some((ps, i) => {
    const es = e[i];
    return ps != '*' && ps != es;
  });

  return !differs;
}
