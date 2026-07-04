// app/api/theme/route.ts

import { NextRequest, NextResponse } from "next/server";

const VALID_THEMES = [
  "light",
  "light-high-contrast",
  "light-medium-contrast",
  "dark",
  "dark-high-contrast",
  "dark-medium-contrast",
] as const;

type Theme = (typeof VALID_THEMES)[number];

export async function POST(request: NextRequest) {
  const { theme } = await request.json();

  if (!VALID_THEMES.includes(theme)) {
    return NextResponse.json({ error: "Invalid theme" }, { status: 400 });
  }

  const response = NextResponse.json({ success: true });

  response.cookies.set({
    name: "theme",
    value: theme as Theme,
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
    maxAge: 60 * 60 * 24 * 365, // 1 year
  });

  return response;
}
