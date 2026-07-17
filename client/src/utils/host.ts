export function isAdminHost(): boolean {
  if (typeof window === 'undefined') return false;
  if (process.env.REACT_APP_FORCE_ADMIN_VIEW === 'true') return true;
  return window.location.hostname.startsWith('admin.');
}

export function getNormalSiteUrl(): string {
  if (typeof window === 'undefined') return '/';

  const { protocol, hostname, port } = window.location;

  if (hostname === 'admin.localhost' || hostname === 'admin.127.0.0.1') {
    return 'http://localhost:3000';
  }

  if (hostname.startsWith('admin.')) {
    const normalHostname = hostname.slice(6);
    const normalPort = port ? `:${port}` : '';
    return `${protocol}//${normalHostname}${normalPort}`;
  }

  return window.location.origin;
}

export function getAdminSiteUrl(): string {
  if (typeof window === 'undefined') return '/';

  const adminUrl = process.env.REACT_APP_ADMIN_URL;
  if (adminUrl) return adminUrl;

  const { protocol, hostname, port } = window.location;

  if (hostname === 'localhost' || hostname === '127.0.0.1') {
    return 'http://admin.localhost:3000';
  }

  const adminHostname = `admin.${hostname}`;
  const adminPort = port ? `:${port}` : '';
  return `${protocol}//${adminHostname}${adminPort}`;
}
