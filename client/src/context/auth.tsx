import { createContext, useState, useContext, useEffect } from "react";
import { API_URL } from "../constants/url";

export type UserAuth = {
  user: {
    id: number;
    username: string;
    email: string;
  };
  accessToken: string;
};

interface AuthState {
  login: (userAuth: UserAuth) => Promise<boolean>;
  logout: () => Promise<void>;
  refresh: () => Promise<string>;
  isLoading: boolean;
  userAuth: UserAuth|null;
}

export const AuthContext = createContext<AuthState|undefined>(undefined);

export function AuthContextProvider({ children }: { children: React.ReactNode }) {
  const [userAuth, setUserAuth] = useState<UserAuth|null>(null);
  const [isLoading, setIsLoading] = useState(true);
  
  
  async function login(userAuth: UserAuth) {
    setUserAuth(userAuth);
    return true;
  }

  async function logout() {
    await fetch(`${API_URL}/api/logout`, { method: "POST" });
    setUserAuth(null);
  }

  async function refresh() {
    const res = await fetch(`${API_URL}/api/refresh`, {
      method: "POST",
      credentials: "include",
    });
    const data = await res.json();
    if (!res.ok) {
      await logout();
      return Promise.reject(new Error("Failed to refresh token: " + data.error));
    }
    
    setUserAuth(prev => {
      if (!prev) return null;
      return { ...prev, accessToken: data.access_token };
    });
    return data.access_token;
  }

  useEffect(() => {
    (async() => {
      try {
        await refresh();
      } catch (err) {
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);

  return (
    <AuthContext.Provider value={{
        login,
        logout,
        refresh,
        isLoading,
        userAuth,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used inside <AuthContextProvider/>");
  return context;
}
