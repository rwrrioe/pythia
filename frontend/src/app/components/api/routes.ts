// src/app/components/api/routes.ts
export const routes = {
  // auth
  login: "/auth/login",
  register: "/auth/register",

  // dashboard
  dashboard: "/dashboard",

  // sessions (ВАЖНО: session, не sessions)
  createSession: "/session/new",

  uploadFile: (sessionId: string | number) => `/session/${sessionId}/upload`,

  finalizeWords: (sessionId: string | number) => `/session/${sessionId}/summary`,

  getFlashcards: (sessionId: string | number) => `/session/${sessionId}/learn/flashcards`,

  getQuiz: (sessionId: string | number) => `/session/${sessionId}/learn/quiz`,

  translateTask: (sessionId: string | number, taskId: string) =>
    `/session/${sessionId}/task/${taskId}/translate`,

  endSession: (sessionId: string | number) => `/session/${sessionId}/end`,


  librarySessions: "/library/session",
  librarySession: (sessionId: string | number) => `/library/session/${sessionId}`,
};

