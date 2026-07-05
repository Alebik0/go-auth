"use client";

import { useEffect } from "react";
import { authApi } from "@/lib/api";
import { useRouter } from "next/navigation";

function LogoutPage() {
  const router = useRouter();

  useEffect(() => {
    authApi
      .logout()
      .then(() => router.push("/register"))
      .catch(() => router.push("/register"));
  }, [router]);

  return <></>;
}

export default LogoutPage;
