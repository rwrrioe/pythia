// src/app/components/flashcards-review.tsx
import * as React from "react";
import { useMemo } from "react";
import { motion, AnimatePresence } from "motion/react";
import { ArrowLeft, RotateCcw, Check, X } from "lucide-react";

import type { SessionWord } from "./session-types";
import { getFlashcards } from "./api/sessions";
import { GreekPattern } from "./greek-pattern";

type Props = {
  sessionId: number | null;
  words?: SessionWord[]; // можно передать из App, чтобы не дергать второй раз
  onDone: () => void;
};

export function FlashcardsReview({ sessionId, words: wordsProp, onDone }: Props) {
  const [words, setWords] = React.useState<SessionWord[]>(wordsProp ?? []);
  const [loading, setLoading] = React.useState<boolean>(!!sessionId && !wordsProp?.length);
  const [error, setError] = React.useState<string | null>(null);

  const [i, setI] = React.useState(0);
  const [flipped, setFlipped] = React.useState(false);

  React.useEffect(() => {
    if (!sessionId) {
      setLoading(false);
      return;
    }
    if (wordsProp?.length) {
      setWords(wordsProp);
      setLoading(false);
      return;
    }

    (async () => {
      try {
        setLoading(true);
        const res = await getFlashcards(sessionId);
        setWords(res.words);
      } catch (e: any) {
        setError(e?.message ?? "Failed to load flashcards");
      } finally {
        setLoading(false);
      }
    })();
  }, [sessionId, wordsProp]);

  const current = words[i];

  const progress = useMemo(() => ((i + 1) / Math.max(1, words.length)) * 100, [i, words.length]);

  if (!sessionId) return <div className="opacity-70">No session selected.</div>;

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background p-4 md:p-8">
        <div className="max-w-4xl mx-auto pt-8 opacity-70">Loading flashcards…</div>
      </div>
    );
  }

  if (error) return <div className="text-red-500">{error}</div>;

  if (!words.length) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background p-4 md:p-8">
        <div className="max-w-4xl mx-auto pt-8 opacity-70">No flashcards yet.</div>
      </div>
    );
  }

  const next = () => {
    if (i + 1 >= words.length) onDone();
    else {
      setI((v) => v + 1);
      setFlipped(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background p-4 md:p-8 relative overflow-hidden">
      <div className="absolute top-0 left-0 right-0">
        <GreekPattern className="w-full h-6 text-primary opacity-40" />
      </div>

      <div className="max-w-4xl mx-auto relative z-10 pt-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <button
            onClick={onDone}
            className="p-2 rounded-lg hover:bg-muted transition-colors"
            title="Back"
          >
            <ArrowLeft className="w-6 h-6 text-foreground" />
          </button>

          <div className="text-right">
            <p className="text-xs text-muted-foreground">Flashcards</p>
            <p className="text-sm text-muted-foreground">
              {i + 1} / {words.length}
            </p>
          </div>
        </div>

        {/* Progress */}
        <div className="mb-6">
          <div className="h-3 bg-card border border-border rounded-full overflow-hidden">
            <motion.div
              className="h-full bg-gradient-to-r from-primary to-accent"
              initial={{ width: 0 }}
              animate={{ width: `${progress}%` }}
              transition={{ duration: 0.25 }}
            />
          </div>
        </div>

        {/* Card */}
        <AnimatePresence mode="wait">
          <motion.div
            key={String((current as any)?.id ?? i)}
            initial={{ opacity: 0, y: 14, scale: 0.995 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.2 }}
            className="space-y-5"
          >
            <div className="perspective-1000">
              <motion.div
                className="relative w-full h-80 md:h-96 cursor-pointer"
                onClick={() => setFlipped((v) => !v)}
                animate={{ rotateY: flipped ? 180 : 0 }}
                transition={{ duration: 0.55 }}
                style={{ transformStyle: "preserve-3d" }}
              >
                {/* Front */}
                <div
                  className="absolute inset-0 bg-card border-2 border-primary/30 rounded-2xl p-8 shadow-xl flex flex-col items-center justify-center backface-hidden"
                  style={{ backfaceVisibility: "hidden" }}
                >
                  <p className="text-sm text-muted-foreground mb-4 uppercase tracking-wide">Word</p>
                  <h2 className="text-5xl md:text-6xl font-bold text-primary">
                    {(current as any)?.word ?? ""}
                  </h2>
                  <p className="text-muted-foreground mt-8 flex items-center gap-2">
                    <RotateCcw className="w-4 h-4" />
                    Click to reveal
                  </p>
                </div>

                {/* Back */}
                <div
                  className="absolute inset-0 bg-gradient-to-br from-accent/20 to-primary/20 border-2 border-primary/50 rounded-2xl p-8 shadow-xl flex flex-col items-center justify-center backface-hidden"
                  style={{ backfaceVisibility: "hidden", transform: "rotateY(180deg)" }}
                >
                  <p className="text-sm text-muted-foreground mb-4 uppercase tracking-wide">Translation</p>
                  <h2 className="text-5xl md:text-6xl font-bold text-primary mb-6">
                    {(current as any)?.translation ?? ""}
                  </h2>

                  {(current as any)?.example ? (
                    <div className="bg-card/80 rounded-lg p-4 max-w-md">
                      <p className="text-sm text-foreground italic">{(current as any)?.example}</p>
                    </div>
                  ) : null}
                </div>
              </motion.div>
            </div>

            {/* Actions (функционал тот же: Start quiz и Next) */}
            <div className="flex gap-3 justify-center">
              <button
                onClick={onDone}
                className="px-6 py-3 bg-card border-2 border-destructive/50 text-destructive rounded-lg font-semibold hover:bg-destructive/10 transition-all inline-flex items-center gap-2"
              >
                <X className="w-5 h-5" />
                Start quiz
              </button>

              <motion.button
                onClick={next}
                whileHover={{ scale: 1.01 }}
                whileTap={{ scale: 0.99 }}
                className="px-6 py-3 bg-gradient-to-r from-secondary to-secondary/85 text-secondary-foreground rounded-lg font-semibold hover:shadow-md transition-all inline-flex items-center gap-2"
              >
                <Check className="w-5 h-5" />
                Next
              </motion.button>
            </div>
          </motion.div>
        </AnimatePresence>
      </div>
    </div>
  );
}
