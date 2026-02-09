// src/app/components/api/library.ts
import { apiFetch } from "./http";
import { routes } from "./routes";

export type LibrarySessionDto = {
  Id: number;
  name?: string;
  StartedAt?: string;
  ended_at?: string;
  Duration?: number;
  Status?: string;
  language?: number; // 1..4
  level?: number;
  Accuracy?: number;
};

export type LibrarySessionsResponse = {
  sessions: LibrarySessionDto[];
};

export type LibraryFlashcardDto = {
  word: string;
  translation: string;
  language: string; // "en" in response example
};

export type LibrarySessionDetailResponse = {
  flashcards: LibraryFlashcardDto[];
  session: LibrarySessionDto;
};

export async function getLibrarySessions(): Promise<LibrarySessionsResponse> {
  return apiFetch<LibrarySessionsResponse>(routes.librarySessions, { method: "GET" });
}

export async function getLibrarySession(sessionId: number | string): Promise<LibrarySessionDetailResponse> {
  return apiFetch<LibrarySessionDetailResponse>(routes.librarySession(sessionId), { method: "GET" });
}
