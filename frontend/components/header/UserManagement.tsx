"use client";

import Link from "next/link";
import { useEffect } from "react";
import { usersApi } from "@/lib/api";
import { useUserStore } from "@/app/contexts/UserContext";

function UserManagement() {
  const userState = useUserStore((s) => s.state);
  const setUser = useUserStore((s) => s.setUser);

  useEffect(() => {
    usersApi
      .getMe()
      .then((response) => setUser(response.data))
      .catch(() => setUser(null));
  }, [setUser]);

  if (userState.loading) {
    return <></>;
  }

  if (userState.user == null) {
    return (
      <Link
        href="/login"
        className="surface-variant on-surface-variant outline rounded border px-3 py-2 text-sm"
      >
        Login
      </Link>
    );
  }

  return (
    <>
      <Link
        href="/users/my"
        className="surface-variant on-surface-variant outline rounded border px-3 py-2 text-sm"
      >
        Profile
      </Link>

      <Link
        href="/logout"
        className="surface-variant on-surface-variant outline rounded border px-3 py-2 text-sm"
      >
        Logout
      </Link>
    </>
  );
}

export default UserManagement;
