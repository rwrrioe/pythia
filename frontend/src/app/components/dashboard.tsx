import * as React from "react";
import { useEffect, useMemo, useState } from "react";
import { motion } from "motion/react";
import { Plus, Flame, Target, BookOpen, Calendar } from "lucide-react";
import { useNavigate } from "react-router-dom";

import { apiFetch } from "./api/http";
import { routes } from "./api/routes";
import type { SessionRecord } from "./session-types";

type DashboardApi = {
  dashboard?: {
    streak?: number;
    words_learned?: number;
    accuracy?: number;
    latest_sessions?: any[];
  };
};

function mapSession(s: any): SessionRecord {
  return {
    id: String(s?.id ?? s?.Id ?? ""),
    name: String(s?.name ?? s?.Name ?? ""),
    status: String(s?.status ?? s?.Status ?? ""),
    language: s?.language ?? s?.Language,
    level: s?.level ?? s?.Level,
    accuracy: s?.accuracy ?? s?.Accuracy ?? 0,
    started_at: s?.started_at ?? s?.StartedAt,
    ended_at: s?.ended_at ?? s?.EndedAt,
    duration: s?.duration ?? s?.Duration ?? 0,
    wordsCount: s?.wordsCount ?? s?.WordsCount ?? 0,
  } as any;
}

