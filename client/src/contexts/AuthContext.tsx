import { createContext, useCallback, useContext, useEffect, useState } from 'react';
import * as api from '../services/api';
import { isAdminHost } from '../utils/host';
import { AuthResponse, User } from '../types';

interface AuthState {
  user: User | null;
  login: (email: string, password: string) => Promise<AuthResponse>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  isAuthenticated: boolean;
  isAdmin: boolean;
  isAdminView: boolean;
  loading: boolean;
}

const AuthContext = createContext<AuthState | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const doLogin = useCallback(async (email: string, password: string) => {
    const data = await api.login(email, password);
    localStorage.setItem('auth_token', data.token);
    setUser({
      id: data.session.user_id,
      email,
      is_admin: data.is_admin,
      created_at: data.session.created_at,
      updated_at: '',
    });
    return data;
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
      api.getCurrentUser()
        .then((data) => {
          setUser(data.user);
        })
        .catch(() => {
          localStorage.removeItem('auth_token');
          setUser(null);
        })
        .finally(() => {
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
        isAuthenticated: !!user,
        isAdmin: !!user?.is_admin,
        isAdminView: isAdminHost(),
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
