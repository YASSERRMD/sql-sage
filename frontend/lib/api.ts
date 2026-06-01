import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";

const baseURL = process.env.NEXT_PUBLIC_API_URL || "";

export const api = axios.create({
  baseURL: `${baseURL}/api/v1`,
  timeout: 120_000,
  headers: { "Content-Type": "application/json" },
});

let accessToken: string | null = null;
let refreshing: Promise<string | null> | null = null;

export function setAccessToken(token: string | null) {
  accessToken = token;
  if (typeof window !== "undefined") {
    if (token) localStorage.setItem("ss_token", token);
    else localStorage.removeItem("ss_token");
  }
}

export function loadAccessToken() {
  if (typeof window === "undefined") return null;
  if (accessToken) return accessToken;
  const stored = localStorage.getItem("ss_token");
  accessToken = stored;
  return stored;
}

api.interceptors.request.use((cfg: InternalAxiosRequestConfig) => {
  const t = loadAccessToken();
  if (t && cfg.headers) {
    cfg.headers.Authorization = `Bearer ${t}`;
  }
  return cfg;
});

api.interceptors.response.use(
  (r) => r,
  async (error: AxiosError) => {
    const original = error.config as InternalAxiosRequestConfig & { _retry?: boolean };
    if (error.response?.status === 401 && !original._retry) {
      original._retry = true;
      const rt = typeof window !== "undefined" ? localStorage.getItem("ss_refresh") : null;
      if (rt) {
        try {
          refreshing =
            refreshing ??
            axios
              .post(`${baseURL}/api/v1/auth/refresh`, { refreshToken: rt })
              .then((res) => {
                const data = res.data as { accessToken: string; refreshToken: string };
                setAccessToken(data.accessToken);
                localStorage.setItem("ss_refresh", data.refreshToken);
                return data.accessToken;
              })
              .catch(() => {
                setAccessToken(null);
                localStorage.removeItem("ss_refresh");
                return null;
              })
              .finally(() => {
                refreshing = null;
              });
          const t = await refreshing;
          if (t) {
            original.headers = original.headers ?? ({} as any);
            original.headers.Authorization = `Bearer ${t}`;
            return api(original);
          }
        } catch {
          /* noop */
        }
      }
    }
    return Promise.reject(error);
  }
);
