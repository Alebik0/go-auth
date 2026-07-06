"use server";

import { cookies } from "next/headers";
import Header from "./header/Header";

const themes = [
  "light",
  "light-high-contrast",
  "light-medium-contrast",
  "dark",
  "dark-high-contrast",
  "dark-medium-contrast",
] as const;

type Theme = (typeof themes)[number];

async function Themed({ children }: Readonly<{ children: React.ReactNode }>) {
  const cookieStore = await cookies();
  const cookieValue = cookieStore.get("theme")?.value ?? "light";
  const themeIndex = themes.indexOf(cookieValue as Theme);
  const theme = themes[themeIndex];

  return (
    <div className={theme}>
      <div className="surface min-h-screen min-h-screen flex flex-col justify-between">
        <header className="shrink-0 surface outline flex items-center justify-between border-b px-6 py-4">
          <Header theme={theme} />
        </header>
        <main className="flex-1 min-h-0 overflow-auto flex flex-col items-center justify-center">
          {children}
        </main>
      </div>
    </div>
  );
}

export default Themed;
