"use client";

import Link from "next/link";
import AuthForm from "./AuthForm";

import { authApi } from "@/lib/api";

function RegisterForm() {
  function onSubmit(login: string, password: string) {
    return authApi.register({
      login: login,
      password: password,
    });
  }

  return (
    <div className="surface outline w-full max-w-md rounded-[25px] border p-8">
      <AuthForm onSubmit={onSubmit} />

      <div className="mt-5 text-center">
        <Link
          href="/login"
          className="mx-auto on-primary-container font-medium transition hover-on-secondary-container"
        >
          I already have an account.
        </Link>
      </div>
    </div>
  );
}

export default RegisterForm;
