// src/app/components/tests-mode.tsx
import * as React from "react";
import { useEffect, useMemo, useState } from "react";
import { motion, AnimatePresence } from "motion/react";
import { ArrowLeft, RotateCcw, Check, X, ChevronRight } from "lucide-react";

import { GreekPattern } from "./greek-pattern";
import { getQuiz } from "./api/sessions";

type Props = {
  sessionId: number | null;
  onComplete: (score: number) => void;
  onBack: () => void;
};

type QuizMode = "flashcard" | "multiple-choice" | "typing";

type QuizQuestion = {
  question: string;
  options: string[];
  answer: string;
};

function normalize(s: any) {
  return String(s ?? "").trim().toLowerCase();
}

export function TestsMode({ sessionId, onComplete, onBack }: Props) {
  const [loading, setLoading] = useState<boolean>(!!sessionId);
  const [error, setError] = useState<string | null>(null);

  const [questions, setQuestions] = useState<QuizQuestion[]>([]);

  const [mode, setMode] = useState<QuizMode>("multiple-choice");
  const [currentIndex, setCurrentIndex] = useState(0);

  const [isFlipped, setIsFlipped] = useState(false);
  const [correctAnswers, setCorrectAnswers] = useState(0);

  const [userAnswer, setUserAnswer] = useState("");
  const [showFeedback, setShowFeedback] = useState(false);
  const [isCorrect, setIsCorrect] = useState(false);

  useEffect(() => {
    if (!sessionId) {
      setLoading(false);
      setQuestions([]);
      return;
    }

    (async () => {
      try {
        setLoading(true);
        setError(null);
        const res = await getQuiz(sessionId); // ✅ HTTP как было
        const list = Array.isArray(res?.questions) ? (res.questions as any[]) : [];
        setQuestions(
          list.map((q: any, i: number) => ({
            question: String(q?.question ?? q?.q ?? `Question ${i + 1}`),
            options: Array.isArray(q?.options) ? q.options.map(String) : [],
            answer: String(q?.answer ?? q?.correct ?? ""),
          }))
        );
      } catch (e: any) {
        setError(e?.message ?? "Failed to load quiz");
        setQuestions([]);
      } finally {
        setLoading(false);
      }
    })();
  }, [sessionId]);

  const total = questions.length;
  const current = questions[currentIndex];

  const progress = useMemo(() => {
    return total ? ((currentIndex + 1) / total) * 100 : 0;
  }, [currentIndex, total]);

  const resetStepState = () => {
    setIsFlipped(false);
    setUserAnswer("");
    setShowFeedback(false);
    setIsCorrect(false);
  };

  const finish = (nextCorrectCount: number) => {
    const score = Math.round((nextCorrectCount / Math.max(1, total)) * 100);
    onComplete(score);
  };

  const next = (didCorrect: boolean) => {
    const nextCorrectCount = didCorrect ? correctAnswers + 1 : correctAnswers;

    if (didCorrect) setCorrectAnswers(nextCorrectCount);

    if (currentIndex < total - 1) {
      setCurrentIndex((v) => v + 1);
      resetStepState();
    } else {
      finish(nextCorrectCount);
    }
  };

  const handleFlip = () => setIsFlipped((v) => !v);

  const handleKnown = () => next(true);
  const handleDontKnow = () => next(false);

  const handleMultipleChoice = (answer: string) => {
    const ok = answer === current?.answer;
    setIsCorrect(ok);
    setShowFeedback(true);
    window.setTimeout(() => next(ok), 1200);
  };

  const handleCheckTyping = () => {
    const ok = normalize(userAnswer) === normalize(current?.answer);
    setIsCorrect(ok);
    setShowFeedback(true);
    window.setTimeout(() => next(ok), 1200);
  };

  if (!sessionId) {
    return <div className="opacity-70 p-4">No session selected.</div>;
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background p-4 md:p-8">
        <div className="max-w-4xl mx-auto pt-8 opacity-70">Loading quiz…</div>
      </div>
    );
  }

  if (error) return <div className="text-red-500 p-4">{error}</div>;

  if (!questions.length) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background p-4 md:p-8">
        <div className="max-w-4xl mx-auto pt-8 opacity-70">No quiz questions yet.</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background p-4 md:p-8 relative overflow-hidden">
      <div className="absolute top-0 left-0 right-0">
        <GreekPattern className="w-full h-6 text-primary opacity-40" />
      </div>

      <div className="max-w-4xl mx-auto relative z-10 pt-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <button onClick={onBack} className="p-2 rounded-lg hover:bg-muted transition-colors">
            <ArrowLeft className="w-6 h-6 text-foreground" />
          </button>

          <div className="flex gap-2">
            <button
              onClick={() => {
                setMode("flashcard");
                resetStepState();
              }}
              className={`px-4 py-2 rounded-lg font-semibold transition-colors ${
                mode === "flashcard" ? "bg-primary text-primary-foreground" : "bg-card text-foreground"
              }`}
            >
              Flashcards
            </button>
            <button
              onClick={() => {
                setMode("multiple-choice");
                resetStepState();
              }}
              className={`px-4 py-2 rounded-lg font-semibold transition-colors ${
                mode === "multiple-choice" ? "bg-primary text-primary-foreground" : "bg-card text-foreground"
              }`}
            >
              Quiz
            </button>
            <button
              onClick={() => {
                setMode("typing");
                resetStepState();
              }}
              className={`px-4 py-2 rounded-lg font-semibold transition-colors ${
                mode === "typing" ? "bg-primary text-primary-foreground" : "bg-card text-foreground"
              }`}
            >
              Type
            </button>
          </div>
        </div>

        {/* Progress */}
        <div className="mb-8">
          <div className="flex justify-between text-sm text-muted-foreground mb-2">
            <span>
              Question {currentIndex + 1} of {total}
            </span>
            <span>{correctAnswers} correct</span>
          </div>
          <div className="h-3 bg-card border border-border rounded-full overflow-hidden">
            <motion.div
              className="h-full bg-gradient-to-r from-primary to-accent"
              initial={{ width: 0 }}
              animate={{ width: `${progress}%` }}
              transition={{ duration: 0.3 }}
            />
          </div>
        </div>

        {/* Content */}
        <AnimatePresence mode="wait">
          {/* FLASHCARD MODE */}
          {mode === "flashcard" && (
            <motion.div
              key="flashcard"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="space-y-6"
            >
              <div className="perspective-1000">
                <motion.div
                  className="relative w-full h-96 cursor-pointer"
                  onClick={handleFlip}
                  animate={{ rotateY: isFlipped ? 180 : 0 }}
                  transition={{ duration: 0.6 }}
                  style={{ transformStyle: "preserve-3d" }}
                >
                  {/* Front */}
                  <div
                    className="absolute inset-0 bg-card border-2 border-primary/30 rounded-2xl p-8
                               shadow-2xl flex flex-col items-center justify-center backface-hidden"
                    style={{ backfaceVisibility: "hidden" }}
                  >
                    <p className="text-sm text-muted-foreground mb-4 uppercase tracking-wide">Question</p>
                    <h2 className="text-2xl md:text-3xl font-bold text-primary mb-8 text-center">
                      {current?.question}
                    </h2>
                    <p className="text-muted-foreground flex items-center gap-2">
                      <RotateCcw className="w-4 h-4" />
                      Click to reveal
                    </p>
                  </div>

                  {/* Back */}
                  <div
                    className="absolute inset-0 bg-gradient-to-br from-accent/20 to-primary/20
                               border-2 border-primary/50 rounded-2xl p-8 shadow-2xl
                               flex flex-col items-center justify-center backface-hidden"
                    style={{ backfaceVisibility: "hidden", transform: "rotateY(180deg)" }}
                  >
                    <p className="text-sm text-muted-foreground mb-4 uppercase tracking-wide">Answer</p>
                    <h2 className="text-4xl md:text-5xl font-bold text-primary text-center">
                      {current?.answer}
                    </h2>
                  </div>
                </motion.div>
              </div>

              {isFlipped && (
                <motion.div initial={{ opacity: 0, y: 18 }} animate={{ opacity: 1, y: 0 }} className="flex gap-4 justify-center">
                  <button
                    onClick={handleDontKnow}
                    className="px-8 py-4 bg-card border-2 border-destructive/50 text-destructive
                               rounded-lg font-semibold hover:bg-destructive/10 transition-all flex items-center gap-2"
                  >
                    <X className="w-5 h-5" />
                    Don't Know
                  </button>
                  <button
                    onClick={handleKnown}
                    className="px-8 py-4 bg-gradient-to-r from-secondary to-secondary/80 text-secondary-foreground
                               rounded-lg font-semibold hover:shadow-lg transition-all flex items-center gap-2"
                  >
                    <Check className="w-5 h-5" />
                    I Know This
                  </button>
                </motion.div>
              )}
            </motion.div>
          )}

          {/* MULTIPLE CHOICE MODE */}
          {mode === "multiple-choice" && (
            <motion.div
              key="multiple-choice"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="space-y-6"
            >
              <div className="bg-card border border-border rounded-2xl p-8 shadow-xl text-center">
                <p className="text-muted-foreground mb-4">Choose the correct answer:</p>
                <h2 className="text-2xl md:text-3xl font-bold text-primary">{current?.question}</h2>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {(Array.isArray(current?.options) ? current.options : []).map((opt, index) => (
                  <motion.button
                    key={`${opt}-${index}`}
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: index * 0.08 }}
                    onClick={() => !showFeedback && handleMultipleChoice(opt)}
                    disabled={showFeedback}
                    className={`p-6 rounded-xl font-semibold text-lg transition-all ${
                      showFeedback
                        ? opt === current?.answer
                          ? "bg-secondary text-secondary-foreground border-2 border-secondary"
                          : "bg-card text-muted-foreground border border-border opacity-50"
                        : "bg-card text-foreground border-2 border-border hover:border-primary hover:shadow-lg"
                    }`}
                  >
                    {opt}
                  </motion.button>
                ))}
              </div>

              <AnimatePresence>
                {showFeedback && (
                  <motion.div
                    initial={{ opacity: 0, y: 16 }}
                    animate={{
                      opacity: 1,
                      y: 0,
                      scale: isCorrect ? [1, 1.05, 1] : [1, 0.95, 1],
                    }}
                    exit={{ opacity: 0 }}
                    transition={{ scale: { duration: 0.25 }, opacity: { duration: 0.2 } }}
                    className={`p-6 rounded-xl text-center shadow-lg ${
                      isCorrect ? "bg-secondary/20 border-2 border-secondary" : "bg-destructive/20 border-2 border-destructive"
                    }`}
                  >
                    <div className="flex items-center justify-center gap-2 mb-2">
                      {isCorrect ? (
                        <>
                          <motion.div initial={{ scale: 0 }} animate={{ scale: 1 }} transition={{ type: "spring", stiffness: 200 }}>
                            <Check className="w-6 h-6 text-secondary" />
                          </motion.div>
                          <p className="text-xl font-bold text-secondary">Correct!</p>
                        </>
                      ) : (
                        <>
                          <motion.div animate={{ x: [-10, 10, -10, 10, 0] }} transition={{ duration: 0.35 }}>
                            <X className="w-6 h-6 text-destructive" />
                          </motion.div>
                          <p className="text-xl font-bold text-destructive">Not quite</p>
                        </>
                      )}
                    </div>
                    <p className="text-sm text-foreground opacity-80">Answer: {current?.answer}</p>
                  </motion.div>
                )}
              </AnimatePresence>
            </motion.div>
          )}

          {/* TYPING MODE */}
          {mode === "typing" && (
            <motion.div
              key="typing"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="space-y-6"
            >
              <div className="bg-card border border-border rounded-2xl p-8 shadow-xl text-center">
                <p className="text-muted-foreground mb-4">Type the correct answer:</p>
                <h2 className="text-2xl md:text-3xl font-bold text-primary mb-6">{current?.question}</h2>

                <input
                  type="text"
                  value={userAnswer}
                  onChange={(e) => setUserAnswer(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && !showFeedback && handleCheckTyping()}
                  disabled={showFeedback}
                  placeholder="Type your answer..."
                  className="w-full max-w-md mx-auto p-4 text-center text-xl bg-input-background
                             border-2 border-border rounded-lg focus:outline-none focus:ring-2
                             focus:ring-primary/50 disabled:opacity-50"
                  autoFocus
                />
              </div>

              {!showFeedback && (
                <div className="text-center">
                  <button
                    onClick={handleCheckTyping}
                    disabled={!userAnswer.trim()}
                    className="px-8 py-3 bg-gradient-to-r from-primary to-accent text-primary-foreground
                               rounded-lg font-semibold hover:shadow-lg transition-all
                               disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 mx-auto"
                  >
                    Check Answer
                    <ChevronRight className="w-5 h-5" />
                  </button>
                </div>
              )}

              <AnimatePresence>
                {showFeedback && (
                  <motion.div
                    initial={{ opacity: 0, y: 16 }}
                    animate={{
                      opacity: 1,
                      y: 0,
                      scale: isCorrect ? [1, 1.05, 1] : [1, 0.95, 1],
                    }}
                    exit={{ opacity: 0 }}
                    transition={{ scale: { duration: 0.25 }, opacity: { duration: 0.2 } }}
                    className={`p-6 rounded-xl ${
                      isCorrect ? "bg-secondary/20 border-2 border-secondary" : "bg-destructive/20 border-2 border-destructive"
                    }`}
                  >
                    <div className="flex items-center justify-center gap-2 mb-2">
                      {isCorrect ? (
                        <>
                          <motion.div initial={{ scale: 0 }} animate={{ scale: 1 }} transition={{ type: "spring", stiffness: 200 }}>
                            <Check className="w-6 h-6 text-secondary" />
                          </motion.div>
                          <p className="text-xl font-bold text-secondary">Perfect!</p>
                        </>
                      ) : (
                        <>
                          <motion.div animate={{ x: [-10, 10, -10, 10, 0] }} transition={{ duration: 0.35 }}>
                            <X className="w-6 h-6 text-destructive" />
                          </motion.div>
                          <p className="text-xl font-bold text-destructive">Correct: {current?.answer}</p>
                        </>
                      )}
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}
