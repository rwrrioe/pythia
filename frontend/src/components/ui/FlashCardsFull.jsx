// src/ui/FlashcardsFull.jsx
import React from "react";

export default function FlashcardsFull({ cards = [] }) {
  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        {cards.map((c, idx) => (
          <div key={c.id || idx} className="bg-white rounded-xl shadow p-6">
            <div className="text-lg font-semibold">{c.token}</div>
            <div className="text-sm text-gray-600 mt-2">{c.translation}</div>
            {c.examples && c.examples.length > 0 && (
              <div className="text-xs text-gray-500 mt-3">Пример: {c.examples[0]}</div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
