import React, { useState } from "react";
import Onboarding from "./components/Onboarding";
import Panel from "./components/Panel";

/*
 App: хранит метрики (онбординг), потом отображет панель с кнопкой.
*/

export default function App() {
  const [metrics, setMetrics] = useState(null);

  return (
    <div className="min-h-screen bg-gradient-to-b from-indigo-50 to-white">
      <header className="max-w-4xl mx-auto p-6">
        <h1 className="text-2xl font-bold text-slate-800">Pythia — The Oracle of Language</h1>
        <p className="text-sm text-slate-500 mt-1"> Загрузите фото текста - автоматически определим незнакомые вам слова, переведем и сделаем флеш-карточки</p>
      </header>

      <main className="max-w-4xl mx-auto p-6">
        {!metrics ? (
          <Onboarding onStart={(m) => setMetrics(m)} />
        ) : (
          <Panel metrics={metrics} />
        )}
      </main>
    </div>
  );
}
