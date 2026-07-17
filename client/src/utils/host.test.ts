import { isAdminHost, getNormalSiteUrl, getAdminSiteUrl } from './host';

describe('isAdminHost', () => {
  const originalHostname = window.location.hostname;

  afterEach(() => {
    Object.defineProperty(window, 'location', {
      value: { ...window.location, hostname: originalHostname },
      writable: true,
    });
    delete (process.env as any).REACT_APP_FORCE_ADMIN_VIEW;
  });

  it('returns false for normal hostnames', () => {
    Object.defineProperty(window, 'location', {
      value: { ...window.location, hostname: 'localhost' },
      writable: true,
    });
    expect(isAdminHost()).toBe(false);
  });

  it('returns true for admin subdomains', () => {
    Object.defineProperty(window, 'location', {
      value: { ...window.location, hostname: 'admin.example.com' },
      writable: true,
    });
    expect(isAdminHost()).toBe(true);
  });

  it('returns true when REACT_APP_FORCE_ADMIN_VIEW is true', () => {
    (process.env as any).REACT_APP_FORCE_ADMIN_VIEW = 'true';
    expect(isAdminHost()).toBe(true);
  });
});

describe('getNormalSiteUrl', () => {
  const originalLocation = window.location;

  afterEach(() => {
    Object.defineProperty(window, 'location', {
      value: originalLocation,
      writable: true,
    });
  });

  it('returns origin for normal host', () => {
    Object.defineProperty(window, 'location', {
      value: {
        protocol: 'https:',
        hostname: 'example.com',
        port: '',
        origin: 'https://example.com',
      },
      writable: true,
    });
    expect(getNormalSiteUrl()).toBe('https://example.com');
  });

  it('strips admin subdomain', () => {
    Object.defineProperty(window, 'location', {
      value: {
        protocol: 'https:',
        hostname: 'admin.example.com',
        port: '',
        origin: 'https://admin.example.com',
      },
      writable: true,
    });
    expect(getNormalSiteUrl()).toBe('https://example.com');
  });

  it('maps admin.localhost to localhost:3000', () => {
    Object.defineProperty(window, 'location', {
      value: {
        protocol: 'http:',
        hostname: 'admin.localhost',
        port: '3000',
        origin: 'http://admin.localhost:3000',
      },
      writable: true,
    });
    expect(getNormalSiteUrl()).toBe('http://localhost:3000');
  });

  it('maps admin.127.0.0.1 to localhost:3000', () => {
    Object.defineProperty(window, 'location', {
      value: {
        protocol: 'http:',
        hostname: 'admin.127.0.0.1',
        port: '3000',
        origin: 'http://admin.127.0.0.1:3000',
      },
      writable: true,
    });
    expect(getNormalSiteUrl()).toBe('http://localhost:3000');
  });
});

describe('getAdminSiteUrl', () => {
  const originalLocation = window.location;

  afterEach(() => {
    Object.defineProperty(window, 'location', {
      value: originalLocation,
      writable: true,
    });
    delete (process.env as any).REACT_APP_ADMIN_URL;
  });

  it('uses REACT_APP_ADMIN_URL when set', () => {
    (process.env as any).REACT_APP_ADMIN_URL = 'https://custom-admin.example.com';
    expect(getAdminSiteUrl()).toBe('https://custom-admin.example.com');
  });

  it('returns admin.localhost for localhost', () => {
    Object.defineProperty(window, 'location', {
      value: {
        protocol: 'http:',
        hostname: 'localhost',
        port: '3000',
        origin: 'http://localhost:3000',
      },
      writable: true,
    });
    expect(getAdminSiteUrl()).toBe('http://admin.localhost:3000');
  });

  it('returns admin.127.0.0.1 for 127.0.0.1', () => {
    Object.defineProperty(window, 'location', {
      value: {
        protocol: 'http:',
        hostname: '127.0.0.1',
        port: '3000',
        origin: 'http://127.0.0.1:3000',
      },
      writable: true,
    });
    expect(getAdminSiteUrl()).toBe('http://admin.localhost:3000');
  });

  it('prefixes hostname with admin', () => {
    Object.defineProperty(window, 'location', {
      value: {
        protocol: 'https:',
        hostname: 'example.com',
        port: '',
        origin: 'https://example.com',
      },
      writable: true,
    });
    expect(getAdminSiteUrl()).toBe('https://admin.example.com');
  });
});
