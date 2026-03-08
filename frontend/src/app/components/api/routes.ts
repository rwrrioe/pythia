// src/app/components/api/routes.ts
export const routes = {
  // auth
  login: "/auth/login",
  register: "/auth/register",

  // dashboard
  dashboard: "/dashboard",

  // sessions (ВАЖНО: session, не sessions)
  createSession: "/session/new",

  uploadFile: (sessionId: string ) => `/session/${sessionId}/upload`,

  finalizeWords: (sessionId: string) => `/session/${sessionId}/summary`,

  getFlashcards: (sessionId: string) => `/session/${sessionId}/learn/flashcards`,

  getQuiz: (sessionId: string ) => `/session/${sessionId}/learn/quiz`,

  translateTask: (sessionId: string, taskId: string) =>
    `/session/${sessionId}/task/${taskId}/translate`,

  endSession: (sessionId: string ) => `/session/${sessionId}/end`,


  librarySessions: "/library/session",
  librarySession: (sessionId: string ) => `/library/session/${sessionId}`,
};

