export type SessionStatus = 'completed' | 'in-progress' | 'scheduled';

export type SessionWord = {
  id?: string | number;
  word: string;
  translation: string;
  example?: string;
  known?: boolean;
};

export interface SessionAttempt {
  dateISO: string; // YYYY-MM-DD
  accuracy: number; // 0-100
}

export interface SessionRecord {
  id: string;
  title: string; // auto-generated from source text or date
  dateISO: string; // YYYY-MM-DD
  language: string;
  difficulty: string;
  wordLimit: number;
  wordsCount: number; // selected words
  status: SessionStatus;
  accuracy?: number; // for completed sessions
  masteredCount?: number;
  words?: SessionWord[];
  attempts?: SessionAttempt[];
}
