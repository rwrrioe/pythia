// src/main.jsx
import React from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import App from "./App";
import FlashcardsPage from "./components/FlashcardsPage";
import "./index.css"; // tailwind entry

createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/flashcards" element={<FlashcardsPage />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);
