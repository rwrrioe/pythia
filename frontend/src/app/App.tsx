// src/app/App.tsx
import * as React from "react";
import { useEffect, useMemo } from "react";
import { Navigate, Route, Routes, useLocation, useNavigate, useParams } from "react-router-dom";

import { AppLayout } from "./components/app-layout";

import { RequireAuth } from "./components/api/require-auth";
import { LoginPage } from "./components/auth/login-page";
import { RegisterPage } from "./components/auth/register-page";

import { Dashboard } from "./components/dashboard";
import { SessionsPage } from "./components/sessions-page";
import { SessionSetup } from "./components/session-setup";
import { SessionCapture } from "./components/session-capture";
import { FinalWordSelection } from "./components/final-word-selection";
import { FlashcardsReview } from "./components/flashcards-review";
import { TestsMode } from "./components/tests-mode";
import { SessionSummary } from "./components/session-summary";
import { SessionDetail } from "./components/session-detail";
import { Profile } from "./components/profile";
import { Stats } from "./components/stats";

import type { SessionRecord, SessionWord } from "./components/session-types";
import type { SessionConfig } from "./components/session-setup";

import { createSession } from "./components/api/sessions";
import { apiFetch } from "./components/api/http";

function PrivateShell({ children }: { children: React.ReactNode }) {
  const location = useLocation();
  const navigate = useNavigate();

  const currentScreen = useMemo(() => {
    const p = location.pathname;
    if (p.startsWith("/dashboard")) return "dashboard";
    if (p.startsWith("/sessions")) return "sessions";
    if (p.startsWith("/stats")) return "stats";
    if (p.startsWith("/profile")) return "profile";
    return "dashboard";
  }, [location.pathname]);

  return (
    <AppLayout
      currentScreen={currentScreen as any}
      onNavigate={(s) => {
        if (s === "dashboard") navigate("/dashboard");
        else if (s === "sessions") navigate("/sessions");
        else if (s === "stats") navigate("/stats");
        else if (s === "profile") navigate("/profile");
        else navigate("/dashboard");
      }}
    >
      {children}
    </AppLayout>
  );
}

/** language mapping: 1 en, 2 de, 3 fr, 4 es */
function mapLangIdToCode(lang: any): string {
  const n = Number(lang);
  if (n === 1) return "en";
  if (n === 2) return "de";
  if (n === 3) return "fr";
  if (n === 4) return "es";
  return String(lang ?? "");
}

function mapLibrarySessionToRecord(s: any): SessionRecord {
  const id = String(s?.Id ?? s?.id ?? "");
  const startedAt = s?.StartedAt ?? s?.started_at ?? s?.startedAt;
  const endedAt = s?.ended_at ?? s?.EndedAt ?? s?.endedAt;

  const name = String(s?.name ?? s?.Name ?? "").trim();
  const title = name || `Session ${id}`;

  return {
    id,
    name,
    title,
    status: String(s?.Status ?? s?.status ?? ""),
    language: mapLangIdToCode(s?.language ?? s?.Language),
    level: s?.level ?? s?.Level ?? 0,
    accuracy: Number(s?.Accuracy ?? s?.accuracy ?? 0),
    started_at: startedAt,
    ended_at: endedAt,
    createdAt: startedAt,
    lastStudiedAt: endedAt ?? startedAt,
    duration: Number(s?.Duration ?? s?.duration ?? 0),
    wordsCount: Number(s?.wordsCount ?? 0),
    wordLimit: Number(s?.wordLimit ?? 0),
  } as any;
}

function mapFlashcardsToWords(flashcards: any[]): SessionWord[] {
  const arr = Array.isArray(flashcards) ? flashcards : [];
  return arr.map((f: any, idx: number) => ({
    id: String(f?.id ?? `${idx}`),
    word: String(f?.word ?? ""),
    translation: String(f?.translation ?? ""),
    language: String(f?.language ?? ""),
    known: !!f?.known,
  })) as any;
}

/**
 * ✅ ВАЖНО: эти компоненты ВНЕ SessionFlowRoutes,
 * иначе при каждом setState SessionFlowRoutes они будут remount,
 * и useEffect будет снова стрелять запросами (лавина).
 */
