// src/app/components/session-page.tsx
import * as React from "react";
import { motion } from "motion/react";
import { Search, ChevronDown, User } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import type { SessionRecord as Session } from "./session-types";

import { apiFetch } from "./api/http";
import { routes } from "./api/routes";

/**
 * ✅ Fix "2000 requests": module-level cache + in-flight promise.
 * Even if component mounts/unmounts rapidly (StrictMode/dev), we only do ONE request.
 */
let _librarySessionsCache: Session[] | null = null;
let _librarySessionsInFlight: Promise<Session[]> | null = null;

async function loadLibrarySessionsOnce(): Promise<Session[]> {
  if (_librarySessionsCache) return _librarySessionsCache;
  if (_librarySessionsInFlight) return _librarySessionsInFlight;

  _librarySessionsInFlight = (async () => {
    const res = await apiFetch<{ sessions: any[] }>(routes.librarySessions, { method: "GET" });
    const list = Array.isArray(res?.sessions) ? res.sessions : [];
    const mapped = list.map(mapLibraryListItemToSession) as Session[];
    _librarySessionsCache = mapped;
    return mapped;
  })().finally(() => {
    _librarySessionsInFlight = null;
  });

  return _librarySessionsInFlight;
}

interface SessionsPageProps {
  sessions?: Session[];
  // ✅ ВАЖНО: App.tsx ожидает id, и сам делает GET /library/session/:id
  onOpenSession: (id: string) => void;
}

type TabType = "flashcards" | "practice" | "folders";

const langIdToCode = (id: any): string => {
  const n = Number(id);
  if (n === 1) return "en";
  if (n === 2) return "de";
  if (n === 3) return "fr";
  if (n === 4) return "es";
  return "en";
};

const levelToDifficulty = (lvl: any): string => {
  const n = Number(lvl);
  const map = ["A1", "A2", "B1", "B2", "C1", "C2"];
  return map[n] ?? "A1";
};

// list item from GET /library/session
function mapLibraryListItemToSession(dto: any): Session {
  const id = String(dto?.Id ?? dto?.id ?? "");
  const started = dto?.StartedAt ?? dto?.started_at ?? new Date().toISOString();
  const ended = dto?.ended_at ?? dto?.EndedAt ?? dto?.endedAt ?? undefined;

  const title = dto?.name && String(dto.name).trim().length > 0 ? String(dto.name) : `Session #${id}`;

  return {
    id,
    title,
    // эти поля SessionRecord не содержит — у тебя они всё равно используются "as any" в других местах
    language: langIdToCode(dto?.language),
    difficulty: levelToDifficulty(dto?.level),
    createdAt: started,
    lastStudiedAt: ended || started,
    wordsCount: Number(dto?.wordsCount ?? 0),
    accuracy: Number(dto?.Accuracy ?? dto?.accuracy ?? 0),
    words: [],
  } as any;
}

