import React, { useMemo, useState } from "react";
import { Calendar, TrendingUp, BookOpen, Search, ArrowLeft } from "lucide-react";

interface LibraryProps {
  onBack: () => void;
}

export function Library({ onBack }: LibraryProps) {
  const [searchQuery, setSearchQuery] = useState("");

  const sessions = useMemo(
    () => [
      { id: 1, date: "2026-01-14", title: "German News Article", words: 12, accuracy: 92, language: "German", status: "completed" },
      { id: 2, date: "2026-01-13", title: "Business Email", words: 15, accuracy: 85, language: "German", status: "completed" },
      { id: 3, date: "2026-01-12", title: "Travel Blog Post", words: 10, accuracy: 90, language: "German", status: "completed" },
      { id: 4, date: "2026-01-11", title: "Technical Document", words: 14, accuracy: 78, language: "German", status: "completed" },
      { id: 5, date: "2026-01-10", title: "Short Story", words: 13, accuracy: 88, language: "German", status: "completed" },
      { id: 6, date: "2026-01-09", title: "Recipe Instructions", words: 11, accuracy: 95, language: "German", status: "completed" },
    ],
    []
  );

  const filteredSessions = useMemo(() => {
    const q = searchQuery.trim().toLowerCase();
    if (!q) return sessions;
    return sessions.filter((s) => s.title.toLowerCase().includes(q));
  }, [searchQuery, sessions]);

  return (
    <div className="min-h-screen bg-background p-6 md:p-8">
      <div className="max-w-5xl mx-auto">
        {/* Header */}
        <div className="mb-6 md:mb-8 flex items-start gap-3">
          <button
            type="button"
            onClick={onBack}
            className="mt-1 inline-flex items-center gap-2 rounded-lg border border-border bg-card px-3 py-2 text-sm text-foreground hover:bg-muted"
          >
            <ArrowLeft className="w-4 h-4" />
            Back
          </button>

          <div className="flex-1">
            <h1 className="text-3xl text-foreground mb-1">Session Library</h1>
            <p className="text-muted-foreground">Review all your past learning sessions and track your progress.</p>
          </div>
        </div>

        {/* Search */}
        <div className="mb-6">
          <div className="relative">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
            <input
              type="text"
              placeholder="Search sessions..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-12 pr-4 py-3 bg-card border border-border rounded-lg text-foreground
                         focus:outline-none focus:ring-2 focus:ring-primary/50"
            />
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
          <div className="bg-card border border-border rounded-lg p-4 text-center">
            <BookOpen className="w-6 h-6 text-primary mx-auto mb-2" />
            <p className="text-2xl font-semibold text-foreground">{sessions.length}</p>
            <p className="text-sm text-muted-foreground">Total Sessions</p>
          </div>
          <div className="bg-card border border-border rounded-lg p-4 text-center">
            <TrendingUp className="w-6 h-6 text-secondary mx-auto mb-2" />
            <p className="text-2xl font-semibold text-foreground">87%</p>
            <p className="text-sm text-muted-foreground">Avg. Accuracy</p>
          </div>
          <div className="bg-card border border-border rounded-lg p-4 text-center">
            <Calendar className="w-6 h-6 text-accent mx-auto mb-2" />
            <p className="text-2xl font-semibold text-foreground">142</p>
            <p className="text-sm text-muted-foreground">Words Learned</p>
          </div>
        </div>

        {/* List */}
        <div className="space-y-3">
          {filteredSessions.length === 0 ? (
            <div className="text-center py-12 text-muted-foreground">
              No sessions found matching "{searchQuery}"
            </div>
          ) : (
            filteredSessions.map((session) => (
              <div
                key={session.id}
                className="bg-card border border-border rounded-lg p-5 hover:shadow-lg transition-shadow cursor-pointer group"
              >
                <div className="flex items-center justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex flex-wrap items-center gap-3 mb-2">
                      <h3 className="font-semibold text-foreground group-hover:text-primary transition-colors truncate">
                        {session.title}
                      </h3>
                      <span className="px-2 py-1 bg-primary/10 text-primary text-xs rounded">
                        {session.language}
                      </span>
                    </div>

                    <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
                      <div className="flex items-center gap-1">
                        <Calendar className="w-4 h-4" />
                        {new Date(session.date).toLocaleDateString("en-US", {
                          month: "short",
                          day: "numeric",
                          year: "numeric",
                        })}
                      </div>
                      <div className="flex items-center gap-1">
                        <BookOpen className="w-4 h-4" />
                        {session.words} words
                      </div>
                    </div>
                  </div>

                  <div className="text-right shrink-0">
                    <p className="text-2xl font-semibold text-primary">{session.accuracy}%</p>
                    <p className="text-xs text-muted-foreground">accuracy</p>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
