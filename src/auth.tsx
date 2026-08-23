import { createContext, useContext, useMemo, useState } from 'react';
import type { ReactNode } from 'react';

const STORAGE_KEY = 'callcalendar.admin.auth';

interface AdminAuthValue {
  isAuthenticated: boolean;
  login: () => void;
  logout: () => void;
}

const AdminAuthContext = createContext<AdminAuthValue | null>(null);

function readAuth(): boolean {
  return localStorage.getItem(STORAGE_KEY) === '1';
}

export function AdminAuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(readAuth);

  const value = useMemo<AdminAuthValue>(
    () => ({
      isAuthenticated,
      login: () => {
        localStorage.setItem(STORAGE_KEY, '1');
        setIsAuthenticated(true);
      },
      logout: () => {
        localStorage.removeItem(STORAGE_KEY);
        setIsAuthenticated(false);
      },
    }),
    [isAuthenticated],
  );

  return <AdminAuthContext.Provider value={value}>{children}</AdminAuthContext.Provider>;
}

export function useAdminAuth(): AdminAuthValue {
  const ctx = useContext(AdminAuthContext);
  if (!ctx) {
    throw new Error('useAdminAuth must be used within AdminAuthProvider');
  }
  return ctx;
}
