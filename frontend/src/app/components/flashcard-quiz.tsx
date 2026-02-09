import * as React from 'react';
import { useState } from 'react';
import { motion, AnimatePresence } from 'motion/react';
import { ArrowLeft, RotateCcw, Check, X, ChevronRight, Trophy } from 'lucide-react';
import { GreekPattern } from './greek-pattern';

interface Word {
  id: string;
  word: string;
  translation: string;
  example: string;
  known: boolean;
}

interface FlashcardQuizProps {
  words: Word[];
  onComplete: (score: number) => void;
  onBack: () => void;
}

type QuizMode = 'flashcard' | 'multiple-choice' | 'typing';

export function FlashcardQuiz({ words, onComplete, onBack }: FlashcardQuizProps) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [mode, setMode] = useState<QuizMode>('flashcard');
  const [isFlipped, setIsFlipped] = useState(false);
  const [correctAnswers, setCorrectAnswers] = useState(0);
  const [userAnswer, setUserAnswer] = useState('');
  const [showFeedback, setShowFeedback] = useState(false);
  const [isCorrect, setIsCorrect] = useState(false);

  const currentWord = words[currentIndex];
  const progress = ((currentIndex + 1) / words.length) * 100;

  const handleFlip = () => {
    setIsFlipped(!isFlipped);
  };

  const handleKnown = () => {
    setCorrectAnswers(prev => prev + 1);
    handleNext();
  };

  const handleNext = () => {
    if (currentIndex < words.length - 1) {
      setCurrentIndex(prev => prev + 1);
      setIsFlipped(false);
      setUserAnswer('');
      setShowFeedback(false);
    } else {
      const score = Math.round((correctAnswers / words.length) * 100);
      onComplete(score);
    }
  };

  const handleCheckAnswer = () => {
    const correct = userAnswer.toLowerCase().trim() === currentWord.translation.toLowerCase().trim();
    setIsCorrect(correct);
    setShowFeedback(true);
    if (correct) {
      setCorrectAnswers(prev => prev + 1);
    }
    setTimeout(() => handleNext(), 1500);
  };

  const generateMultipleChoiceOptions = () => {
    const options = [currentWord.translation];
    const otherWords = words.filter(w => w.id !== currentWord.id);
    while (options.length < 4 && otherWords.length > 0) {
      const randomIndex = Math.floor(Math.random() * otherWords.length);
      const option = otherWords[randomIndex].translation;
      if (!options.includes(option)) {
        options.push(option);
      }
      otherWords.splice(randomIndex, 1);
    }
    return options.sort(() => Math.random() - 0.5);
  };

  const handleMultipleChoice = (answer: string) => {
    const correct = answer === currentWord.translation;
    setIsCorrect(correct);
    setShowFeedback(true);
    if (correct) {
      setCorrectAnswers(prev => prev + 1);
    }
    setTimeout(() => handleNext(), 1500);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background p-4 md:p-8 relative overflow-hidden">
      {/* Background pattern */}
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
              onClick={() => setMode('flashcard')}
              className={`px-4 py-2 rounded-lg font-semibold transition-colors ${
                mode === 'flashcard' ? 'bg-primary text-primary-foreground' : 'bg-card text-foreground'
              }`}
            >
              Flashcards
            </button>
            <button
              onClick={() => setMode('multiple-choice')}
              className={`px-4 py-2 rounded-lg font-semibold transition-colors ${
                mode === 'multiple-choice' ? 'bg-primary text-primary-foreground' : 'bg-card text-foreground'
              }`}
            >
              Quiz
            </button>
            <button
              onClick={() => setMode('typing')}
              className={`px-4 py-2 rounded-lg font-semibold transition-colors ${
                mode === 'typing' ? 'bg-primary text-primary-foreground' : 'bg-card text-foreground'
              }`}
            >
              Type
            </button>
          </div>
        </div>

        {/* Progress Bar */}
        <div className="mb-8">
          <div className="flex justify-between text-sm text-muted-foreground mb-2">
            <span>Word {currentIndex + 1} of {words.length}</span>
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
          {mode === 'flashcard' && (
            <motion.div
              key="flashcard"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="space-y-6"
            >
              {/* Flashcard */}
              <div className="perspective-1000">
                <motion.div
                  className="relative w-full h-96 cursor-pointer"
                  onClick={handleFlip}
                  animate={{ rotateY: isFlipped ? 180 : 0 }}
                  transition={{ duration: 0.6 }}
                  style={{ transformStyle: 'preserve-3d' }}
                >
                  {/* Front */}
                  <div
                    className="absolute inset-0 bg-card border-2 border-primary/30 rounded-2xl p-8 
                               shadow-2xl flex flex-col items-center justify-center backface-hidden"
                    style={{ backfaceVisibility: 'hidden' }}
                  >
                    <p className="text-sm text-muted-foreground mb-4 uppercase tracking-wide">Word</p>
                    <h2 className="text-5xl md:text-6xl font-bold text-primary mb-8" 
                        style={{ fontFamily: 'var(--font-heading)' }}>
                      {currentWord.word}
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
                    style={{ 
                      backfaceVisibility: 'hidden',
                      transform: 'rotateY(180deg)'
                    }}
                  >
                    <p className="text-sm text-muted-foreground mb-4 uppercase tracking-wide">Translation</p>
                    <h2 className="text-5xl md:text-6xl font-bold text-primary mb-6"
                        style={{ fontFamily: 'var(--font-heading)' }}>
                      {currentWord.translation}
                    </h2>
                    <div className="bg-card/80 rounded-lg p-4 max-w-md">
                      <p className="text-sm text-foreground italic">{currentWord.example}</p>
                    </div>
                  </div>
                </motion.div>
              </div>

              {/* Actions */}
              {isFlipped && (
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="flex gap-4 justify-center"
                >
                  <button
                    onClick={handleNext}
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

          {mode === 'multiple-choice' && (
            <motion.div
              key="multiple-choice"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="space-y-6"
            >
              {/* Question */}
              <div className="bg-card border border-border rounded-2xl p-8 shadow-xl text-center">
                <p className="text-muted-foreground mb-4">What is the translation of:</p>
                <h2 className="text-4xl md:text-5xl font-bold text-primary" 
                    style={{ fontFamily: 'var(--font-heading)' }}>
                  {currentWord.word}
                </h2>
              </div>

              {/* Options */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {generateMultipleChoiceOptions().map((option, index) => (
                  <motion.button
                    key={index}
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: index * 0.1 }}
                    onClick={() => !showFeedback && handleMultipleChoice(option)}
                    disabled={showFeedback}
                    className={`p-6 rounded-xl font-semibold text-lg transition-all ${
                      showFeedback
                        ? option === currentWord.translation
                          ? 'bg-secondary text-secondary-foreground border-2 border-secondary'
                          : 'bg-card text-muted-foreground border border-border opacity-50'
                        : 'bg-card text-foreground border-2 border-border hover:border-primary hover:shadow-lg'
                    }`}
                  >
                    {option}
                  </motion.button>
                ))}
              </div>

              {/* Feedback */}
              <AnimatePresence>
                {showFeedback && (
                  <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ 
                      opacity: 1, 
                      y: 0,
                      scale: isCorrect ? [1, 1.05, 1] : [1, 0.95, 1]
                    }}
                    exit={{ opacity: 0 }}
                    transition={{ 
                      scale: { duration: 0.3 },
                      opacity: { duration: 0.2 }
                    }}
                    className={`p-6 rounded-xl text-center shadow-lg ${
                      isCorrect 
                        ? 'bg-secondary/20 border-2 border-secondary' 
                        : 'bg-destructive/20 border-2 border-destructive'
                    }`}
                  >
                    <div className="flex items-center justify-center gap-2 mb-2">
                      {isCorrect ? (
                        <>
                          <motion.div
                            initial={{ scale: 0 }}
                            animate={{ scale: 1 }}
                            transition={{ type: "spring", stiffness: 200 }}
                          >
                            <Check className="w-6 h-6 text-secondary" />
                          </motion.div>
                          <p className="text-xl font-bold text-secondary">Correct!</p>
                        </>
                      ) : (
                        <>
                          <motion.div
                            animate={{ x: [-10, 10, -10, 10, 0] }}
                            transition={{ duration: 0.4 }}
                          >
                            <X className="w-6 h-6 text-destructive" />
                          </motion.div>
                          <p className="text-xl font-bold text-destructive">Not quite</p>
                        </>
                      )}
                    </div>
                    <p className="text-sm text-foreground italic">{currentWord.example}</p>
                  </motion.div>
                )}
              </AnimatePresence>
            </motion.div>
          )}

          {mode === 'typing' && (
            <motion.div
              key="typing"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="space-y-6"
            >
              {/* Question */}
              <div className="bg-card border border-border rounded-2xl p-8 shadow-xl text-center">
                <p className="text-muted-foreground mb-4">Type the translation:</p>
                <h2 className="text-4xl md:text-5xl font-bold text-primary mb-6" 
                    style={{ fontFamily: 'var(--font-heading)' }}>
                  {currentWord.word}
                </h2>
                <input
                  type="text"
                  value={userAnswer}
                  onChange={(e) => setUserAnswer(e.target.value)}
                  onKeyPress={(e) => e.key === 'Enter' && !showFeedback && handleCheckAnswer()}
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
                    onClick={handleCheckAnswer}
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

              {/* Feedback */}
              <AnimatePresence>
                {showFeedback && (
                  <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ 
                      opacity: 1, 
                      y: 0,
                      scale: isCorrect ? [1, 1.05, 1] : [1, 0.95, 1]
                    }}
                    exit={{ opacity: 0 }}
                    transition={{ 
                      scale: { duration: 0.3 },
                      opacity: { duration: 0.2 }
                    }}
                    className={`p-6 rounded-xl ${
                      isCorrect ? 'bg-secondary/20 border-2 border-secondary' : 'bg-destructive/20 border-2 border-destructive'
                    }`}
                  >
                    <div className="flex items-center justify-center gap-2 mb-2">
                      {isCorrect ? (
                        <>
                          <motion.div
                            initial={{ scale: 0 }}
                            animate={{ scale: 1 }}
                            transition={{ type: "spring", stiffness: 200 }}
                          >
                            <Check className="w-6 h-6 text-secondary" />
                          </motion.div>
                          <p className="text-xl font-bold text-secondary">Perfect!</p>
                        </>
                      ) : (
                        <>
                          <motion.div
                            animate={{ x: [-10, 10, -10, 10, 0] }}
                            transition={{ duration: 0.4 }}
                          >
                            <X className="w-6 h-6 text-destructive" />
                          </motion.div>
                          <p className="text-xl font-bold text-destructive">The correct answer is: {currentWord.translation}</p>
                        </>
                      )}
                    </div>
                    <p className="text-sm text-foreground italic text-center">{currentWord.example}</p>
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