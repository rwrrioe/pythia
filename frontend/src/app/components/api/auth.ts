import { apiFetch, setToken } from "./http";
import { routes } from "./routes";

/* ---------- Auth ---------- */

export async function login(email: string, password: string) {
  const res = await apiFetch<{ token: string }>(routes.login, {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });

  setToken(res.token);
  return res;
}

export async function register(email: string, password: string) {
  return apiFetch<{ id: number }>(routes.register, {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

/* ---------- Auth state ---------- */

// Источник истины — наличие JWT
export function isAuthenticated(): boolean {
  return !!localStorage.getItem("token");
}
