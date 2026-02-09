// src/app/components/auth/register-page.tsx
import * as React from "react";
import { useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { apiFetch } from "../api/http";

export function RegisterPage() {
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const canSubmit = useMemo(
    () => email.trim().length > 0 && password.trim().length > 0 && !busy,
    [email, password, busy]
  );

  async function onRegister(e?: React.FormEvent) {
    e?.preventDefault();
    if (!canSubmit) return;

    setBusy(true);
    setError(null);

    try {
      const res = await apiFetch<{ id: number }>("/auth/register", {
        method: "POST",
        body: JSON.stringify({ email: email.trim(), password }),
      });

      // userflow: register -> login
      console.log("Registered user id:", res.id);
      navigate("/login", { replace: true });
    } catch (err: any) {
      setError(err?.message ?? "Registration failed");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-6">
      <div className="w-full max-w-md bg-card border border-border rounded-2xl p-6 shadow-sm">
        <h1 className="text-2xl text-foreground mb-1">Register</h1>
        <p className="text-sm text-muted-foreground">
          Create an account to start sessions.
        </p>

        <form className="mt-6 space-y-3" onSubmit={onRegister}>
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
              autoComplete="new-password"
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
            {busy ? "Creating..." : "Create account"}
          </button>

          <div className="text-sm text-muted-foreground text-center">
            Already have an account?{" "}
            <Link className="text-primary hover:underline" to="/login">
              Login
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