function SessionsIndexRoute({
  sessions,
  onOpenSessionId,
}: {
  sessions: SessionRecord[];
  onOpenSessionId: (id: string) => void;
}) {
  return (
    <PrivateShell>
      <SessionsPage
        sessions={sessions as any}
        onOpenSession={(id: string) => onOpenSessionId(String(id))}
      />
    </PrivateShell>
  );
}

function SessionDetailRouteStable({
  sessions,
  activeSession,
  loadLibrarySessionDetail,
}: {
  sessions: SessionRecord[];
  activeSession: SessionRecord | null;
  loadLibrarySessionDetail: (id: string) => Promise<SessionRecord>;
}) {
  const navigate = useNavigate();
  const { sessionId } = useParams();
  const id = String(sessionId ?? "");

  const localFromList = sessions.find((s) => String((s as any).id) === id) ?? null;
  const local = localFromList ?? activeSession;

  const [resolved, setResolved] = React.useState<SessionRecord | null>(local);
  const [loading, setLoading] = React.useState<boolean>(false);
  const [err, setErr] = React.useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    void (async () => {
      if (!id) return;

      const hasWords = !!(local as any)?.words?.length;
      if (local && hasWords) {
        setResolved(local);
        return;
      }

      try {
        setLoading(true);
        setErr(null);
        const s = await loadLibrarySessionDetail(id);
        if (!mounted) return;
        setResolved(s);
      } catch (e: any) {
        console.error(e);
        if (!mounted) return;
        setErr(e?.message ?? "Failed to load session");
        setResolved(local ?? null);
      } finally {
        if (mounted) setLoading(false);
      }
    })();

    return () => {
      mounted = false;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  if (loading && !resolved) {
    return (
      <PrivateShell>
        <div className="min-h-screen bg-background p-8">
          <div className="max-w-3xl mx-auto">
            <p className="text-muted-foreground">Loading session…</p>
          </div>
        </div>
      </PrivateShell>
    );
  }

  if (!resolved) {
    return (
      <PrivateShell>
        <div className="min-h-screen bg-background p-8">
          <div className="max-w-3xl mx-auto">
            <p className="text-muted-foreground">{err ?? "Session not found."}</p>
            <button
              onClick={() => navigate("/sessions")}
              className="mt-4 px-4 py-2 bg-card border border-border rounded-lg hover:bg-muted transition-colors"
            >
              Back to Sessions
            </button>
          </div>
        </div>
      </PrivateShell>
    );
  }

  return (
    <PrivateShell>
      <SessionDetail
        session={resolved as any}
        onBack={() => navigate("/sessions")}
        onContinueFlashcards={() => navigate("/session/flashcards")}
        onStartTest={() => navigate("/session/tests")}
        onReviewWords={() => {}}
      />
    </PrivateShell>
  );
}

function SessionFlowRoutes() {
  const navigate = useNavigate();

  const [sessions, setSessions] = React.useState<SessionRecord[]>([]);
  const [activeSession, setActiveSession] = React.useState<SessionRecord | null>(null);

  const [sessionConfig, setSessionConfig] = React.useState<SessionConfig | null>(null);

  const [capturedWords, setCapturedWords] = React.useState<SessionWord[]>([]);
  const [finalWords, setFinalWords] = React.useState<SessionWord[]>([]);
  const [sessionScore, setSessionScore] = React.useState(0);

  const [captureSessionId, setCaptureSessionId] = React.useState<number | null>(null);

  const startNewSession = () => {
    setSessionConfig(null);
    setCapturedWords([]);
    setFinalWords([]);
    setSessionScore(0);
    setCaptureSessionId(null);
    setActiveSession(null);
    navigate("/session/new");
  };

  const handleSessionSetupContinue = async (config: SessionConfig) => {
    setSessionConfig(config);

    const res = await createSession({
      language: config.language,
      wordLimit: config.wordLimit,
      durationSeconds: 30 * 60,
    });

    setCaptureSessionId(res.session_id);
    navigate("/session/capture");
  };

  const handleCaptureFinish = (_words: SessionWord[], _title: string, sessionId: number) => {
    setCaptureSessionId(sessionId);
    navigate("/session/flashcards");
  };

  const handleStartTests = () => navigate("/session/tests");

  const handleTestsComplete = (score: number) => {
    setSessionScore(score);
    navigate("/session/summary");
  };

  const backToDashboard = () => navigate("/dashboard");

  const handleFinalizeConfirm = (words: SessionWord[]) => {
    setFinalWords(words);
    navigate("/session/flashcards");
  };

  async function loadLibrarySessionDetail(sessionId: string) {
    const res = await apiFetch<{ session?: any; flashcards?: any[] }>(`/library/session/${sessionId}`, {
      method: "GET",
    });

    const base = mapLibrarySessionToRecord(res?.session ?? { Id: sessionId });
    const words = mapFlashcardsToWords(res?.flashcards ?? []);
    const merged = { ...(base as any), words, wordsCount: words.length } as any;

    setActiveSession(merged);
    setSessions((prev) => {
      const idx = prev.findIndex((s) => String((s as any).id) === String(sessionId));
      if (idx === -1) return prev;
      const next = prev.slice();
      next[idx] = { ...(next[idx] as any), ...(merged as any) };
      return next;
    });

    return merged as SessionRecord;
  }

  return (
    <Routes>
      <Route
        path="/dashboard"
        element={
          <PrivateShell>
            <Dashboard />
          </PrivateShell>
        }
      />

      {/* ✅ /sessions теперь НЕ объявлен как inner component -> нет remount loop -> нет лавины */}
      <Route
        path="/sessions"
        element={
          <SessionsIndexRoute
            sessions={sessions}
            onOpenSessionId={(id) => {
              void (async () => {
                try {
                  await loadLibrarySessionDetail(String(id));
                } catch (e) {
                  console.error(e);
                } finally {
                  navigate(`/sessions/${id}`);
                }
              })();
            }}
          />
        }
      />

      <Route
        path="/sessions/:sessionId"
        element={
          <SessionDetailRouteStable
            sessions={sessions}
            activeSession={activeSession}
            loadLibrarySessionDetail={loadLibrarySessionDetail}
          />
        }
      />

      <Route
        path="/stats"
        element={
          <PrivateShell>
            <Stats onBack={backToDashboard} />
          </PrivateShell>
        }
      />

      <Route
        path="/profile"
        element={
          <PrivateShell>
            <Profile onBack={backToDashboard} />
          </PrivateShell>
        }
      />

      <Route path="/session/new" element={<SessionSetup onContinue={handleSessionSetupContinue} onBack={backToDashboard} />} />

      <Route
        path="/session/capture"
        element={
          sessionConfig && captureSessionId !== null ? (
            <SessionCapture
              sessionId={captureSessionId}
              config={sessionConfig as any}
              durationSeconds={30 * 60}
              onFinish={handleCaptureFinish}
              onCancel={backToDashboard}
            />
          ) : (
            <Navigate to="/session/new" replace />
          )
        }
      />

      <Route
        path="/session/finalize"
        element={
          <FinalWordSelection
            words={capturedWords}
            wordLimit={sessionConfig?.wordLimit ?? 15}
            onConfirm={handleFinalizeConfirm}
            onBack={() => navigate("/session/capture")}
          />
        }
      />

      <Route path="/session/flashcards" element={<FlashcardsReview sessionId={captureSessionId} onDone={handleStartTests} />} />

      <Route
        path="/session/tests"
        element={
          <TestsMode
            sessionId={captureSessionId}
            onBack={() => navigate("/session/flashcards")}
            onComplete={handleTestsComplete}
          />
        }
      />

      <Route
        path="/session/summary"
        element={<SessionSummary words={finalWords as any} score={sessionScore} onBackToDashboard={backToDashboard} />}
      />

      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  );
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/login" replace />} />

      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />

      <Route
        path="/*"
        element={
          <RequireAuth>
            <SessionFlowRoutes />
          </RequireAuth>
        }
      />
    </Routes>
  );
}
