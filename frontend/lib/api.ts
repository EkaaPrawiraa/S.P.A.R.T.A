/**
 * Centralized API client with Bearer token injection and envelope unwrapping
 */

interface ApiResponse<T> {
  status: "success" | "error";
  data?: T;
  message?: string;
}

function normalizeToken(input: string): string {
  let token = input.trim();
  if (token.length === 0) return "";

  // Strip accidental surrounding quotes (common when values are copied/pasted or URL-decoded).
  if (
    (token.startsWith('"') && token.endsWith('"')) ||
    (token.startsWith("'") && token.endsWith("'"))
  ) {
    token = token.slice(1, -1).trim();
  }

  // Support raw token or "Bearer <token>" or accidental "Bearer Bearer <token>".
  const fields = token.split(/\s+/).filter(Boolean);
  if (fields.length === 1) return fields[0];
  if (fields.length === 2 && /^bearer$/i.test(fields[0])) return fields[1];
  if (
    fields.length === 3 &&
    /^bearer$/i.test(fields[0]) &&
    /^bearer$/i.test(fields[1])
  ) {
    return fields[2];
  }

  // Fallback: return trimmed original.
  return token;
}

class ApiClient {
  private baseUrl: string;

  constructor(
    baseUrl: string = process.env.NEXT_PUBLIC_API_BASE_URL ||
      "http://localhost:8080",
  ) {
    this.baseUrl = baseUrl;
  }

  private getToken(): string | null {
    if (typeof window === "undefined") return null;

    const fromStorage = localStorage.getItem("auth_token");
    if (fromStorage && fromStorage.trim().length > 0) {
      const normalized = normalizeToken(fromStorage);
      return normalized.length > 0 ? normalized : null;
    }

    // Fallback for cases where middleware auth cookie exists but localStorage was cleared.
    const cookie = document.cookie
      .split(";")
      .map((c) => c.trim())
      .find((c) => c.startsWith("auth_token="));
    if (!cookie) return null;

    const raw = cookie.slice("auth_token=".length);
    try {
      const decoded = decodeURIComponent(raw);
      const normalized = normalizeToken(decoded);
      return normalized.length > 0 ? normalized : null;
    } catch {
      const normalized = normalizeToken(raw);
      return normalized.length > 0 ? normalized : null;
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const token = this.getToken();

    const headers = new Headers(options.headers);
    headers.set("Content-Type", "application/json");
    if (token) headers.set("Authorization", `Bearer ${token}`);

    const response = await fetch(url, {
      ...options,
      headers,
    });

    let json: ApiResponse<T>;
    try {
      json = (await response.json()) as ApiResponse<T>;
    } catch {
      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }
      throw new Error("Invalid API response");
    }

    if (!response.ok) {
      if (response.status === 401 && typeof window !== "undefined") {
        // Cookie-only middleware can allow /app while the JWT is expired/invalid.
        // If API returns 401, force re-auth to avoid the UI being stuck in error state.
        const { clearAuth } = await import("@/lib/auth");
        clearAuth();
        // Avoid infinite loops if login itself fails.
        if (!window.location.pathname.startsWith("/login")) {
          window.location.assign("/login");
        }
      }
      throw new Error(json.message || `API error: ${response.status}`);
    }

    if (json.status === "error") {
      throw new Error(json.message || "API returned error");
    }

    return json.data as T;
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "GET" });
  }

  async post<T>(endpoint: string, body?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async put<T>(endpoint: string, body?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "DELETE" });
  }
}

export const api = new ApiClient();
