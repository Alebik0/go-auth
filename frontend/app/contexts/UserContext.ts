"use client";

import type { UserData } from "@/lib/api";
import { create } from "zustand";

type UserState = {
  user: UserData | null;
  loading: boolean;
};

type UserStore = {
  state: UserState;
  setUser: (user: UserData | null) => void;
};

export const useUserStore = create<UserStore>((set) => ({
  state: { user: null, loading: true },
  setUser: (user) => set({ state: { user: user, loading: false } }),
}));
