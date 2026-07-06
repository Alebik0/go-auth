"use client";

import { useState } from "react";
import { usersApi } from "@/lib/api";
import { useRouter } from "next/navigation";
import { useUserStore } from "@/app/contexts/UserContext";

type SubmitFunction = (login: string, password: string) => Promise<void>;

function AuthForm({ onSubmit }: { onSubmit: SubmitFunction }) {
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

    onSubmit(login, password)
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
            setSubmitError("There is no user found.");
            break;
          case 409:
            setSubmitError("Login is already taken.");
            break;
          default:
            setSubmitError("Internal error, try later.");
            break;
        }
      });
  }

  return (
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
  );
}

export default AuthForm;
