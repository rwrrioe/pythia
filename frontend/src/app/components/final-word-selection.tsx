import * as React from "react";
import type { SessionWord } from "./session-types";

type Props = {
  words: SessionWord[];
  wordLimit: number;
  onConfirm: (words: SessionWord[]) => void;
  onBack: () => void;
};

export function FinalWordSelection({ words, wordLimit, onConfirm, onBack }: Props) {
  const [selectedIds, setSelectedIds] = React.useState<Set<string>>(() => {
    const init = new Set<string>();
    for (let i = 0; i < words.length; i++) {
      init.add(String(words[i].id ?? `${words[i].word}-${i}`));
      if (init.size >= wordLimit) break;
    }
    return init;
  });

  const selectedCount = selectedIds.size;

  const toggle = (id: string) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else {
        if (next.size >= wordLimit) return next;
        next.add(id);
      }
      return next;
    });
  };

  const confirm = () => {
    const out: SessionWord[] = [];
    words.forEach((w, idx) => {
      const id = String(w.id ?? `${w.word}-${idx}`);
      if (selectedIds.has(id)) out.push(w);
    });
    onConfirm(out);
  };

  return (
    <div className="min-h-screen bg-background p-4 md:p-8">
      <div className="max-w-3xl mx-auto">
        <div className="flex items-center justify-between gap-3">
          <div>
            <h1 className="text-2xl font-semibold text-foreground">Select final words</h1>
            <p className="text-sm text-muted-foreground">
              Pick up to {wordLimit}. Selected: {selectedCount}
            </p>
          </div>

          <div className="flex gap-2">
            <button className="px-4 py-2 rounded-lg border border-border bg-card hover:bg-muted" onClick={onBack}>
              Back
            </button>
            <button
              className="px-4 py-2 rounded-lg bg-secondary text-secondary-foreground hover:opacity-95 disabled:opacity-50"
              onClick={confirm}
              disabled={selectedCount === 0}
            >
              Confirm
            </button>
          </div>
        </div>

        <div className="mt-6 space-y-2">
          {words.map((w, idx) => {
            const id = String(w.id ?? `${w.word}-${idx}`);
            const checked = selectedIds.has(id);
            const disabled = !checked && selectedCount >= wordLimit;

            return (
              <label key={id} className="flex items-start gap-3 p-3 rounded-lg border border-border bg-card">
                <input
                  type="checkbox"
                  checked={checked}
                  disabled={disabled}
                  onChange={() => toggle(id)}
                  className="mt-1"
                />
                <div>
                  <div className="font-medium text-foreground">{w.word}</div>
                  <div className="text-sm text-muted-foreground">{w.translation}</div>
                  {w.example ? <div className="text-xs text-muted-foreground mt-1">{w.example}</div> : null}
                </div>
              </label>
            );
          })}
        </div>
      </div>
    </div>
  );
}
