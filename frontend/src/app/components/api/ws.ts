export type SessionWsMessage = {
  session_id?: string;
  task_id?: string;
  stage?: string;
  status?: string;
  error?: string;
  words?: any[];
  [k: string]: any;
};

const WS_ORIGIN = "ws://localhost:8080";

function toWsUrl(baseHttpUrl: string) {
  if (baseHttpUrl.startsWith("https://")) return "wss://" + baseHttpUrl.slice("https://".length);
  if (baseHttpUrl.startsWith("http://")) return "ws://" + baseHttpUrl.slice("http://".length);
  if (baseHttpUrl.startsWith("ws://") || baseHttpUrl.startsWith("wss://")) return baseHttpUrl;
  return "ws://" + baseHttpUrl;
}

function joinPath(base: string, path: string) {
  const b = base.replace(/\/+$/, "");
  const p = path.startsWith("/") ? path : `/${path}`;
  return `${b}${p}`;
}

export function openSessionWs(
    session_id: string,
    onMessage: (msg: SessionWsMessage) => void
) {
  const origin = WS_ORIGIN;

  const wsPath = ((import.meta as any)?.env?.VITE_WS_PATH ?? "/ws") as string;

  const wsBase = toWsUrl(origin);
  const baseWithPath =
      wsPath && wsPath !== "/" ? joinPath(wsBase, wsPath) : wsBase.replace(/\/+$/, "");

  const url = `${baseWithPath}?session_id=${encodeURIComponent(session_id)}`;

  const ws = new WebSocket(url);

  ws.onmessage = (e) => {
    try {
      onMessage(JSON.parse(e.data));
    } catch {
      onMessage({ raw: e.data } as any);
    }
  };

  ws.onerror = (e) => console.error("WS error", e);

  ws.onclose = (e) =>
      console.warn("WS closed", { code: (e as any).code, reason: (e as any).reason });

  return ws;
}