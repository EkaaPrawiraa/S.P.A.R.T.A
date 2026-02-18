/**
 * Auth state management helpers
 */

interface AuthState {
  token: string;
  userId: string;
}

const STORAGE_KEY = "auth_token";
const USER_ID_KEY = "auth_user_id";

function setCookie(name: string, value: string) {
  const secure =
    typeof window !== "undefined" && window.location.protocol === "https:";
  document.cookie = `${name}=${encodeURIComponent(value)}; Path=/; SameSite=Lax${secure ? "; Secure" : ""}`;
}

function deleteCookie(name: string) {
  document.cookie = `${name}=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax`;
}

export function setAuthToken(token: string, userId: string) {
  if (typeof window !== "undefined") {
    localStorage.setItem(STORAGE_KEY, token);
    localStorage.setItem(USER_ID_KEY, userId);
    setCookie(STORAGE_KEY, token);
  }
}

export function getAuthToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(STORAGE_KEY);
}

export function getUserId(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(USER_ID_KEY);
}

export function getAuthState(): AuthState | null {
  const token = getAuthToken();
  const userId = getUserId();

  if (!token || !userId) return null;

  return { token, userId };
}

export function clearAuth() {
  if (typeof window !== "undefined") {
    localStorage.removeItem(STORAGE_KEY);
    localStorage.removeItem(USER_ID_KEY);
    deleteCookie(STORAGE_KEY);
  }
}

export function isAuthenticated(): boolean {
  return getAuthState() !== null;
}
