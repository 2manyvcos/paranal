export function paranalReconnectInterest(): string {
  return `paranal://reconnect`;
}

export function paranalStatesInterest(registry: string): string {
  return `paranal://registry/${encodeURIComponent(registry)}/states`;
}

export function paranalStateInterest(registry: string, id: string): string {
  return `paranal://registry/${encodeURIComponent(registry)}/states/${encodeURIComponent(id)}`;
}
