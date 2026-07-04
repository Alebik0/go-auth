"use client";

import { useState } from "react";
import { authApi } from "@/lib/api";
import { useRouter } from "next/navigation";

function LoginForm() {
  const router = useRouter();
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");

  const [errors, setErrors] = useState<{
    login?: string;
    password?: string;
    submit?: string;
  }>({});

  const validate = () => {
    const nextErrors: typeof errors = {};

    if (!login) {
      nextErrors.login = "Login is required";
    } else {
      if (login.length > 32) {
        nextErrors.login = "Maximum length is 32 characters";
      }

      if (!/^[A-Za-z0-9]+$/.test(login)) {
        nextErrors.login = "Only letters and digits are allowed";
      }
    }

    if (!password) {
      nextErrors.password = "Password is required";
    } else if (password.length > 64) {
      nextErrors.password = "Maximum length is 64 characters";
    }

    setErrors(nextErrors);

    return Object.keys(nextErrors).length === 0;
  };

  const handleSubmit = (e: React.SubmitEvent) => {
    e.preventDefault();

    if (!validate()) return;

    authApi
      .register({
        login: login,
        password: password,
      })
      .then(() => router.push("/user/my"))
      .catch(() => (errors.submit = "Failed to submit"));
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="surface outline w-full max-w-md rounded-[25px] border p-8"
    >
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
          onChange={(e) => setLogin(e.target.value)}
          className="surface-variant on-surface-variant outline border-1 m-[1px] w-full rounded-[25px] px-5 py-3 outline-none transition-colors duration-150 ease-in-out focus-primary-outline focus:border-2 focus:m-0"
          placeholder="Enter your login"
        />

        {errors.login && (
          <p className="on-error-container mt-2 text-sm">{errors.login}</p>
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
          onChange={(e) => setPassword(e.target.value)}
          className="surface-variant on-surface-variant border-1 m-[1px] w-full rounded-[25px] px-5 py-3 outline-none transition-colors duration-150 ease-in-out focus-primary-outline focus:border-2 focus:m-0"
          placeholder="Enter your password"
        />

        {errors.password && (
          <p className="on-error-container mt-2 text-sm">{errors.password}</p>
        )}
      </div>

      <button
        type="submit"
        className="primary-container on-primary-container w-full rounded-[25px] px-5 py-3 font-semibold text-black transition-all duration-150 ease-in-out hover-secondary-container hover-on-secondary-container"
      >
        Sign In
      </button>
      {errors.submit && (
        <p className="on-error-container mt-2 text-sm">{errors.submit}</p>
      )}
    </form>
  );
}

export default LoginForm;
