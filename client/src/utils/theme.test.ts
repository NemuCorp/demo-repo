import { Theme, applyTheme, getInitialTheme, getStoredTheme, getSystemTheme, THEME_STORAGE_KEY } from './theme';

describe('theme utils', () => {
  afterEach(() => {
    window.localStorage.clear();
    (window.matchMedia as any) = undefined;
    document.documentElement.removeAttribute('data-theme');
  });

  it('getInitialTheme returns the stored theme when present', () => {
    window.localStorage.setItem(THEME_STORAGE_KEY, 'dark');
    expect(getInitialTheme()).toBe('dark');
  });

  it('falls back to the matchMedia result when storage is empty', () => {
    window.matchMedia = jest.fn().mockReturnValue({ matches: true });
    expect(getInitialTheme()).toBe('dark');
  });

  it('getSystemTheme returns light when matchMedia is undefined', () => {
    (window.matchMedia as any) = undefined;
    expect(getSystemTheme()).toBe('light');
  });

  it('getStoredTheme rejects invalid strings by returning null', () => {
    window.localStorage.setItem(THEME_STORAGE_KEY, 'blue');
    expect(getStoredTheme()).toBeNull();
  });

  it('applyTheme sets the data-theme attribute on document.documentElement', () => {
    applyTheme('dark');
    expect(document.documentElement.getAttribute('data-theme')).toBe('dark');
    applyTheme('light');
    expect(document.documentElement.getAttribute('data-theme')).toBe('light');
  });

  it('getStoredTheme returns null when localStorage access throws', () => {
    const spy = jest.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new Error('denied');
    });
    expect(getStoredTheme()).toBeNull();
    spy.mockRestore();
  });
});