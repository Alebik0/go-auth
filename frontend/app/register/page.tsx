"use client";

import { useEffect, useState } from "react";
import { usersApi } from "@/lib/api";
import { useRouter } from "next/navigation";
import LoginForm from "@/components/LoginForm";

enum LoginState {
  LOADING,
  LOGOUT,
}

function RegisterPage() {
  const router = useRouter();
  const [loginState, setLoginState] = useState(LoginState.LOADING);

  useEffect(() => {
    usersApi
      .getMe()
      .then(() => router.push(`/user/my`))
      .catch(() => setLoginState(LoginState.LOGOUT));
  }, [router]);

  return (
    <>
      {loginState == LoginState.LOADING && <p>Loading...</p>}
      {loginState == LoginState.LOGOUT && <p>Loading...</p>}
      <LoginForm />
    </>
  );
}

export default RegisterPage;
