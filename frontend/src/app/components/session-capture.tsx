// src/app/components/session-capture.tsx
import * as React from "react";
import { useEffect, useMemo, useRef, useState } from "react";
import { motion } from "motion/react";
import { ArrowLeft, Upload, FileText, Clock } from "lucide-react";
import type { SessionWord } from "./session-types";
import { GreekPattern } from "./greek-pattern";
import { PifMascot } from "./pif-mascot";

import {
  uploadFile,
  translateTask,
  endSessionAndFetchFlashcards,
  langIdToCode,
  type LangId,
  type LangCode,
} from "./api/sessions";
import { openSessionWs, type SessionWsMessage } from "./api/ws";

export interface SessionConfig {
  language: LangId; // 1..4 строго
  difficulty: string;
  wordLimit: number;
}

interface SessionCaptureProps {
  config: SessionConfig;
  durationSeconds: number;
  sessionId: number;
  onFinish: (words: SessionWord[], title: string, sessionId: number) => void;
  onCancel: () => void;
}

function formatMMSS(totalSeconds: number) {
  const s = Math.max(0, totalSeconds);
  const mm = Math.floor(s / 60);
  const ss = s % 60;
  return `${String(mm).padStart(2, "0")}:${String(ss).padStart(2, "0")}`;
}

function normalizeWord(s: string) {
  return String(s ?? "").trim().toLowerCase();
}

function mergeUnique(prev: SessionWord[], incoming: SessionWord[]) {
  const seen = new Set<string>();
  const out: SessionWord[] = [];
  for (const w of [...incoming, ...prev]) {
    const key = normalizeWord((w as any)?.word ?? (w as any)?.Word ?? "");
    if (!key) continue;
    if (seen.has(key)) continue;
    seen.add(key);
    out.push({
      ...w,
      word: (w as any).word ?? (w as any).Word ?? "",
      translation: (w as any).translation ?? (w as any).Translation ?? "",
    } as any);
  }
  return out;
}

function mapLevel(difficulty: string): string {
  const d = String(difficulty || "").trim().toUpperCase();
  if (d.includes("A1")) return "A1";
  if (d.includes("A2")) return "A2";
  if (d.includes("B1")) return "B1";
  if (d.includes("B2")) return "B2";
  if (d.includes("C1")) return "C1";
  return "B1";
}

