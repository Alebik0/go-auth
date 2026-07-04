"use client";

import { useEffect, useState } from "react";
import { usersApi } from "@/lib/api";
import { useRouter } from "next/navigation";

enum LoginState {
  LOADING,
  LOGIN,
  LOGOUT,
}

function HomePage() {
  const router = useRouter();
  const [loginState, setLoginState] = useState(LoginState.LOADING);

  useEffect(() => {
    usersApi
      .getMe()
      .then(() => setLoginState(LoginState.LOGIN))
      .catch(() => router.push("/register"));
  }, [router]);

  return (
    <>
      <header>
        <h1>Just a simple header</h1>
      </header>
      <section>
        <h2>Hello world</h2>
        {loginState}
      </section>
    </>
  );
}

export default HomePage;