export function Dashboard() {
  const navigate = useNavigate();

  const [sessions, setSessions] = useState<SessionRecord[]>([]);
  const [stats, setStats] = useState<{ streak: number; words_learned: number; accuracy: number } | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const res = await apiFetch<DashboardApi>(routes.dashboard, { method: "GET" });
        const d = res?.dashboard ?? {};
        setStats({
          streak: Number(d?.streak ?? 0),
          words_learned: Number(d?.words_learned ?? 0),
          accuracy: Number(d?.accuracy ?? 0),
        });

        const ls = Array.isArray(d?.latest_sessions) ? d.latest_sessions : [];
        setSessions(ls.map(mapSession));
      } catch (e) {
        console.error(e);
        setStats(null);
        setSessions([]);
      }
    })();
  }, []);

  const tiles = useMemo(() => {
    return [
      { label: "Current Streak", value: stats?.streak ?? 0, icon: Flame, suffix: "days" },
      { label: "Words Learned", value: stats?.words_learned ?? 0, icon: BookOpen, suffix: "" },
      { label: "Accuracy", value: stats?.accuracy ?? 0, icon: Target, suffix: "%" },
    ] as const;
  }, [stats]);

  const openSession = async (sessionId: string) => {
    // ✅ требование: при клике из dashboard делать GET /api/library/session/:id
    await apiFetch(routes.librarySession(sessionId), { method: "GET" });
    navigate(`/sessions/${sessionId}`);
  };

  return (
    <div className="min-h-screen bg-background p-4 md:p-8">
      <div className="max-w-5xl mx-auto">
        {/* Header */}
        <motion.div initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }} className="mb-6 md:mb-8">
          <h1 className="text-2xl md:text-3xl text-foreground mb-2">Dashboard</h1>
          <p className="text-sm md:text-base text-muted-foreground">
            Welcome back! Ready to continue your learning journey?
          </p>
        </motion.div>

        {/* Stats Grid */}
        <div className="grid grid-cols-3 gap-3 md:gap-4 mb-6 md:mb-8">
          {tiles.map((stat, index) => (
            <motion.div
              key={stat.label}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.05 }}
              className="bg-card border border-border rounded-lg p-4 md:p-6 shadow"
            >
              <div className="flex items-center gap-2 md:gap-3 mb-2 md:mb-3">
                <div className="p-1.5 md:p-2 rounded-lg bg-primary/10">
                  <stat.icon className="w-4 h-4 md:w-5 md:h-5 text-primary" />
                </div>
              </div>
              <p className="text-2xl md:text-3xl font-semibold text-foreground">
                {stat.value}
                {stat.suffix}
              </p>
              <p className="text-xs md:text-sm text-muted-foreground mt-1">{stat.label}</p>
            </motion.div>
          ))}
        </div>

        {/* Main Content Grid */}
        <div className="grid md:grid-cols-3 gap-4 md:gap-6">
          {/* New Session Card */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="md:col-span-1"
          >
            <button
              onClick={() => navigate("/session/new")}
              className="w-full h-full min-h-[200px] md:min-h-[280px] bg-gradient-to-br from-primary/10 to-accent/10 
                         border-2 border-primary/30 rounded-xl p-6 
                         hover:shadow-lg hover:scale-[1.02] transition-all duration-300
                         flex flex-col items-center justify-center gap-3 md:gap-4 group"
            >
              <div
                className="w-12 h-12 md:w-16 md:h-16 rounded-full bg-gradient-to-br from-primary to-accent 
                           flex items-center justify-center shadow-lg 
                           group-hover:scale-110 transition-transform"
              >
                <Plus className="w-6 h-6 md:w-8 md:h-8 text-white" />
              </div>
              <div className="text-center">
                <h2 className="text-lg md:text-xl font-semibold text-foreground mb-1 md:mb-2">New Session</h2>
                <p className="text-xs md:text-sm text-muted-foreground">Start learning new words</p>
              </div>
            </button>
          </motion.div>

          {/* Recent Sessions List */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="md:col-span-2 bg-card border border-border rounded-xl p-4 md:p-6 shadow"
          >
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-base md:text-lg font-semibold text-foreground">Recent Sessions</h2>
              <Calendar className="w-4 h-4 md:w-5 md:h-5 text-muted-foreground" />
            </div>

            <div className="space-y-2 md:space-y-3">
              {sessions.length === 0 ? (
                <div className="text-sm text-muted-foreground p-4 border border-dashed rounded-lg bg-muted/20">
                  No sessions yet.
                </div>
              ) : (
                sessions.slice(0, 4).map((s, index) => {
                  const dateISO = (s as any)?.started_at ?? (s as any)?.StartedAt ?? new Date().toISOString();
                  const title = (s as any)?.name?.trim?.() ? (s as any).name : `Session #${s.id}`;
                  const words = Number((s as any)?.wordsCount ?? 0);
                  const accuracy = Number((s as any)?.accuracy ?? 0);

                  return (
                    <motion.button
                      key={s.id}
                      initial={{ opacity: 0, x: -10 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: 0.4 + index * 0.05 }}
                      onClick={() => openSession(String(s.id))}
                      className="w-full flex items-center justify-between p-3 bg-muted/30 rounded-lg 
                                 border border-border/50 hover:bg-muted/50 transition-colors cursor-pointer text-left"
                    >
                      <div className="flex-1 min-w-0">
                        <p className="font-medium text-foreground text-xs md:text-sm truncate">{title}</p>
                        <div className="flex items-center gap-2 md:gap-3 mt-1">
                          <p className="text-xs text-muted-foreground">
                            {new Date(dateISO).toLocaleDateString("en-US", { month: "short", day: "numeric" })}
                          </p>
                          <p className="text-xs text-muted-foreground">{words} words</p>
                        </div>
                      </div>
                      <div className="text-right ml-2">
                        <p className="text-base md:text-lg font-semibold text-primary">{accuracy}%</p>
                      </div>
                    </motion.button>
                  );
                })
              )}
            </div>
          </motion.div>
        </div>

        {/* Next Review (как в макете; пока статично, сеть не трогаю) */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5 }}
          className="mt-4 md:mt-6 bg-secondary/10 border border-secondary/30 rounded-lg p-4 md:p-5"
        >
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs md:text-sm text-muted-foreground">Next Review Session</p>
              <p className="text-lg md:text-xl font-semibold text-foreground mt-1">Tomorrow, January 18</p>
            </div>
            <div className="text-right">
              <p className="text-2xl md:text-3xl font-semibold text-secondary">24</p>
              <p className="text-xs text-muted-foreground">words ready</p>
            </div>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
