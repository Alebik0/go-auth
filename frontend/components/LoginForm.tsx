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
      className='w-full max-w-md rounded-[25px] bg-[#1A1D24] border border-[#2A2F38] p-8 shadow-2xl'
    >
      <h1 className='mb-8 text-center text-3xl font-bold text-white'>
        Sign In
      </h1>

      <div className='mb-6'>
        <label
          htmlFor='login'
          className='mb-2 block text-sm font-medium text-gray-300'
        >
          Login
        </label>

        <input
          id='login'
          type='text'
          value={login}
          maxLength={32}
          autoComplete='username'
          onChange={(e) => setLogin(e.target.value)}
          className='w-full rounded-[25px] border border-gray-700 bg-[#12151B] px-5 py-3 text-white outline-none transition focus:border-[#A0D49B] focus:ring-2 focus:ring-[#A0D49B]/30'
          placeholder='Enter your login'
        />

        {errors.login && (
          <p className='mt-2 text-sm text-red-400'>{errors.login}</p>
        )}
      </div>

      <div className='mb-8'>
        <label
          htmlFor='password'
          className='mb-2 block text-sm font-medium text-gray-300'
        >
          Password
        </label>

        <input
          id='password'
          type='password'
          value={password}
          maxLength={64}
          autoComplete='current-password'
          onChange={(e) => setPassword(e.target.value)}
          className='w-full rounded-[25px] border border-gray-700 bg-[#12151B] px-5 py-3 text-white outline-none transition focus:border-[#A0D49B] focus:ring-2 focus:ring-[#A0D49B]/30'
          placeholder='Enter your password'
        />

        {errors.password && (
          <p className='mt-2 text-sm text-red-400'>{errors.password}</p>
        )}
      </div>

      <button
        type='submit'
        className='w-full rounded-[25px] bg-[#A0D49B] px-5 py-3 font-semibold text-black transition hover:brightness-110 active:scale-[0.98]'
      >
        Sign In
      </button>
      {errors.submit && (
        <p className='mt-2 text-sm text-red-400'>{errors.submit}</p>
      )}
    </form>
  );
}

export default LoginForm;