export function SessionCapture({
  config,
  durationSeconds,
  sessionId,
  onFinish,
  onCancel,
}: SessionCaptureProps) {
  const fileRef = useRef<HTMLInputElement | null>(null);

  const [busy, setBusy] = useState(false);

  // UI status
  const [stage, setStage] = useState<string | null>(null);
  const [status, setStatus] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const [lastTaskId, setLastTaskId] = useState<string | null>(null);
  const translatedTasksRef = useRef<Set<string>>(new Set());
  const ocrDoneRef = useRef<Set<string>>(new Set());

  const wsRef = useRef<WebSocket | null>(null);

  const [startedAt] = useState<number>(() => Date.now());
  const [remaining, setRemaining] = useState<number>(durationSeconds);

  const [allWords, setAllWords] = useState<SessionWord[]>([]);

  const langCode: LangCode = useMemo(() => langIdToCode(config.language), [config.language]);
  const headerLang = langCode.toUpperCase();
  const headerLevel = String(config.difficulty || "").toUpperCase();

  // WS connect
  useEffect(() => {
    try {
      const ws = openSessionWs(sessionId, (msg: SessionWsMessage) => {
        console.log("[WS]", msg);

        if (msg?.error) {
          setError(String(msg.error));
          return;
        }

        if (msg?.stage) setStage(String(msg.stage));
        if (msg?.status) setStatus(String(msg.status));
        if (msg?.task_id) setLastTaskId(String(msg.task_id));

        if (msg?.stage === "ocr" && msg?.status === "done") {
          const taskId = String(msg?.task_id ?? "");
          if (!taskId) return;

          if (ocrDoneRef.current.has(taskId)) return;
          ocrDoneRef.current.add(taskId);

          void (async () => {
            try {
              setError(null);
              await translateTask(sessionId, taskId, {
                duration: "6 months",
                lang: langCode,
                level: mapLevel(config.difficulty),
              });
            } catch (e: any) {
              console.error("translateTask REST failed", e);
              setError(e?.message ?? "translateTask failed");
            }
          })();

          return;
        }

        // Translate done -> words arrive via WS
        if (msg?.stage === "translate" && msg?.status === "done") {
          const taskId = String(msg?.task_id ?? "");
          if (taskId) translatedTasksRef.current.add(taskId);

          if (Array.isArray(msg?.words)) {
            setAllWords((prev) => mergeUnique(prev, msg.words as any));
          }
          return;
        }

        if (msg?.type === "words" && Array.isArray((msg as any).words)) {
          setAllWords((prev) => mergeUnique(prev, (msg as any).words as any));
        }
      });

      wsRef.current = ws;
    } catch (e) {
      console.warn("WS not available", e);
      setError("WebSocket not available");
    }

    return () => {
      wsRef.current?.close();
      wsRef.current = null;
    };
  }, [sessionId, langCode, config.difficulty]);

  // Timer (auto end)
  useEffect(() => {
    const id = window.setInterval(() => {
      const elapsed = Math.floor((Date.now() - startedAt) / 1000);
      const left = Math.max(0, durationSeconds - elapsed);
      setRemaining(left);
      if (left <= 0) {
        window.clearInterval(id);
        void handleEndSession(true);
      }
    }, 250);

    return () => window.clearInterval(id);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const progressPct = useMemo(() => {
    const used = durationSeconds - remaining;
    const pct = (used / Math.max(1, durationSeconds)) * 100;
    return Math.max(0, Math.min(100, pct));
  }, [durationSeconds, remaining]);

  const handleUploadClick = () => fileRef.current?.click();

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0];
    e.target.value = "";
    if (!f) return;

    setBusy(true);
    setError(null);

    setStage("ocr");
    setStatus("processing");

    try {
      const res = await uploadFile(sessionId, f, langCode);
      if (res?.task_id) setLastTaskId(String(res.task_id));
    } catch (e: any) {
      console.error("uploadFile failed", e);
      setError(e?.message ?? "uploadFile failed");
      setStatus(null);
    } finally {
      setBusy(false);
    }
  };

  const handleClear = () => {
    setAllWords([]);
    setLastTaskId(null);
    setStage(null);
    setStatus(null);
    setError(null);
    translatedTasksRef.current = new Set();
    ocrDoneRef.current = new Set();
  };

  const handleEndSession = async (auto = false) => {
    if (busy) return; // hard anti-double-click
    setBusy(true);
    setError(null);

    try {
      // Requirement #2: PATCH /end then GET /learn/flashcards
      const cards = await endSessionAndFetchFlashcards(sessionId);

      const title =
        `${langCode} • ${new Date().toISOString().slice(0, 10)}` +
        (auto ? " (auto)" : "");

      onFinish(cards.words, title, sessionId);
    } catch (err: any) {
      console.error("endSession/getFlashcards failed", err);
      setError(err?.message ?? "Failed to end session");
    } finally {
      setBusy(false);
    }
  };

  const hasWords = allWords.length > 0;

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background p-4 md:p-8 relative overflow-hidden">
      <div className="absolute top-0 left-0 right-0">
        <GreekPattern className="w-full h-6 text-primary opacity-40" />
      </div>

      <div className="max-w-5xl mx-auto relative z-10 pt-6 md:pt-8">
        <div className="flex items-start justify-between gap-3 mb-4 md:mb-6">
          <div className="flex items-start gap-2 md:gap-3">
            <button
              onClick={onCancel}
              className="p-2 rounded-lg hover:bg-muted transition-colors"
              disabled={busy}
            >
              <ArrowLeft className="w-5 h-5 text-foreground" />
            </button>
            <div>
              <h1 className="text-2xl md:text-3xl text-foreground">Session Capture</h1>

              <div className="mt-2 flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                <span className="px-2 py-1 bg-card border border-border rounded">
                  {headerLang} • {headerLevel}
                </span>
                <span className="px-2 py-1 bg-card border border-border rounded">
                  {allWords.length} words found
                </span>
                <span className="px-2 py-1 bg-card border border-border rounded">
                  session: {sessionId}
                </span>
                {lastTaskId ? (
                  <span className="px-2 py-1 bg-card border border-border rounded">
                    task: {lastTaskId}
                  </span>
                ) : null}
                {stage || status ? (
                  <span className="px-2 py-1 bg-card border border-border rounded">
                    {stage ? `stage: ${stage}` : ""}
                    {stage && status ? " • " : ""}
                    {status ? `status: ${status}` : ""}
                  </span>
                ) : null}
                {error ? (
                  <span className="px-2 py-1 bg-destructive/10 border border-destructive/30 rounded text-destructive">
                    {error}
                  </span>
                ) : null}
              </div>
            </div>
          </div>

          <div className="text-right shrink-0">
            <div className="flex items-center justify-end gap-2 text-muted-foreground">
              <Clock className="w-4 h-4" />
              <span className="text-xs">Time left</span>
            </div>
            <div className="text-xl md:text-2xl font-semibold text-foreground">
              {formatMMSS(remaining)}
            </div>
          </div>
        </div>

        <div className="mb-5 md:mb-6">
          <div className="flex justify-between text-xs text-muted-foreground mb-2">
            <span>{allWords.length} words found</span>
            <span>{remaining > 0 ? "Capture in progress" : "Time is up"}</span>
          </div>
          <div className="h-2 bg-card border border-border rounded-full overflow-hidden">
            <motion.div
              className="h-full bg-gradient-to-r from-primary to-accent"
              initial={{ width: 0 }}
              animate={{ width: `${progressPct}%` }}
              transition={{ duration: 0.25 }}
            />
          </div>
        </div>

        {!hasWords ? (
          <motion.div initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }} className="space-y-6">
            <div className="flex justify-center">
              <PifMascot
                message={
                  status === "processing"
                    ? (stage === "translate"
                        ? "Translating… waiting for WS words."
                        : "OCR processing… waiting for WS done → then auto-translate."
                      )
                    : "Upload an image/PDF. OCR runs first, then client triggers translate after WS OCR done."
                }
                variant="happy"
                size="lg"
              />
            </div>

            <div className="bg-card border border-border rounded-xl shadow p-4 md:p-6">
              <button
                onClick={handleUploadClick}
                className="w-full border-2 border-dashed border-primary/30 rounded-xl p-10 md:p-12
                           hover:border-primary/50 hover:bg-muted/20 transition-all text-center"
                disabled={busy}
              >
                <div className="flex flex-col items-center gap-4">
                  <div className="p-3 rounded-xl bg-primary/10">
                    <Upload className="w-7 h-7 text-primary" />
                  </div>
                  <div>
                    <div className="text-base md:text-lg font-semibold text-foreground">
                      Upload image / PDF
                    </div>
                    <div className="text-sm text-muted-foreground mt-1">
                      REST 202 → WS processing/done → auto POST translate → WS words
                    </div>
                  </div>
                  <div className="px-6 py-2 rounded-lg bg-primary text-primary-foreground font-semibold">
                    {busy ? "Processing…" : "Choose File"}
                  </div>
                </div>
              </button>

              <input
                ref={fileRef}
                type="file"
                accept="image/*,.pdf"
                className="hidden"
                onChange={handleFileChange}
              />
            </div>

            <div className="bg-card border border-border rounded-xl shadow p-4 md:p-6 opacity-60">
              <div className="flex items-center gap-2 mb-3">
                <FileText className="w-4 h-4 text-muted-foreground" />
                <span className="text-sm font-semibold text-foreground">
                  Paste text (not wired yet)
                </span>
              </div>

              <textarea
                value={""}
                readOnly
                placeholder="Coming soon…"
                className="w-full h-28 p-4 bg-input-background border border-border rounded-lg resize-none"
              />
            </div>

            <div className="flex justify-end gap-2">
              <button
                onClick={handleClear}
                disabled={busy}
                className="px-5 py-2 rounded-lg font-semibold transition-all
                           bg-card border border-border hover:bg-muted disabled:opacity-50"
              >
                Clear
              </button>

              <button
                onClick={() => void handleEndSession(false)}
                disabled={busy}
                className="px-5 py-2 rounded-lg font-semibold transition-all
                           bg-secondary text-secondary-foreground hover:shadow-lg disabled:opacity-50"
              >
                End session
              </button>
            </div>
          </motion.div>
        ) : (
          <motion.div initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }} className="space-y-4">
            <div className="bg-card border border-border rounded-xl shadow p-4 md:p-5">
              <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
                <div>
                  <p className="text-sm text-muted-foreground">Words collected</p>
                  <p className="text-2xl font-semibold text-foreground">{allWords.length} total</p>
                </div>

                <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-2 sm:justify-end">
                  <button
                    onClick={handleClear}
                    disabled={busy}
                    className="w-full sm:w-auto px-4 py-2 bg-card border border-border rounded-lg hover:bg-muted transition-colors disabled:opacity-50"
                  >
                    Clear
                  </button>

                  <button
                    onClick={() => void handleEndSession(false)}
                    disabled={busy}
                    className="w-full sm:w-auto px-4 py-2 rounded-lg font-semibold transition-all bg-secondary text-secondary-foreground hover:shadow-lg disabled:opacity-50"
                  >
                    End session
                  </button>
                </div>
              </div>
            </div>

            <div className="bg-card border border-border rounded-xl shadow p-4 md:p-6">
              <h2 className="text-lg font-semibold text-foreground">Upload more</h2>

              <button
                onClick={handleUploadClick}
                disabled={busy}
                className="mt-4 w-full border-2 border-dashed border-primary/30 rounded-xl p-5 md:p-6
                           hover:border-primary/50 hover:bg-muted/20 transition-all text-left disabled:opacity-50"
              >
                <div className="flex items-start gap-4">
                  <div className="p-3 rounded-xl bg-primary/10">
                    <Upload className="w-6 h-6 text-primary" />
                  </div>
                  <div className="flex-1">
                    <div className="font-semibold text-foreground">Click to upload an image</div>
                    <div className="text-xs text-muted-foreground mt-1">
                      OCR → WS done → auto translate → words merge.
                    </div>
                  </div>
                  <div className="text-xs text-muted-foreground">{busy ? "Processing…" : "Ready"}</div>
                </div>
              </button>

              <input
                ref={fileRef}
                type="file"
                accept="image/*,.pdf"
                className="hidden"
                onChange={handleFileChange}
              />
            </div>

            <div className="bg-card border border-border rounded-xl shadow p-4 md:p-6">
              <div className="flex items-center justify-between mb-3">
                <h2 className="text-lg font-semibold text-foreground">Words (with translations)</h2>
                <span className="text-xs text-muted-foreground">From WS translate done</span>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                {allWords.map((w, idx) => (
                  <div
                    key={String((w as any).id ?? `${w.word}-${idx}`)}
                    className="bg-muted/20 border border-border rounded-lg p-4"
                  >
                    <div className="font-semibold text-foreground">{(w as any).word ?? ""}</div>
                    <div className="text-sm text-muted-foreground">{(w as any).translation ?? ""}</div>
                  </div>
                ))}
              </div>
            </div>
          </motion.div>
        )}
      </div>
    </div>
  );
}
