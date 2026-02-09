// src/app/components/auth/require-auth.tsx
import * as React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { getToken } from "./http"; 

export function RequireAuth({ children }: { children: React.ReactNode }) {
  const location = useLocation();
  const token = getToken();

  if (!token) {
    return (
      <Navigate
        to="/login"
        replace
        state={{ from: location.pathname }}
      />
    );
  }

  return <>{children}</>;
}
