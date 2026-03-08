// src/app/components/api/sessions.ts
import { apiFetch } from "./http";
import { routes } from "./routes";
import type { SessionWord } from "../session-types";

/** Session id is UUID string */
export type SessionId = string;

/** Strict language types */
export type LangId = 1 | 2 | 3 | 4;
export type LangCode = "en" | "de" | "fr" | "es";

export function langIdToCode(id: LangId): LangCode {
  switch (id) {
    case 1:
      return "en";
    case 2:
      return "de";
    case 3:
      return "fr";
    case 4:
      return "es";
  }
}

export function ensureLangId(x: unknown): LangId {
  const n = typeof x === "number" ? x : Number(x);
  if (n === 1 || n === 2 || n === 3 || n === 4) return n;
  return 2; // default safe language: de
}

/** Payload you send to backend create session.
 * IMPORTANT: backend expects lang_id as int (1..4) */
export type CreateSessionRequest = {
  language: LangId;
  difficulty: string;
  wordLimit: number;
};

export type AnalyzeRequest = {
  Level: string; // "B1"
  Durating: string; // backend field is Durating (typo kept)
  Lang: string; // "de"
};

export type QuizQuestion = {
  answer: string;
  question: string;
  options: string[];
};

/** 202 Accepted shapes (REST) */
export type UploadOcrAccepted = {
  session_id: SessionId;
  stage: "ocr" | string;
  task_id: string;
};

export type TranslateAccepted = {
  session_id: SessionId;
  stage: "translate" | string;
  task_id: string;
};

/** Helper: safely extract session id from different backend shapes */
function pickSessionId(res: any, fallback?: SessionId): SessionId {
  const id =
      res?.session_id ??
      res?.sessionId ??
      res?.id ??
      res?.session?.id ??
      fallback;

  if (typeof id !== "string") {
    // if backend accidentally returned uuid.UUID as something non-string,
    // fail fast so it's not silently encoded as [object Object]
    throw new Error(`Invalid session_id in response: ${JSON.stringify(id)}`);
  }
  return id;
}

export async function createSession(payload: {
  language: number; // 1..4
  wordLimit: number; // slider
  durationSeconds: number; // 1800
}): Promise<{ session_id: SessionId }> {
  const body = {
    lang_id: Number(payload.language),
    words_count: Number(payload.wordLimit),
    durating: Number(payload.durationSeconds), // IMPORTANT: number, not string
  };

  const res = await apiFetch<any>(routes.createSession, {
    method: "POST",
    body: JSON.stringify(body),
  });

  return { session_id: pickSessionId(res) };
}

export async function uploadFile(
    sessionId: SessionId,
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
    session_id: pickSessionId(res, sessionId),
    stage: String(res?.stage ?? "ocr"),
    task_id: String(res?.task_id),
  };
}

// IMPORTANT: REST returns 202 {session_id, task_id, stage:"translate"}
// Words will arrive via WebSocket (processing/done)
export async function translateTask(
    sessionId: SessionId,
    taskId: string,
    payload: { duration: string; lang: LangCode; level: string }
): Promise<TranslateAccepted> {
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
    session_id: pickSessionId(res, sessionId),
    stage: String(res?.stage ?? "translate"),
    task_id: String(res?.task_id ?? taskId),
  };
}

export async function endSession(sessionId: SessionId) {
  // backend route: PATCH /api/session/:id/end
  return apiFetch<any>(routes.endSession(sessionId), { method: "PATCH" });
}

export async function getFlashcards(
    sessionId: SessionId
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

/** Helper: end session then immediately fetch flashcards */
export async function endSessionAndFetchFlashcards(
    sessionId: SessionId
): Promise<{ words: SessionWord[] }> {
  await endSession(sessionId);
  return await getFlashcards(sessionId);
}

export async function getQuiz(
    sessionId: SessionId
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