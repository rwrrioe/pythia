import * as React from 'react';
import { useState } from 'react';
import { motion } from 'motion/react';
import { Upload, FileText, ArrowLeft, Check } from 'lucide-react';
import { GreekPattern } from './greek-pattern';
import { PifMascot } from './pif-mascot';

interface Word {
  id: string;
  word: string;
  translation: string;
  example: string;
  known: boolean;
}

interface TextUploadProps {
  onWordsSelected: (words: Word[]) => void;
  onBack: () => void;
}

export function TextUpload({ onWordsSelected, onBack }: TextUploadProps) {
  const [text, setText] = useState('');
  const [extractedWords, setExtractedWords] = useState<Word[]>([]);
  const [selectedWords, setSelectedWords] = useState<Set<string>>(new Set());

  const handleExtractWords = () => {
    // Simulate word extraction
    const mockWords: Word[] = [
      { id: '1', word: 'Philosophy', translation: 'Filosofía', example: 'Ancient Greek philosophy shaped Western thought.', known: false },
      { id: '2', word: 'Democracy', translation: 'Democracia', example: 'Democracy originated in ancient Athens.', known: false },
      { id: '3', word: 'Oracle', translation: 'Oráculo', example: 'The Oracle of Delphi was famous throughout Greece.', known: false },
      { id: '4', word: 'Temple', translation: 'Templo', example: 'The Parthenon was a temple dedicated to Athena.', known: false },
      { id: '5', word: 'Mythology', translation: 'Mitología', example: 'Greek mythology influenced many cultures.', known: false },
      { id: '6', word: 'Acropolis', translation: 'Acrópolis', example: 'The Acropolis overlooks the city of Athens.', known: false },
      { id: '7', word: 'Theater', translation: 'Teatro', example: 'Ancient Greek theater was performed in amphitheaters.', known: false },
      { id: '8', word: 'Symposium', translation: 'Simposio', example: 'A symposium was a drinking party in ancient Greece.', known: false },
      { id: '9', word: 'Gymnasium', translation: 'Gimnasio', example: 'The gymnasium was a place for physical and intellectual training.', known: false },
      { id: '10', word: 'Laurel', translation: 'Laurel', example: 'Victors were crowned with laurel wreaths.', known: false },
    ];
    setExtractedWords(mockWords);
    setSelectedWords(new Set(mockWords.map(w => w.id)));
  };

  const toggleWord = (id: string) => {
    const newSelected = new Set(selectedWords);
    if (newSelected.has(id)) {
      newSelected.delete(id);
    } else {
      newSelected.add(id);
    }
    setSelectedWords(newSelected);
  };

  const handleContinue = () => {
    const words = extractedWords.filter(w => selectedWords.has(w.id));
    onWordsSelected(words);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background p-4 md:p-8 relative">
      <div className="absolute top-0 left-0 right-0">
        <GreekPattern className="w-full h-6 text-primary opacity-40" />
      </div>

      <div className="max-w-4xl mx-auto relative z-10 pt-8">
        {/* Header */}
        <div className="flex items-center gap-4 mb-8">
          <button onClick={onBack} className="p-2 rounded-lg hover:bg-muted transition-colors">
            <ArrowLeft className="w-6 h-6 text-foreground" />
          </button>
          <div>
            <h1 className="text-3xl text-foreground mb-1">Upload Your Text</h1>
            <p className="text-muted-foreground">Paste text or upload a document to extract words</p>
          </div>
        </div>

        {extractedWords.length === 0 ? (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="space-y-6"
          >
            {/* Pif Mascot */}
            <div className="flex justify-center mb-6">
              <PifMascot 
                message="Hi! Upload some text and I'll help you find the best words to learn!" 
                variant="happy"
                size="lg"
              />
            </div>

            {/* Upload Options */}
            <div className="bg-card border-2 border-dashed border-primary/30 rounded-xl p-12 text-center">
              <div className="mb-6">
                <Upload className="w-16 h-16 text-primary mx-auto mb-4" />
                <h2 className="text-xl font-semibold text-foreground mb-2">Upload Image or Document</h2>
                <p className="text-muted-foreground">JPG, PNG, or PDF • OCR will extract text</p>
              </div>
              <input 
                type="file" 
                className="hidden" 
                id="file-upload"
                accept="image/*,.pdf"
              />
              <label 
                htmlFor="file-upload"
                className="inline-block px-8 py-3 bg-primary text-primary-foreground rounded-lg font-semibold 
                           hover:bg-primary/90 transition-colors cursor-pointer"
              >
                Choose File
              </label>
            </div>

            {/* Or Divider */}
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-border"></div>
              </div>
              <div className="relative flex justify-center">
                <span className="px-4 bg-background text-muted-foreground">or</span>
              </div>
            </div>

            {/* Text Input */}
            <div className="bg-card border border-border rounded-xl p-6">
              <div className="flex items-center gap-3 mb-4">
                <FileText className="w-6 h-6 text-primary" />
                <h2 className="text-xl font-semibold text-foreground">Paste Text Directly</h2>
              </div>
              <textarea
                value={text}
                onChange={(e) => setText(e.target.value)}
                placeholder="Paste your text here... The oracle will find the words for you to learn."
                className="w-full h-48 p-4 bg-input-background border border-border rounded-lg 
                           focus:outline-none focus:ring-2 focus:ring-primary/50 resize-none"
              />
              <button
                onClick={handleExtractWords}
                disabled={!text.trim()}
                className="mt-4 w-full py-3 bg-gradient-to-r from-primary to-accent text-primary-foreground 
                           rounded-lg font-semibold hover:shadow-lg transition-all duration-300
                           disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Extract Words
              </button>
            </div>
          </motion.div>
        ) : (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="space-y-6"
          >
            {/* Word Count */}
            <div className="bg-primary/10 border border-primary/30 rounded-lg p-4 flex items-center justify-between">
              <p className="text-foreground">
                <span className="font-bold text-2xl text-primary">{selectedWords.size}</span> words selected
              </p>
              <button
                onClick={handleContinue}
                disabled={selectedWords.size === 0}
                className="px-6 py-2 bg-gradient-to-r from-primary to-accent text-primary-foreground 
                           rounded-lg font-semibold hover:shadow-lg transition-all
                           disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Continue to Learn
              </button>
            </div>

            {/* Word List */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {extractedWords.map((word, index) => (
                <motion.div
                  key={word.id}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.05 }}
                  onClick={() => toggleWord(word.id)}
                  className={`bg-card border-2 rounded-lg p-4 cursor-pointer transition-all ${
                    selectedWords.has(word.id)
                      ? 'border-primary shadow-lg'
                      : 'border-border hover:border-primary/50'
                  }`}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <h3 className="text-lg font-semibold text-foreground">{word.word}</h3>
                      <p className="text-sm text-muted-foreground">{word.translation}</p>
                    </div>
                    {selectedWords.has(word.id) && (
                      <div className="p-1 bg-primary rounded-full">
                        <Check className="w-4 h-4 text-primary-foreground" />
                      </div>
                    )}
                  </div>
                </motion.div>
              ))}
            </div>
          </motion.div>
        )}
      </div>
    </div>
  );
}