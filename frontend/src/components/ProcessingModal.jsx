// ProcessingModal.jsx
import React, { useEffect, useState } from "react";
import Flashcards from "./FlashcardsPage";

export default function ProcessingModal({ metrics, taskId, onClose }) {
  const [ws, setWs] = useState(null);
  const [wsStatus, setWsStatus] = useState("closed");
  const [stage, setStage] = useState(null);
  const [messages, setMessages] = useState([]);
  const [finished, setFinished] = useState(false);
  const [waitingExamples, setWaitingExamples] = useState(false);

  const [showCards, setShowCards] = useState(false);
  const [cardsData, setCardsData] = useState([]);

  const pushMessage = (msg) => setMessages((m) => [...m, msg]);

  useEffect(() => {
    if (!taskId) return;
    const socket = new WebSocket(`ws://localhost:8080/ws?task_id=${taskId}`);
    setWs(socket);
    setWsStatus("connecting");

    socket.onopen = () => setWsStatus("open");
    socket.onclose = () => setWsStatus("closed");
    socket.onerror = () => setWsStatus("error");

    socket.onmessage = (evt) => {
      try {
        const payload = JSON.parse(evt.data);
        if (payload.stage) setStage(payload.stage);

        if (payload.words && payload.stage === "translate") {
          pushMessage({ type: "words", words: payload.words });
          setWaitingExamples(true);
        }

        if ((payload.words && payload.stage === "writing examples") || (payload.examples && payload.stage === "writing examples")) {
          const examplesArr = payload.words ?? payload.examples;

          setMessages(prev =>
            prev.map(msg => {
              if (msg.type !== "words") return msg;
              return {
                ...msg,
                words: msg.words.map(w => {
                  const exObj = examplesArr.find(e => e.word === w.word);
                  return exObj ? { ...w, example: exObj.example } : w;
                })
              };
            })
          );

          setWaitingExamples(false);
        }

        if (payload.status === "done") {
          if (payload.stage === "ocr") startTranslate();
          else if (payload.stage === "translate") startExamples();
          else if (payload.stage === "writing examples") setFinished(true);
        }

      } catch (e) {
        console.error("WS parse error", e);
      }
    };

    return () => {
      try { socket.close(); } catch {}
    };
  }, [taskId]);

  const startTranslate = async () => {
    if (!taskId) return;
    pushMessage({ type: "system", text: "Запрос translate отправлен" });
    try {
      const body = { ...metrics, taskID: taskId, task_id: taskId };
      await fetch("http://localhost:8080/api/translate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });
    } catch (e) {
      pushMessage({ type: "system", text: "Ошибка вызова translate" });
    }
  };

  const startExamples = async () => {
    if (!taskId) return;
    pushMessage({ type: "system", text: "Запрос translate/examples отправлен" });
    try {
      const body = { ...metrics, taskID: taskId, task_id: taskId };
      await fetch("http://localhost:8080/api/translate/examples", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });
    } catch (e) {
      pushMessage({ type: "system", text: "Ошибка вызова translate/examples" });
    }
  };

  const loadFlashcards = async () => {
    try {
      const res = await fetch("http://localhost:8080/api/flashcards", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ analyzerequest: metrics.analyzerequest }),
      });
      const data = await res.json();
      setCardsData(data.flashcards || []);
      navigate('/flashcards', { state: { cards: cardsData, metrics, taskId } });
    } catch (e) {
      console.error("Ошибка загрузки карточек:", e);
      alert("Не удалось загрузить карточки");
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
      <div className="w-full max-w-3xl bg-white rounded-xl shadow-lg overflow-auto max-h-[90vh]">
        <div className="p-4 border-b flex justify-between items-center">
          <div>
            <h3 className="text-lg font-semibold">Обработка страницы</h3>
            <div className="text-xs text-slate-500">WS: {wsStatus} • Stage: {stage || "—"}</div>
          </div>
          <button onClick={onClose} className="text-slate-600 hover:text-slate-900">Закрыть ✕</button>
        </div>

        <div className="p-6 space-y-4">
          {stage && !finished && (
            <div className="rounded-lg p-4 bg-slate-50 border flex items-center gap-3">
              <div className="w-10 h-10 rounded-full border flex items-center justify-center">
                <div className="animate-spin">⏳</div>
              </div>
              <div>
                <div className="font-medium">Обработка: {stage}</div>
                <div className="text-sm text-slate-500">Ищем незнакомые слова</div>
              </div>
            </div>
          )}

          <div className="space-y-4">
            {messages.filter(m => m.type === "words").map((m, idx) => (
              <div key={"w"+idx} className="bg-white border rounded p-4">
                <div className="font-semibold mb-2">Найденные слова</div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  {m.words.map((w, i) => (
                    <div key={i} className="border rounded p-3">
                      <div className="font-medium">{w.word}</div>
                      <div className="text-sm text-slate-600">{w.translation}</div>
                      {w.example && (
                        <div className="text-xs text-slate-400 mt-2">• {w.example}</div>
                      )}
                      {!w.example && waitingExamples && (
                        <div className="text-xs text-slate-400 mt-2">Загружаем пример...</div>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>

          {finished && (
            <div className="flex justify-center mt-6">
              <button
                onClick={loadFlashcards}
                className="px-8 py-4 bg-indigo-600 text-white text-lg font-semibold rounded hover:bg-indigo-700"
              >
                Показать карточки
              </button>
            </div>
          )}

          {showCards && <Flashcards cards={cardsData} onClose={() => setShowCards(false)} />}
        </div>
      </div>
    </div>
  );
}
