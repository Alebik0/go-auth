"use server";

import Link from "next/link";
import ThemeButton from "./ThemeButton";

async function Header({ theme }: Readonly<{ theme: string }>) {
  return (
    <header className="surface outline flex items-center justify-between border-b px-6 py-4">
      <h1 className="on-surface text-xl font-semibold">My App</h1>

      <div className="flex items-center gap-4">
        <ThemeButton theme={theme} />

        <Link
          href="/logout"
          className="surface-variant on-surface-variant outline rounded border px-3 py-2 text-sm"
        >
          Logout
        </Link>
      </div>
    </header>
  );
}

export default Header;
