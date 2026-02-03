import React, { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import FlashcardsFull from "./ui/FlashCardsFull";
import Quiz from "./ui/Quiz";

export default function FlashcardsPage() {
  const location = useLocation();
  const navigate = useNavigate();
  // cards должны передаваться через state при навигации: navigate('/cards', { state: { cards, metrics, taskId }})
  const { cards = [], metrics = {}, taskId } = location.state || {};

  const [showQuiz, setShowQuiz] = useState(false);
  const [quizPayload, setQuizPayload] = useState(null);
  const [wsStatus, setWsStatus] = useState("closed");
  const [wsMessages, setWsMessages] = useState([]);

  useEffect(() => {
    // если пользователь попал на страницу напрямую без данных — вернуться на /
    if (!cards || cards.length === 0) {
      // можете менять поведение — я верну на главную
      // navigate("/");
    }
  }, [cards]);

  // открываем ws, если taskId есть — слушаем глобальный WS для карточек/quiz уведомлений
  useEffect(() => {
    if (!taskId) return;
    const socket = new WebSocket(`ws://localhost:8080/ws?task_id=${taskId}`);
    setWsStatus("connecting");
    socket.onopen = () => setWsStatus("open");
    socket.onclose = () => setWsStatus("closed");
    socket.onerror = () => setWsStatus("error");
    socket.onmessage = (evt) => {
      try {
        const payload = JSON.parse(evt.data);
        setWsMessages((s) => [...s, payload]);
        // можете обработать специфичные stage: quiz / flashcards
      } catch (e) {
        console.error("WS error parse", e);
      }
    };
    return () => socket.close();
  }, [taskId]);

  // Подготовим payload для теста из имеющихся cards.
  // Пример состава quiz: из вашего примера payload должен содержать массив quiz объектов.
  function buildQuizFromCards(cards) {
    // Простая генерация: вопрос — перевод (rus), options — mix of tokens
    const tokens = cards.map(c => c.token);
    const quiz = cards.map(c => {
      // собрать 3 случайных неправильных опции
      const wrong = tokens
        .filter(t => t !== c.token)
        .sort(() => Math.random() - 0.5)
        .slice(0, 3);
      const options = [c.token, ...wrong].sort(() => Math.random() - 0.5);
      return {
        answer: c.token,
        question: c.translation || c.translation || "—",
        options
      };
    });
    return { quiz, stage: "quiz", status: "done" };
  }

  async function startQuiz() {
    if (!cards || cards.length === 0) {
      alert("Нет карточек для квиза");
      return;
    }
    const payload = buildQuizFromCards(cards);
    setQuizPayload(payload);

    try {
      const res = await fetch("http://localhost:8080/api/flashcards/tests/quiz", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const t = await res.text();
        console.warn("quiz POST failed:", t);
      }
    } catch (e) {
      console.error("Quiz POST error:", e);
    } finally {
      // открываем локальный квиз UI
      setShowQuiz(true);
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b">
        <div className="max-w-6xl mx-auto px-6 py-4 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold">Flashcards — Pythia</h1>
            <div className="text-sm text-gray-500">WS: {wsStatus} • {taskId ? `task ${taskId}` : "no task"}</div>
          </div>
          <div className="flex items-center gap-3">
            <button onClick={() => navigate(-1)} className="px-3 py-2 border rounded">Назад</button>
            <button onClick={startQuiz} className="px-4 py-2 bg-indigo-600 text-white rounded shadow">Start Quiz</button>
          </div>
        </div>
      </header>

      <main className="max-w-6xl mx-auto px-6 py-8">
        <FlashcardsFull cards={cards} />
      </main>

      {showQuiz && quizPayload && (
        <Quiz
          payload={quizPayload}
          onClose={() => setShowQuiz(false)}
          // можно передать wsMessages/ taskId если хотим слушать бек события в тесте
          taskId={taskId}
        />
      )}
    </div>
  );
}
