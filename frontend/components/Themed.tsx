"use server";

import { cookies } from "next/headers";
import Header from "./Header";

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
      <Header theme={theme} />
      <main className="surface min-h-screen flex items-center justify-center px-4">
        {children}
      </main>
    </div>
  );
}

export default Themed;
