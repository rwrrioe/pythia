// src/ui/Quiz.jsx
import React, { useEffect, useMemo, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

export default function Quiz({ payload, onClose, taskId }) {
  // payload.quiz = [{ answer, question, options: [] }, ...]
  const quiz = payload?.quiz || [];
  const [index, setIndex] = useState(0);
  const [selected, setSelected] = useState(null);
  const [results, setResults] = useState([]); // {correct: bool, selected}
  const [showResult, setShowResult] = useState(false);

  useEffect(() => {
    setIndex(0); setSelected(null); setResults([]); setShowResult(false);
  }, [payload]);

  const current = quiz[index];

  function choose(option) {
    if (selected !== null) return; // блокируем повторный выбор
    setSelected(option);
    const correct = option === current.answer;
    setResults(r => [...r, { selected: option, correct }]);
    setShowResult(true);

    // отправим (по желанию) локально событие в WS / UI — если нужен callback
    // например: fetch ws or POST result to server

    // дальше автоматически перейти через паузу
    setTimeout(() => {
      setShowResult(false);
      setSelected(null);
      if (index < quiz.length - 1) setIndex(i => i + 1);
      else {
        // финал — можно отправить результаты на бек
        // POST /api/flashcards/tests/quiz/result (опционально)
      }
    }, 900);
  }

  const score = results.filter(r => r.correct).length;
  const percent = Math.round((score / results.length || 0) * 100);

  return (
    <div className="fixed inset-0 z-60 bg-black/90 flex items-center justify-center p-6">
      <div className="w-full max-w-3xl bg-white rounded-xl overflow-hidden shadow-2xl">
        <div className="p-6 border-b flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold">Quiz</h2>
            <div className="text-sm text-gray-500">{index+1}/{quiz.length}</div>
          </div>
          <div className="flex items-center gap-3">
            <div className="text-sm text-gray-600">{results.length > 0 ? `${score}/${results.length}` : ""}</div>
            <button onClick={onClose} className="px-3 py-2 rounded border">Закрыть</button>
          </div>
        </div>

        <div className="p-8">
          <AnimatePresence mode="wait">
            <motion.div
              key={index}
              initial={{ opacity: 0, x: 40 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -40 }}
              transition={{ duration: 0.35 }}
              className="min-h-[200px] flex flex-col items-center justify-center"
            >
              <div className="text-xl text-gray-700 mb-6 text-center font-medium">{current.question}</div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 w-full max-w-xl">
                {current.options.map((opt, i) => {
                  const isSelected = selected === opt;
                  const correct = selected !== null ? opt === current.answer : false;
                  const wrongSelected = isSelected && selected !== current.answer;

                  return (
                    <motion.button
                      key={i}
                      onClick={() => choose(opt)}
                      whileHover={{ scale: 1.02 }}
                      whileTap={{ scale: 0.98 }}
                      className={`p-4 rounded-lg border text-left font-medium text-lg ${
                        selected === null ? "bg-white" : correct ? "bg-emerald-500 text-white border-emerald-500" : wrongSelected ? "bg-rose-500 text-white border-rose-500" : "opacity-60 bg-white"
                      }`}
                    >
                      {opt}
                    </motion.button>
                  );
                })}
              </div>
            </motion.div>
          </AnimatePresence>
        </div>

        <div className="p-4 border-t flex items-center justify-between">
          <div className="text-sm text-gray-600">Score: {score} / {results.length}</div>
          {index === quiz.length - 1 && results.length === quiz.length && (
            <div className="text-sm font-semibold text-indigo-600">Final: {percent}%</div>
          )}
        </div>
      </div>
    </div>
  );
}
