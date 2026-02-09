import * as React from 'react';
import { motion } from 'motion/react';
import { Trophy, Target, Calendar, TrendingUp, ChevronRight } from 'lucide-react';
import { GreekPattern, LaurelWreath } from './greek-pattern';
import { PifMascot } from './pif-mascot';

interface Word {
  id: string;
  word: string;
  translation: string;
  example: string;
  known: boolean;
}

interface SessionSummaryProps {
  words: Word[];
  score: number;
  onBackToDashboard: () => void;
}

export function SessionSummary({ words, score, onBackToDashboard }: SessionSummaryProps) {
  const getEncouragement = () => {
    if (score >= 90) return { 
      text: "Excellent! The Oracle is pleased!", 
      color: "text-secondary",
      pifMessage: "Amazing work! You're becoming a master!"
    };
    if (score >= 75) return { 
      text: "Well done, scholar!", 
      color: "text-primary",
      pifMessage: "Great job! Keep up the good work!"
    };
    if (score >= 60) return { 
      text: "Good progress!", 
      color: "text-accent",
      pifMessage: "Nice effort! You're improving!"
    };
    return { 
      text: "Keep practicing!", 
      color: "text-foreground",
      pifMessage: "Don't give up! Practice makes perfect!"
    };
  };

  const encouragement = getEncouragement();
  const nextReviewDate = new Date();
  nextReviewDate.setDate(nextReviewDate.getDate() + 1);

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background p-4 md:p-8 relative overflow-hidden">
      {/* Background pattern */}
      <div className="absolute top-0 left-0 right-0">
        <GreekPattern className="w-full h-6 text-primary opacity-40" />
      </div>

      {/* Decorative laurel wreaths */}
      <div className="absolute top-20 left-10 opacity-10">
        <LaurelWreath className="w-48 h-48 text-primary" />
      </div>
      <div className="absolute bottom-20 right-10 opacity-10">
        <LaurelWreath className="w-48 h-48 text-secondary" />
      </div>

      <div className="max-w-4xl mx-auto relative z-10 pt-8">
        {/* Trophy Section */}
        <motion.div
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.5 }}
          className="text-center mb-8"
        >
          <div className="inline-flex items-center justify-center w-32 h-32 bg-gradient-to-br from-primary to-accent rounded-full mb-6 shadow-2xl">
            <Trophy className="w-16 h-16 text-primary-foreground" />
          </div>
          <h1 className="text-4xl md:text-5xl font-bold text-primary mb-3" style={{ fontFamily: 'var(--font-heading)' }}>
            Session Complete!
          </h1>
          <p className={`text-2xl ${encouragement.color} font-semibold mb-6`}>
            {encouragement.text}
          </p>
          
          {/* Pif Mascot with encouragement */}
          <div className="flex justify-center">
            <PifMascot 
              message={encouragement.pifMessage}
              variant={score >= 75 ? "happy" : "encouraging"}
              size="md"
            />
          </div>
        </motion.div>

        {/* Score Card */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="bg-gradient-to-br from-primary/10 to-accent/10 border-2 border-primary/30 rounded-2xl p-8 mb-8 shadow-xl text-center"
        >
          <p className="text-muted-foreground mb-2">Your Score</p>
          <p className="text-7xl font-bold text-primary mb-2" style={{ fontFamily: 'var(--font-heading)' }}>
            {score}%
          </p>
          <p className="text-foreground">
            <span className="font-semibold text-2xl">{words.length}</span> words studied
          </p>
        </motion.div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="bg-card border border-border rounded-xl p-6 text-center shadow-lg"
          >
            <div className="inline-flex items-center justify-center w-12 h-12 bg-secondary/10 rounded-full mb-3">
              <Target className="w-6 h-6 text-secondary" />
            </div>
            <p className="text-3xl font-bold text-foreground mb-1">{Math.round(score / 10)}</p>
            <p className="text-sm text-muted-foreground">Accuracy Level</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.4 }}
            className="bg-card border border-border rounded-xl p-6 text-center shadow-lg"
          >
            <div className="inline-flex items-center justify-center w-12 h-12 bg-primary/10 rounded-full mb-3">
              <Calendar className="w-6 h-6 text-primary" />
            </div>
            <p className="text-3xl font-bold text-foreground mb-1">
              {nextReviewDate.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
            </p>
            <p className="text-sm text-muted-foreground">Next Review</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.5 }}
            className="bg-card border border-border rounded-xl p-6 text-center shadow-lg"
          >
            <div className="inline-flex items-center justify-center w-12 h-12 bg-accent/10 rounded-full mb-3">
              <TrendingUp className="w-6 h-6 text-accent" />
            </div>
            <p className="text-3xl font-bold text-foreground mb-1">+{words.length}</p>
            <p className="text-sm text-muted-foreground">Total Learned</p>
          </motion.div>
        </div>

        {/* Words Review */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6 }}
          className="bg-card border border-border rounded-xl p-6 mb-8 shadow-xl"
        >
          <h2 className="text-2xl font-bold text-foreground mb-4" style={{ fontFamily: 'var(--font-heading)' }}>
            Words from This Session
          </h2>
          <div className="space-y-3">
            {words.map((word, index) => (
              <motion.div
                key={word.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 0.7 + index * 0.05 }}
                className="flex items-center justify-between p-4 bg-muted/30 rounded-lg border border-border/50 hover:bg-muted/50 transition-colors"
              >
                <div>
                  <p className="font-semibold text-foreground">{word.word}</p>
                  <p className="text-sm text-muted-foreground">{word.translation}</p>
                </div>
                <div className="w-2 h-2 bg-primary rounded-full" />
              </motion.div>
            ))}
          </div>
        </motion.div>

        {/* Review Schedule Info */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.8 }}
          className="bg-secondary/10 border border-secondary/30 rounded-xl p-6 mb-8 text-center"
        >
          <h3 className="text-lg font-semibold text-foreground mb-2">
            Spaced Repetition Schedule
          </h3>
          <p className="text-muted-foreground text-sm mb-4">
            Review these words tomorrow, then in 3 days, and finally in 7 days for optimal retention
          </p>
          <div className="flex justify-center gap-2">
            <div className="px-3 py-1 bg-secondary/20 rounded-full text-sm text-secondary font-semibold">
              Day 1
            </div>
            <div className="px-3 py-1 bg-secondary/20 rounded-full text-sm text-secondary font-semibold">
              Day 3
            </div>
            <div className="px-3 py-1 bg-secondary/20 rounded-full text-sm text-secondary font-semibold">
              Day 7
            </div>
          </div>
        </motion.div>

        {/* Action Button */}
        <motion.button
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.9 }}
          onClick={onBackToDashboard}
          className="w-full py-4 bg-gradient-to-r from-primary to-accent text-primary-foreground 
                     rounded-xl font-semibold text-lg shadow-lg hover:shadow-xl 
                     transform hover:scale-105 transition-all duration-300 
                     border border-primary/20 relative overflow-hidden group flex items-center justify-center gap-2"
        >
          <span className="relative z-10">Return to Sanctuary</span>
          <ChevronRight className="w-5 h-5 relative z-10" />
          <div className="absolute inset-0 bg-gradient-to-r from-accent to-primary opacity-0 
                          group-hover:opacity-100 transition-opacity duration-300" />
        </motion.button>
      </div>
    </div>
  );
}