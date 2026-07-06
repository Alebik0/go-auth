"use client";

import Link from "next/link";

import { useState } from "react";
import { authApi, usersApi } from "@/lib/api";
import { useRouter } from "next/navigation";
import { useUserStore } from "@/app/contexts/UserContext";

function LoginForm() {
  const router = useRouter();
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [loginError, setLoginError] = useState<string | null>(null);
  const [passwordError, setPasswordError] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const setUser = useUserStore((s) => s.setUser);

  function clearErrors() {
    setLoginError(null);
    setPasswordError(null);
    setSubmitError(null);
  }

  function validate() {
    let success = true;

    if (!login) {
      setLoginError("Login is required");
      success = false;
    } else {
      if (login.length > 32) {
        setLoginError("Maximum length is 32 characters");
        success = false;
      }

      if (!/^[A-Za-z0-9]+$/.test(login)) {
        setLoginError("Only letters and digits are allowed");
        success = false;
      }
    }

    if (!password) {
      setPasswordError("Password is required");
      success = false;
    } else if (password.length > 64) {
      setPasswordError("Maximum length is 64 characters");
      success = false;
    }

    return success;
  }

  function handleSubmit(event: React.SubmitEvent) {
    event.preventDefault();

    if (!validate()) return;

    authApi
      .login({
        login: login,
        password: password,
      })
      .then(() => {
        usersApi
          .getMe()
          .then((response) => setUser(response.data))
          .catch(() => setUser(null));
        router.push("/users/my");
      })
      .catch((error) => {
        switch (error.response?.status) {
          case 404:
            // Login is already taken
            setSubmitError("There is no user found.");
            break;
          default:
            // Internal error
            setSubmitError("Internal error, try later.");
            break;
        }
      });
  }

  return (
    <div className="surface outline w-full max-w-md rounded-[25px] border p-8">
      <form onSubmit={handleSubmit}>
        <h1 className="on-surface mb-8 text-center text-3xl font-bold">
          Sign In
        </h1>

        <div className="mb-6">
          <label
            htmlFor="login"
            className="on-surface mb-2 block text-sm font-medium"
          >
            Login
          </label>

          <input
            id="login"
            type="text"
            value={login}
            maxLength={32}
            autoComplete="username"
            onChange={(e) => {
              clearErrors();
              setLogin(e.target.value);
            }}
            className="surface-variant on-surface-variant outline border-1 m-[1px] w-full rounded-[25px] px-5 py-3 outline-none transition-colors duration-150 ease-in-out focus-primary-outline focus:border-2 focus:m-0"
            placeholder="Enter your login"
          />

          {loginError && (
            <p className="on-error-container mt-2 text-sm">{loginError}</p>
          )}
        </div>

        <div className="mb-8">
          <label
            htmlFor="password"
            className="on-surface mb-2 block text-sm font-medium"
          >
            Password
          </label>

          <input
            id="password"
            type="password"
            value={password}
            maxLength={64}
            autoComplete="current-password"
            onChange={(e) => {
              clearErrors();
              setPassword(e.target.value);
            }}
            className="surface-variant on-surface-variant border-1 m-[1px] w-full rounded-[25px] px-5 py-3 outline-none transition-colors duration-150 ease-in-out focus-primary-outline focus:border-2 focus:m-0"
            placeholder="Enter your password"
          />

          {passwordError && (
            <p className="on-error-container mt-2 text-sm">{passwordError}</p>
          )}
        </div>

        <button
          type="submit"
          className="primary-container on-primary-container w-full rounded-[25px] px-5 py-3 font-semibold text-black transition-all duration-150 ease-in-out hover-secondary-container hover-on-secondary-container"
        >
          Sign In
        </button>
        {submitError && (
          <p className="on-error-container mt-2 text-sm">{submitError}</p>
        )}
      </form>

      <div className="mt-5 text-center">
        <Link
          href="/register"
          className="mx-auto on-primary-container font-medium transition hover-on-secondary-container"
        >
          Create new account.
        </Link>
      </div>
    </div>
  );
}

export default LoginForm;
