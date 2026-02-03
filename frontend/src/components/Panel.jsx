import React, { useRef, useState } from "react";
import ProcessingModal from "./ProcessingModal";
import { UploadCloud } from "lucide-react";

export default function Panel({ metrics }) {
  const [open, setOpen] = useState(false);
  const [taskId, setTaskId] = useState(null);
  const fileRef = useRef(null);
  const [uploading, setUploading] = useState(false);

  const handleFiles = async (files) => {
    if (!files || files.length === 0) return;

    const file = files[0];
    const formData = new FormData();
    formData.append("task_id", ""); // бэк создаст ID
    formData.append("file", file);

    setUploading(true);

    try {
      const res = await fetch("http://localhost:8080/api/upload", {
        method: "POST",
        body: formData,
      });
      if (!res.ok) throw new Error(res.statusText);

      const data = await res.json().catch(() => ({}));
      setTaskId(data.task_id || null);
      setOpen(true);
    } catch (err) {
      console.error("Upload error:", err);
      alert("Ошибка при загрузке файла: " + err.message);
    } finally {
      setUploading(false);
    }
  };

  const handleDrop = (e) => {
    e.preventDefault();
    handleFiles(e.dataTransfer.files);
  };

  const handleDragOver = (e) => e.preventDefault();

  return (
    <div className="bg-white rounded-xl shadow p-8 flex flex-col items-center gap-6">
      <div className="text-center">
        <h2 className="text-xl font-semibold">
          Готово — загрузите страницу учебника
        </h2>
        <p className="text-sm text-slate-500 mt-2">
          Мы распознаем текст, предложим слова, переводы и примеры в контексте.
        </p>
      </div>

      <div
        className={`w-full max-w-md border-2 border-dashed rounded-xl p-6 flex flex-col items-center justify-center cursor-pointer hover:border-indigo-500 hover:bg-indigo-50 transition-all duration-200 ${
          uploading ? "opacity-70" : ""
        }`}
        onClick={() => fileRef.current.click()}
        onDrop={handleDrop}
        onDragOver={handleDragOver}
      >
        <UploadCloud size={24} className="text-indigo-500 mb-2" />
        {uploading ? (
          <p className="text-indigo-500 font-medium">Загрузка файла...</p>
        ) : (
          <p className="text-gray-700 font-medium text-center">
            Кликни сюда или перетащи фото
          </p>
        )}
        <input
          ref={fileRef}
          type="file"
          accept="image/*,application/pdf"
          className="hidden"
          onChange={(e) => handleFiles(e.target.files)}
        />
      </div>

      {open && (
        <ProcessingModal
          metrics={metrics}
          taskId={taskId}
          onClose={() => setOpen(false)}
        />
      )}
    </div>
  );
}
