import { createContext, useState, useContext, useEffect, useCallback, useRef } from "react";
import { API_URL } from "../constants/url";
import type { User } from "../types/user";


interface AuthState {
  isAuthenticated: boolean,
  user: User | null;
  accessToken: string | null;
  isLoading: boolean;
  login: (form: FormData) => Promise<void>;
  register: (form: FormData) => Promise<void>;
  logout: () => Promise<void>;
  refresh: () => Promise<string | null>;
  authFetch: (url: string, options: RequestInit) => Promise<Response>;
}

export const AuthContext = createContext<AuthState | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const refreshing = useRef<Promise<string | null>>(null);

  const refresh = useCallback(async () => {
    if (refreshing.current) return refreshing.current;

    async function tryRefresh() {
      try {
        const res = await fetch(`${API_URL}/api/refresh`, {
          method: "POST",
          credentials: "include"
        });
        const data = await res.json();

        if (!res.ok)
          throw new Error(`Failed to refresh: ${data.error}`);
        setAccessToken(data.access_token);
        setUser(data.user);
        return data.access_token;
      } catch (err) {
        setAccessToken(null);
        setUser(null);
        return null;
      } finally {
        refreshing.current = null;
      }
    }
    refreshing.current = tryRefresh();
    return refreshing.current;
  }, []);

  const authFetch = useCallback(async (url: string, options: RequestInit = {}) => {
    const headers = new Headers(options.headers);
    if (accessToken) {
      headers.set("Authorization", `Bearer ${accessToken}`);
    }
    let res = await fetch(url, { ...options, headers });
    if (res.status === 401) {
      const newToken = await refresh();
      if (newToken) {
        headers.set("Authorization", `Bearer ${newToken}`);
        res = await fetch(url, { ...options, headers });
      }
    }
    return res;
  }, [accessToken, refresh]);

  async function login(form: FormData) {
    const res = await fetch(`${API_URL}/api/login`, {
      method: "POST",
      body: form,
      credentials: "include"
    });
    const data = await res.json();

    if (!res.ok)
      throw new Error(`Failed to login ${data.error}`);

    setAccessToken(data.access_token);
    setUser(data.user);
  }

  async function register(form: FormData) {
    const res = await fetch(`${API_URL}/api/register`, {
      method: "POST",
      body: form,
    });

    if (!res.ok) {
      const data = await res.json();
      throw new Error(`Failed to register ${data.error}`);
    }

  }

  async function logout() {
    await fetch(`${API_URL}/api/logout`, { method: "POST" });
    setUser(null);
    setAccessToken(null);
  }


  useEffect(() => {
    const initAuth = async () => {
      try {
        await refresh();
      } catch (err) {
        console.error("Initial session check failed", err);
      } finally {
        setIsLoading(false);
      }
    };
    initAuth();
  }, [refresh]);

  return (
    <AuthContext.Provider value={{
      isAuthenticated: accessToken !== null,
      user,
      isLoading,
      accessToken,
      login,
      register,
      logout,
      refresh,
      authFetch
    }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthState {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used inside <AuthProvider/>");
  return context;
}
