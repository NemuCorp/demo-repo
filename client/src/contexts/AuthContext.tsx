import { createContext, useCallback, useContext, useEffect, useState } from 'react';
import * as api from '../services/api';
import { User } from '../types';

interface AuthState {
  user: User | null;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  isAuthenticated: boolean;
  loading: boolean;
}

const AuthContext = createContext<AuthState | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const doLogin = useCallback(async (email: string, password: string) => {
    const data = await api.login(email, password);
    localStorage.setItem('auth_token', data.token);
    setUser({ id: data.session.user_id, email, created_at: data.session.created_at, updated_at: '' });
  }, []);

  const doRegister = useCallback(async (email: string, password: string) => {
    const data = await api.register(email, password);
    setUser(data.user);
  }, []);

  const doLogout = useCallback(async () => {
    try {
      await api.logout();
    } catch {
      // ignore errors on logout
    }
    localStorage.removeItem('auth_token');
    setUser(null);
  }, []);

  useEffect(() => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      api.getCart().then(() => {
        setUser({ id: 0, email: '', created_at: '', updated_at: '' });
      }).catch(() => {
        localStorage.removeItem('auth_token');
      }).finally(() => {
        setLoading(false);
      });
    } else {
      setLoading(false);
    }
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        login: doLogin,
        register: doRegister,
        logout: doLogout,
        isAuthenticated: !!localStorage.getItem('auth_token'),
        loading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return ctx;
}
