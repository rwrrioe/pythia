# Pythia — Frontend (React + Vite + Tailwind)

This repository contains a production-ready frontend for *Pythia — The Oracle of The Language*.
It implements:
- onboarding (collects metrics: lang, level, duration, book)
- chat-like interface (file upload to `/api/upload` and websocket `/ws?task_id=...`)
- buttons to request `/api/translate` and `/api/translate/examples`
- in-chat display of words, context examples, marking know/unknown and export of flashcards JSON
- uses Tailwind for styling and Vite for fast builds

---

## Files
- `index.html` — app entry
- `src/` — React source
  - `App.jsx` — root
  - `components/Onboarding.jsx`
  - `components/ChatPane.jsx`
- `vite.config.js`, `package.json`, `tailwind.config.cjs`, `postcss.config.cjs`

---

## How it communicates with your backend (API contract)

1. **Upload page image**
   - Endpoint: `POST /api/upload`
   - Form data fields:
     - `task_id` (string)
     - `file` (file) — image or PDF
   - Response: `200 OK` or error. Backend will process OCR and emit progress via websocket.

2. **WebSocket**
   - Connect to: `ws://<host>/ws?task_id=<id>` (use `wss://` under HTTPS)
   - Server sends JSON messages. Expected shapes:
     - `{ "status":"processing"|"error"|"done", "stage":"ocr"|"translate"|"writing examples" }`
     - `{ "words": [ { "id","token","guess","translation","examples":[...] }, ... ] }`
     - `{ "examples": [ { "sentence","source" }, ... ] }`

3. **Request translation**
   - Endpoint: `POST /api/translate`
   - Body (JSON): `{ "task_id","level","lang","duration","book" }`
   - Response: may return `{ "words": [...] }` immediately or the server may push via websocket.

4. **Request context examples**
   - Endpoint: `POST /api/translate/examples`
   - Body: same as `/api/translate`
   - Response: `{ "examples": [...] }` or via websocket.

---

## Development

Prerequisites: Node >=18, npm or pnpm.

1. Install:
```bash
cd pythia_frontend
npm install
```

2. Run dev server:
```bash
npm run dev
# opens at http://localhost:5173
```

3. Build for production:
```bash
npm run build
npm run preview
```

## Deployment recommendations (production ready)

- Build static files with `npm run build` and serve with any static hosting (NGINX, Vercel, Netlify).
- If you need a single server, build and copy `/dist` to your backend's static folder.
- Use HTTPS and a reverse proxy so WebSocket upgrades work (configure nginx for `proxy_set_header Upgrade $http_upgrade;`).
- Add authentication on API endpoints (JWT/Cookie) before public launch.
- Rate-limit uploads and protect file endpoints (virus scanning for PDFs).
- For big OCR jobs return early and stream progress via WS or SSE. Persist task state in Redis (you mentioned using Redis).

---

## Notes & ideas to extend
- Add client-side caching of seen words (idb-keyval).
- Export direct Anki deck (.apkg) via backend when user requests export.
- Add spaced repetition algorithm on backend and sync with user's account.
- Allow multi-page upload and batch processing.

---