"use server";

import ThemeButton from "../ThemeButton";
import UserManagement from "./UserManagement";

async function Header({ theme }: Readonly<{ theme: string }>) {
  return (
    <>
      <h1 className="on-surface text-xl font-semibold">My App</h1>

      <div className="flex items-center gap-4">
        <ThemeButton theme={theme} />
        <UserManagement />
      </div>
    </>
  );
}

export default Header;
