// src/app/components/session-detail.tsx
import { motion, AnimatePresence } from "motion/react";
import {
  ArrowLeft,
  BookOpen,
  Target,
  Calendar,
  Play,
  Brain,
  ChevronLeft,
  ChevronRight,
  Grid3x3,
  Shuffle,
  RotateCw,
} from "lucide-react";
import * as React from "react";
import { useState } from "react";

import type { SessionRecord, SessionWord } from "./session-types";
import { apiFetch } from "./api/http";
import { routes } from "./api/routes";

interface SessionDetailProps {
  session: SessionRecord;
  onBack: () => void;
  onContinueFlashcards: () => void;
  onStartTest: () => void;
  onReviewWords: () => void;
}

export function SessionDetail({ session, onBack, onContinueFlashcards, onStartTest }: SessionDetailProps) {
  const [currentCardIndex, setCurrentCardIndex] = useState(0);
  const [isFlipped, setIsFlipped] = useState(false);

  const sessionId = String((session as any).id ?? (session as any).Id ?? "");
  const title = String((session as any).title ?? (session as any).name ?? `Session ${sessionId}`);

  // даты: используем то, что реально есть в SessionRecord
  const createdAt =
    (session as any).createdAt ??
    (session as any).started_at ??
    (session as any).StartedAt ??
    null;

  const lastStudiedAt =
    (session as any).lastStudiedAt ??
    (session as any).ended_at ??
    (session as any).EndedAt ??
    null;

  const dateToShow = lastStudiedAt || createdAt;

  const accuracy =
    typeof (session as any).accuracy === "number"
      ? (session as any).accuracy
      : Number((session as any).Accuracy ?? 0);

  const words: SessionWord[] = ((session as any).words ?? []) as SessionWord[];
  const wordsCount = Number((session as any).wordsCount ?? words.length);

  const handlePrevCard = () => {
    if (words.length === 0) return;
    setCurrentCardIndex((prev) => (prev > 0 ? prev - 1 : words.length - 1));
    setIsFlipped(false);
  };

  const handleNextCard = () => {
    if (words.length === 0) return;
    setCurrentCardIndex((prev) => (prev < words.length - 1 ? prev + 1 : 0));
    setIsFlipped(false);
  };

  const currentWord = words[currentCardIndex];

  const handleFlashcards = async () => {
    await apiFetch(routes.getFlashcards(sessionId), { method: "GET" });
    onContinueFlashcards();
  };

  const handleTest = async () => {
    await apiFetch(routes.getQuiz(sessionId), { method: "GET" });
    onStartTest();
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="border-b border-border bg-card">
        <div className="max-w-6xl mx-auto px-4 md:px-8 py-4">
          <button
            onClick={onBack}
            className="text-muted-foreground hover:text-foreground transition-colors mb-3 inline-flex items-center gap-2"
          >
            <ArrowLeft className="w-4 h-4" />
            Back to Sessions
          </button>

          <div className="flex items-start justify-between mb-4">
            <div>
              <h1 className="text-2xl md:text-3xl font-semibold text-foreground mb-2">{title}</h1>

              <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
                <div className="flex items-center gap-1.5">
                  <Calendar className="w-4 h-4" />
                  <span>
                    {dateToShow
                      ? new Date(dateToShow).toLocaleDateString("en-US", {
                          month: "long",
                          day: "numeric",
                          year: "numeric",
                        })
                      : "—"}
                  </span>
                </div>

                <div className="flex items-center gap-1.5">
                  <BookOpen className="w-4 h-4" />
                  <span>{wordsCount} words</span>
                </div>

                {Number.isFinite(accuracy) && (
                  <div className="flex items-center gap-1.5">
                    <Target className="w-4 h-4 text-primary" />
                    <span className="text-primary font-semibold">{accuracy}% accuracy</span>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Action buttons */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            <button
              onClick={handleFlashcards}
              className="py-3 px-4 rounded-lg bg-primary/10 border border-primary/30 hover:bg-primary/20 transition-all flex items-center justify-center gap-2 group"
            >
              <Brain className="w-5 h-5 text-primary" />
              <span className="font-medium text-foreground">Flashcards</span>
            </button>

            <button
              className="py-3 px-4 rounded-lg bg-card border border-border hover:bg-muted/40 transition-all flex items-center justify-center gap-2"
              // Learn пока без сетевой логики (как в макете)
            >
              <BookOpen className="w-5 h-5 text-muted-foreground" />
              <span className="font-medium text-foreground">Learn</span>
            </button>

            <button
              onClick={handleTest}
              className="py-3 px-4 rounded-lg bg-card border border-border hover:bg-muted/40 transition-all flex items-center justify-center gap-2"
            >
              <Play className="w-5 h-5 text-muted-foreground" />
              <span className="font-medium text-foreground">Test</span>
            </button>

            <button
              className="py-3 px-4 rounded-lg bg-card border border-border hover:bg-muted/40 transition-all flex items-center justify-center gap-2"
              // Match пока без логики
            >
              <Grid3x3 className="w-5 h-5 text-muted-foreground" />
              <span className="font-medium text-foreground">Match</span>
            </button>
          </div>
        </div>
      </div>

      {/* Main */}
      <div className="max-w-6xl mx-auto px-4 md:px-8 py-8">
        {words.length > 0 ? (
          <>
            {/* Big card */}
            <div className="mb-8">
              <div className="flex items-center justify-center mb-6">
                <motion.div
                  className="relative w-full max-w-3xl"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                >
                  <motion.button
                    onClick={() => setIsFlipped((v) => !v)}
                    className="w-full bg-card border border-border rounded-2xl shadow-lg hover:shadow-xl transition-shadow cursor-pointer"
                    style={{ minHeight: "320px" }}
                    whileHover={{ scale: 1.01 }}
                    whileTap={{ scale: 0.99 }}
                  >
                    <AnimatePresence mode="wait">
                      {!isFlipped ? (
                        <motion.div
                          key="front"
                          initial={{ rotateY: 90, opacity: 0 }}
                          animate={{ rotateY: 0, opacity: 1 }}
                          exit={{ rotateY: -90, opacity: 0 }}
                          transition={{ duration: 0.2 }}
                          className="p-8 md:p-12 flex flex-col items-center justify-center"
                          style={{ minHeight: "320px" }}
                        >
                          <p className="text-xs text-muted-foreground mb-4">WORD</p>
                          <p className="text-3xl md:text-4xl font-semibold text-foreground text-center">
                            {(currentWord as any)?.word ?? ""}
                          </p>
                        </motion.div>
                      ) : (
                        <motion.div
                          key="back"
                          initial={{ rotateY: 90, opacity: 0 }}
                          animate={{ rotateY: 0, opacity: 1 }}
                          exit={{ rotateY: -90, opacity: 0 }}
                          transition={{ duration: 0.2 }}
                          className="p-8 md:p-12 flex flex-col items-center justify-center"
                          style={{ minHeight: "320px" }}
                        >
                          <p className="text-xs text-muted-foreground mb-4">TRANSLATION</p>
                          <p className="text-3xl md:text-4xl font-semibold text-foreground text-center">
                            {(currentWord as any)?.translation ?? ""}
                          </p>
                        </motion.div>
                      )}
                    </AnimatePresence>
                  </motion.button>

                  <div className="text-center mt-4">
                    <p className="text-sm text-muted-foreground">
                      Click card to flip •
                      <span className="ml-1 px-2 py-0.5 rounded bg-muted/50 border border-border text-xs">
                        SPACE
                      </span>
                    </p>
                  </div>
                </motion.div>
              </div>

              {/* Navigation */}
              <div className="flex items-center justify-center gap-4 mb-8">
                <button
                  onClick={handlePrevCard}
                  className="p-3 rounded-lg bg-card border border-border hover:bg-muted/40 transition-all"
                  disabled={words.length <= 1}
                >
                  <ChevronLeft className="w-5 h-5 text-foreground" />
                </button>

                <div className="flex items-center gap-3">
                  <span className="text-sm font-medium text-foreground">
                    {currentCardIndex + 1} / {words.length}
                  </span>

                  <button className="p-2 rounded-lg hover:bg-muted/40 transition-all" type="button">
                    <Shuffle className="w-4 h-4 text-muted-foreground" />
                  </button>

                  <button className="p-2 rounded-lg hover:bg-muted/40 transition-all" type="button">
                    <RotateCw className="w-4 h-4 text-muted-foreground" />
                  </button>
                </div>

                <button
                  onClick={handleNextCard}
                  className="p-3 rounded-lg bg-card border border-border hover:bg-muted/40 transition-all"
                  disabled={words.length <= 1}
                >
                  <ChevronRight className="w-5 h-5 text-foreground" />
                </button>
              </div>
            </div>

            {/* Words list */}
            <div className="mt-12">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-semibold text-foreground">Words in this set ({words.length})</h2>
              </div>

              <div className="space-y-3">
                {words.map((word: any, index: number) => (
                  <motion.button
                    key={String(word?.id ?? index)}
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.03 }}
                    onClick={() => {
                      setCurrentCardIndex(index);
                      setIsFlipped(false);
                      window.scrollTo({ top: 0, behavior: "smooth" });
                    }}
                    className={`w-full bg-card border rounded-xl p-5 hover:shadow-md transition-all text-left group ${
                      index === currentCardIndex ? "border-primary/50 shadow-sm" : "border-border"
                    }`}
                  >
                    <div className="flex items-start gap-6">
                      <div className="flex-1 min-w-0">
                        <p className="text-xs text-muted-foreground mb-1">Word</p>
                        <p className="text-lg font-semibold text-foreground mb-3">{word?.word ?? ""}</p>

                        <div className="border-t border-border pt-3">
                          <p className="text-xs text-muted-foreground mb-1">Translation</p>
                          <p className="text-base text-foreground">{word?.translation ?? ""}</p>
                        </div>
                      </div>

                      <div className="flex flex-col items-end gap-2">
                        {!!word?.known && (
                          <span className="text-xs px-2 py-1 rounded-full bg-secondary/15 border border-secondary/40 text-secondary">
                            Known
                          </span>
                        )}
                        {index === currentCardIndex && <div className="w-2 h-2 rounded-full bg-primary mt-2" />}
                      </div>
                    </div>
                  </motion.button>
                ))}
              </div>
            </div>
          </>
        ) : (
          <div className="text-center py-20 text-muted-foreground">
            <BookOpen className="w-16 h-16 mx-auto mb-4 opacity-40" />
            <p className="text-lg">No words saved for this session yet.</p>
          </div>
        )}
      </div>
    </div>
  );
}
