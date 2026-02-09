// src/app/components/auth/login-page.tsx
import * as React from "react";
import { useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { apiFetch, setToken } from "../api/http";

export function LoginPage() {
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const canSubmit = useMemo(
    () => email.trim().length > 0 && password.trim().length > 0 && !busy,
    [email, password, busy]
  );

  async function onLogin(e?: React.FormEvent) {
    e?.preventDefault();
    if (!canSubmit) return;

    setBusy(true);
    setError(null);

    try {
      const res = await apiFetch<{ token: string }>("/auth/login", {
        method: "POST",
        body: JSON.stringify({ email: email.trim(), password }),
      });

      setToken(res.token);

      // ✅ твой редирект
      navigate("/dashboard", { replace: true });
    } catch (err: any) {
      setError(err?.message ?? "Login failed");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-6">
      <div className="w-full max-w-md bg-card border border-border rounded-2xl p-6 shadow-sm">
        <h1 className="text-2xl text-foreground mb-1">Login</h1>
        <p className="text-sm text-muted-foreground">
          Sign in to continue to Pythia.
        </p>

        <form className="mt-6 space-y-3" onSubmit={onLogin}>
          <div className="space-y-1">
            <label className="text-sm text-foreground">Email</label>
            <input
              className="w-full rounded-xl border border-border bg-input-background px-3 py-2 outline-none focus:ring-2 focus:ring-primary/30"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              autoComplete="email"
              placeholder="you@example.com"
            />
          </div>

          <div className="space-y-1">
            <label className="text-sm text-foreground">Password</label>
            <input
              className="w-full rounded-xl border border-border bg-input-background px-3 py-2 outline-none focus:ring-2 focus:ring-primary/30"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="current-password"
              placeholder="••••••••"
              type="password"
            />
          </div>

          {error && (
            <div className="rounded-xl border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
              {error}
            </div>
          )}

          <button
            type="submit"
            disabled={!canSubmit}
            className="w-full rounded-xl bg-stone-900 text-white py-2 disabled:opacity-50"
          >
            {busy ? "Signing in..." : "Login"}
          </button>

          <div className="text-sm text-muted-foreground text-center">
            No account?{" "}
            <Link className="text-primary hover:underline" to="/register">
              Register
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
