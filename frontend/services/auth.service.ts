import { api, setAccessToken } from "@/lib/api";
import type { AuthResponse, User } from "@/types";

export const authService = {
  async login(email: string, password: string): Promise<AuthResponse> {
    const { data } = await api.post<AuthResponse>("/auth/login", { email, password });
    setAccessToken(data.accessToken);
    if (typeof window !== "undefined") {
      localStorage.setItem("ss_refresh", data.refreshToken);
    }
    return data;
  },
  async refresh(refreshToken: string): Promise<AuthResponse> {
    const { data } = await api.post<AuthResponse>("/auth/refresh", { refreshToken });
    setAccessToken(data.accessToken);
    if (typeof window !== "undefined") {
      localStorage.setItem("ss_refresh", data.refreshToken);
    }
    return data;
  },
  async logout(): Promise<void> {
    const rt = typeof window !== "undefined" ? localStorage.getItem("ss_refresh") : null;
    try {
      await api.post("/auth/logout", { refreshToken: rt });
    } catch {
      /* noop */
    }
    setAccessToken(null);
    if (typeof window !== "undefined") {
      localStorage.removeItem("ss_refresh");
    }
  },
  async me(): Promise<User> {
    const { data } = await api.get<User>("/auth/me");
    return data;
  },
};
