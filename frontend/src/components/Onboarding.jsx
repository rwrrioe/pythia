import React from "react";

export default function Onboarding({ onStart }) {
  // defaults
  const [level, setLevel] = React.useState("A1");
  const [lang, setLang] = React.useState("de");
  const [durating, setDurating] = React.useState("30");
  const [book, setBook] = React.useState("");

  const submit = (e) => {
    e.preventDefault();
    // note: send these exact keys later (level, lang, durating, book, taskID)
    onStart({ level, lang, durating, book });
  };

  return (
    <form onSubmit={submit} className="bg-white rounded-xl shadow p-6 space-y-4">
      <h2 className="text-lg font-semibold">Начать новую сессию</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div>
          <label className="block text-sm text-slate-600">Уровень (CEFR)</label>
          <select value={level} onChange={(e)=>setLevel(e.target.value)} className="mt-1 w-full border rounded px-2 py-2">
            <option>A1</option><option>A2</option><option>B1</option><option>B2</option><option>C1</option><option>C2</option>
          </select>
        </div>
        <div>
          <label className="block text-sm text-slate-600">Язык</label>
          <select value={lang} onChange={(e)=>setLang(e.target.value)} className="mt-1 w-full border rounded px-2 py-2">
            <option value="de">Немецкий (de)</option>
            <option value="en">Английский (en)</option>
            <option value="es">Испанский (es)</option>
          </select>
        </div>

        <div>
          <label className="block text-sm text-slate-600">Длительность изучения. Например: 5 месяцев</label>
          <input value={durating} onChange={(e)=>setDurating(e.target.value)} className="mt-1 w-full border rounded px-2 py-2" />
        </div>

        <div>
          <label className="block text-sm text-slate-600">Книга / Источник</label>
          <input value={book} onChange={(e)=>setBook(e.target.value)} className="mt-1 w-full border rounded px-2 py-2" placeholder="Название книги" />
        </div>
      </div>

      <div className="flex justify-end">
        <button type="submit" className="px-4 py-2 bg-indigo-600 text-white rounded shadow">Начать</button>
      </div>
    </form>
  );
}
