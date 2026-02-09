// src/app/components/api/sessions.ts
import { apiFetch } from "./http";
import { routes } from "./routes";
import type { SessionWord } from "../session-types";

/** Strict language types (no more `any`) */
export type LangId = 1 | 2 | 3 | 4;
export type LangCode = "en" | "de" | "fr" | "es";

export function langIdToCode(id: LangId): LangCode {
  switch (id) {
    case 1: return "en";
    case 2: return "de";
    case 3: return "fr";
    case 4: return "es";
  }
}

export function ensureLangId(x: unknown): LangId {
  const n = typeof x === "number" ? x : Number(x);
  if (n === 1 || n === 2 || n === 3 || n === 4) return n;
  // default safe language
  return 2; // de
}

/** Payload you send to backend create session.
 *  IMPORTANT: backend expects LangId as int (1..4). */
export type CreateSessionRequest = {
  language: LangId;
  difficulty: string;
  wordLimit: number;
};

export type AnalyzeRequest = {
  Level: string;      // "B1"
  Durating: string;   // backend field is Durating (typo kept)
  Lang: string;       // "de"
};

export type QuizQuestion = {
  answer: string;
  question: string;
  options: string[];
};

export type UploadOcrAccepted = {
  session_id: number;
  stage: "ocr" | string;
  task_id: string;
};

export type TranslateAccepted = {
  session_id: number;
  stage: "translate" | string;
  task_id: string;
};


export async function createSession(payload: {
  language: number;        // 1..4
  wordLimit: number;       // slider
  durationSeconds: number; // 1800
}): Promise<{ session_id: number }> {
  const body = {
    lang_id: Number(payload.language),
    words_count: Number(payload.wordLimit),
    durating: Number(payload.durationSeconds), // ВАЖНО: число, не строка
  };

  const res = await apiFetch<any>(routes.createSession, {
    method: "POST",
    body: JSON.stringify(body),
  });

  const id = Number(res?.session_id ?? res?.id ?? res?.sessionId);
  return { session_id: id };
}

export async function uploadFile(
  sessionId: number,
  file: File,
  lang: LangCode
): Promise<UploadOcrAccepted> {
  const form = new FormData();
  form.append("file", file);
  form.append("lang", lang);

  const res = await apiFetch<any>(routes.uploadFile(sessionId), {
    method: "POST",
    body: form,
  });

  return {
    session_id: Number(res?.session_id ?? sessionId),
    stage: String(res?.stage ?? "ocr"),
    task_id: String(res?.task_id),
  };
}

// IMPORTANT: REST returns 202 {session_id, task_id, stage:"translate"}
// Words will arrive via WebSocket (processing/done)
export async function translateTask(
  sessionId: number,
  taskId: string,
  payload: { duration: string; lang: LangCode; level: string }
): Promise<TranslateAccepted> {
  // map to backend AnalyzeRequest struct:
  // Level, Durating, Lang
  const body: AnalyzeRequest = {
    Level: String(payload.level),
    Durating: String(payload.duration),
    Lang: String(payload.lang),
  };

  const res = await apiFetch<any>(routes.translateTask(sessionId, taskId), {
    method: "POST",
    body: JSON.stringify(body),
  });

  return {
    session_id: Number(res?.session_id ?? sessionId),
    stage: String(res?.stage ?? "translate"),
    task_id: String(res?.task_id ?? taskId),
  };
}

export async function endSession(sessionId: number) {
  // backend route: PATCH /api/session/:id/end
  return apiFetch<any>(routes.endSession(sessionId), { method: "PATCH" });
}

export async function getFlashcards(
  sessionId: number
): Promise<{ words: SessionWord[] }> {
  const res = await apiFetch<any>(routes.getFlashcards(sessionId), {
    method: "GET",
  });

  const list =
    (Array.isArray(res?.words) && res.words) ||
    (Array.isArray(res?.flashcards) && res.flashcards) ||
    (Array.isArray(res) && res) ||
    [];

  const words: SessionWord[] = list.map((x: any, i: number) => ({
    id: x?.id ?? i,
    word: String(x?.word ?? x?.front ?? x?.term ?? ""),
    translation: String(x?.translation ?? x?.back ?? x?.meaning ?? ""),
    example: x?.example ? String(x.example) : undefined,
  }));

  return { words };
}

/** Helper: end session then immediately fetch flashcards (your requirement #2). */
export async function endSessionAndFetchFlashcards(
  sessionId: number
): Promise<{ words: SessionWord[] }> {
  await endSession(sessionId);          // wait for 200
  return await getFlashcards(sessionId); // then fetch flashcards
}

export async function getQuiz(
  sessionId: number
): Promise<{ questions: QuizQuestion[] }> {
  const res = await apiFetch<any>(routes.getQuiz(sessionId), { method: "GET" });
  const questions: QuizQuestion[] = Array.isArray(res?.questions)
    ? res.questions
    : Array.isArray(res)
      ? res
      : Array.isArray(res?.quiz)
        ? res.quiz
        : [];
  return { questions };
}
