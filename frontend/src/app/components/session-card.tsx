import * as React from 'react';
import { motion } from 'motion/react';
import { BookOpen, Calendar, CheckCircle2, Clock, CalendarDays } from 'lucide-react';
import type { SessionRecord } from './session-types';

interface SessionCardProps {
  session: SessionRecord;
  onOpen: (id: string) => void;
}

function statusLabel(status: SessionRecord['status']) {
  switch (status) {
    case 'completed':
      return 'Completed';
    case 'in-progress':
      return 'In progress';
    case 'scheduled':
      return 'Scheduled';
    default:
      return 'Session';
  }
}

function StatusIcon({ status }: { status: SessionRecord['status'] }) {
  if (status === 'completed') return <CheckCircle2 className="w-4 h-4 text-secondary" />;
  if (status === 'in-progress') return <Clock className="w-4 h-4 text-primary" />;
  return <CalendarDays className="w-4 h-4 text-muted-foreground" />;
}

export function SessionCard({ session, onOpen }: SessionCardProps) {
  const dateLabel = new Date(session.dateISO).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });

  return (
    <motion.button
      type="button"
      onClick={() => onOpen(session.id)}
      whileHover={{ y: -2 }}
      whileTap={{ scale: 0.99 }}
      className="w-full text-left bg-card border border-border rounded-lg shadow-sm hover:shadow-md transition-shadow overflow-hidden"
      aria-label={`Open session ${session.title}`}
    >
      {/* Olive / gold accent line */}
      <div className="h-full w-1 bg-gradient-to-b from-primary/70 to-secondary/60 float-left" />

      <div className="p-4 pl-5">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2 mb-1">
              <StatusIcon status={session.status} />
              <p className="text-xs text-muted-foreground">{statusLabel(session.status)}</p>
            </div>

            <p className="font-semibold text-foreground truncate">{session.title}</p>

            <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-muted-foreground">
              <span className="inline-flex items-center gap-1">
                <Calendar className="w-3.5 h-3.5" />
                {dateLabel}
              </span>
              <span className="inline-flex items-center gap-1">
                <BookOpen className="w-3.5 h-3.5" />
                {session.wordsCount} / {session.wordLimit}
              </span>
              <span className="px-2 py-0.5 rounded bg-primary/10 text-primary">
                {session.language}
              </span>
            </div>
          </div>

          <div className="text-right">
            {typeof session.accuracy === 'number' ? (
              <>
                <p className="text-lg font-semibold text-primary">{session.accuracy}%</p>
                <p className="text-xs text-muted-foreground">accuracy</p>
              </>
            ) : (
              <>
                <p className="text-lg font-semibold text-muted-foreground">—</p>
                <p className="text-xs text-muted-foreground">accuracy</p>
              </>
            )}
          </div>
        </div>
      </div>
    </motion.button>
  );
}