export function SessionsPage({ sessions: sessionsProp = [], onOpenSession }: SessionsPageProps) {
  const [q, setQ] = useState("");
  const [activeTab, setActiveTab] = useState<TabType>("flashcards");
  const [sortBy] = useState<"recent" | "title" | "created">("recent");

  const [sessions, setSessions] = useState<Session[]>(sessionsProp);

  // ✅ On mount: load list ONCE for whole app (cache+inflight handles StrictMode)
  useEffect(() => {
    let alive = true;

    (async () => {
      try {
        const list = await loadLibrarySessionsOnce();
        if (!alive) return;
        setSessions(list);
      } catch (e) {
        console.error(e);
        if (!alive) return;
        setSessions([]);
      }
    })();

    return () => {
      alive = false;
    };
  }, []);

  const filtered = useMemo(() => {
    const qq = q.trim().toLowerCase();
    if (!qq) return sessions;

    return sessions.filter((s: any) => {
      return (
        String(s.title ?? "").toLowerCase().includes(qq) ||
        String(s.language ?? "").toLowerCase().includes(qq) ||
        String(s.difficulty ?? "").toLowerCase().includes(qq)
      );
    });
  }, [q, sessions]);

  const groupedSessions = useMemo(() => {
    const groups: Record<string, Session[]> = {};
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    filtered.forEach((session: any) => {
      const d = new Date(session.lastStudiedAt || session.createdAt);
      d.setHours(0, 0, 0, 0);

      let groupKey = "";
      if (d.getTime() === today.getTime()) groupKey = "TODAY";
      else {
        const monthYear = d.toLocaleDateString("en-US", { month: "long", year: "numeric" }).toUpperCase();
        groupKey = `IN ${monthYear}`;
      }

      if (!groups[groupKey]) groups[groupKey] = [];
      groups[groupKey].push(session);
    });

    return groups;
  }, [filtered]);

  const tabs = [
    { id: "flashcards" as TabType, label: "Flashcards" },
    { id: "practice" as TabType, label: "Practice tests" },
    { id: "folders" as TabType, label: "Folders" },
  ];

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-6xl mx-auto px-4 md:px-8 py-6 md:py-8">
        <motion.div initial={{ opacity: 0, y: -10 }} animate={{ opacity: 1, y: 0 }} className="mb-6">
          <h1 className="text-2xl md:text-3xl text-foreground mb-6">Your library</h1>

          <div className="flex items-center gap-6 border-b border-border mb-6">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`pb-3 text-sm font-medium transition-colors relative ${
                  activeTab === tab.id ? "text-foreground" : "text-muted-foreground hover:text-foreground"
                }`}
              >
                {tab.label}
                {activeTab === tab.id && (
                  <motion.div layoutId="activeTab" className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary" />
                )}
              </button>
            ))}
          </div>

          <div className="flex flex-col md:flex-row gap-3 items-start md:items-center justify-between">
            <div className="relative">
              <button className="px-4 py-2 bg-card border border-border rounded-lg text-sm font-medium text-foreground hover:bg-muted/40 transition-colors inline-flex items-center gap-2">
                {sortBy === "recent" ? "Recent" : sortBy}
                <ChevronDown className="w-4 h-4" />
              </button>
            </div>

            <div className="relative w-full md:w-80">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
              <input
                value={q}
                onChange={(e) => setQ(e.target.value)}
                placeholder="Search flashcards"
                className="w-full pl-10 pr-3 py-2 bg-card border border-border rounded-lg text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/40"
              />
            </div>
          </div>
        </motion.div>

        <div className="space-y-8">
          {Object.entries(groupedSessions).map(([groupLabel, groupSessions]) => (
            <div key={groupLabel}>
              <h2 className="text-xs font-semibold text-muted-foreground mb-3 tracking-wide">{groupLabel}</h2>

              <div className="space-y-2">
                {groupSessions.map((session: any) => (
                  <motion.button
                    key={session.id}
                    initial={{ opacity: 0, y: 5 }}
                    animate={{ opacity: 1, y: 0 }}
                    onClick={() => onOpenSession(String(session.id))}
                    className="w-full bg-card hover:bg-muted/30 border border-border rounded-lg p-4 transition-all text-left group"
                  >
                    <div className="flex items-start gap-3">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-xs text-muted-foreground">
                            {session.wordsCount} {session.wordsCount === 1 ? "term" : "terms"}
                          </span>
                          <div className="flex items-center gap-1.5">
                            <div className="w-5 h-5 rounded-full bg-primary/20 flex items-center justify-center">
                              <User className="w-3 h-3 text-primary" />
                            </div>
                            <span className="text-xs text-muted-foreground">Pythia</span>
                          </div>
                          {session.language && (
                            <span className="text-xs px-2 py-0.5 rounded bg-muted/50 border border-border/60 text-muted-foreground">
                              {session.language}
                            </span>
                          )}
                        </div>
                        <h3 className="text-base font-semibold text-foreground group-hover:text-primary transition-colors">
                          {session.title}
                        </h3>
                      </div>

                      {typeof session.accuracy === "number" && (
                        <div className="text-right shrink-0">
                          <p className="text-lg font-semibold text-primary">{session.accuracy}%</p>
                          <p className="text-xs text-muted-foreground">accuracy</p>
                        </div>
                      )}
                    </div>
                  </motion.button>
                ))}
              </div>
            </div>
          ))}
        </div>

        {filtered.length === 0 && (
          <div className="text-center py-12 text-muted-foreground">
            <p>No sessions found{q ? ` matching "${q}"` : ""}.</p>
          </div>
        )}
      </div>
    </div>
  );
}
