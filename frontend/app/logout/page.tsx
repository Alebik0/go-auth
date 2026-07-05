"use client";

import { useEffect } from "react";
import { authApi } from "@/lib/api";
import { useRouter } from "next/navigation";
import { useUserStore } from "../contexts/UserContext";

function LogoutPage() {
  const router = useRouter();
  const setUser = useUserStore((s) => s.setUser);

  useEffect(() => {
    authApi
      .logout()
      .then(() => {
        setUser(null);
        router.push("/register");
      })
      .catch(() => {
        setUser(null);
        router.push("/register");
      });
  }, [router, setUser]);

  return <></>;
}

export default LogoutPage;
