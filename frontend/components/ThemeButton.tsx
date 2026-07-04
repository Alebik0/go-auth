"use client";

import axios from "axios";
import { useRouter } from "next/navigation";
import { useState } from "react";

const themes = [
  "light",
  "light-high-contrast",
  "light-medium-contrast",
  "dark",
  "dark-high-contrast",
  "dark-medium-contrast",
] as const;

type Theme = (typeof themes)[number];

function ThemeButton({ theme }: Readonly<{ theme: string }>) {
  const router = useRouter();
  const [isError, setIsError] = useState(false);

  async function onNextTheme() {
    const themeIndex = themes.indexOf(theme as Theme);
    const nextIndex = (themeIndex + 1) % themes.length;
    const nextTheme = themes[nextIndex];

    axios
      .post("/api/theme", { theme: nextTheme }, { timeout: 3000 })
      .then(() => router.refresh())
      .catch(() => setIsError(true));
  }

  const classes = `outline rounded border px-3 py-2 text-sm ${
    isError
      ? "error-container on-error-container"
      : "surface-variant on-surface-variant"
  }`;

  return (
    <button onClick={onNextTheme} className={classes}>
      Theme: {theme}
    </button>
  );
}

export default ThemeButton;
