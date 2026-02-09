import * as React from "react";
import { useState } from "react";
import { motion } from "motion/react";
import { ArrowRight, Book, Target, Globe } from "lucide-react";

import type { LangId } from "./api/sessions";

interface SessionSetupProps {
  onContinue: (config: SessionConfig) => void | Promise<void>; // разрешаем async
  onBack: () => void;
}

export interface SessionConfig {
  language: LangId;      // 1..4
  difficulty: string;    // "a2" | "b1" | "b2" etc
  wordLimit: number;     // 10..15
}

export function SessionSetup({ onContinue, onBack }: SessionSetupProps) {
  const [language, setLanguage] = useState<LangId>(2); // default de
  const [difficulty, setDifficulty] = useState("a2");
  const [wordLimit, setWordLimit] = useState(12);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onContinue({ language, difficulty, wordLimit });
  };

  const languages: Array<{ id: LangId; label: string }> = [
    { id: 2, label: "German" },
    { id: 1, label: "English" },
    { id: 3, label: "French" },
    { id: 4, label: "Spanish" },
  ];

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-8"
        >
          <button
            onClick={onBack}
            className="text-muted-foreground hover:text-foreground transition-colors mb-4"
          >
            ← Back to Dashboard
          </button>
          <h1 className="text-3xl text-foreground mb-2">New Learning Session</h1>
          <p className="text-muted-foreground">
            Configure your session parameters. The oracle will select the most important words for you.
          </p>
        </motion.div>

        {/* Form */}
        <motion.form
          onSubmit={handleSubmit}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="bg-card border border-border rounded-xl p-8 shadow-lg"
        >
          {/* Language Selection */}
          <div className="mb-6">
            <label className="flex items-center gap-2 text-foreground mb-3">
              <Globe className="w-5 h-5 text-primary" />
              Learning Language
            </label>
            <div className="grid grid-cols-2 gap-3">
              {languages.map((lang) => (
                <button
                  key={lang.id}
                  type="button"
                  onClick={() => setLanguage(lang.id)}
                  className={`
                    p-4 rounded-lg border-2 transition-all duration-200
                    ${language === lang.id
                      ? "border-primary bg-primary/10 text-primary"
                      : "border-border bg-muted/50 text-foreground hover:border-primary/50"
                    }
                  `}
                >
                  {lang.label}
                </button>
              ))}
            </div>
          </div>

          {/* Difficulty Level */}
          <div className="mb-6">
            <label className="flex items-center gap-2 text-foreground mb-3">
              <Target className="w-5 h-5 text-primary" />
              Difficulty Level
            </label>
            <div className="grid grid-cols-3 gap-3">
              {[
                { value: "a2", label: "A2", desc: "Elementary" },
                { value: "b1", label: "B1", desc: "Intermediate" },
                { value: "b2", label: "B2", desc: "Upper Int." },
              ].map((level) => (
                <button
                  key={level.value}
                  type="button"
                  onClick={() => setDifficulty(level.value)}
                  className={`
                    p-4 rounded-lg border-2 transition-all duration-200 text-center
                    ${difficulty === level.value
                      ? "border-primary bg-primary/10 text-primary"
                      : "border-border bg-muted/50 text-foreground hover:border-primary/50"
                    }
                  `}
                >
                  <div className="font-semibold">{level.label}</div>
                  <div className="text-xs mt-1 opacity-70">{level.desc}</div>
                </button>
              ))}
            </div>
          </div>

          {/* Word Limit */}
          <div className="mb-8">
            <label className="flex items-center gap-2 text-foreground mb-3">
              <Book className="w-5 h-5 text-primary" />
              Words per Session: <span className="text-primary font-semibold">{wordLimit}</span>
            </label>
            <input
              type="range"
              min="10"
              max="15"
              value={wordLimit}
              onChange={(e) => setWordLimit(Number(e.target.value))}
              className="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer accent-primary"
            />
            <div className="flex justify-between text-xs text-muted-foreground mt-2">
              <span>10 words</span>
              <span>15 words</span>
            </div>
          </div>

          {/* Submit Button */}
          <motion.button
            type="submit"
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            className="w-full py-4 bg-gradient-to-r from-primary to-accent text-primary-foreground 
                       rounded-lg font-semibold shadow-lg hover:shadow-xl transition-all duration-300
                       flex items-center justify-center gap-2"
          >
            Continue to Upload
            <ArrowRight className="w-5 h-5" />
          </motion.button>
        </motion.form>

        {/* Info Box */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3 }}
          className="mt-6 p-4 bg-secondary/10 border border-secondary/30 rounded-lg"
        >
          <p className="text-sm text-muted-foreground text-center">
            💡 The oracle selects only the most important words from your text, ensuring focused and efficient learning.
          </p>
        </motion.div>
      </div>
    </div>
  );
}
