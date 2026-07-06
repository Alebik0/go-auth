"use server";

import axios, { AxiosError } from "axios";
import { UserData } from "@/lib/api";

import Error500Card from "@/components/errors/500/Error500Card";
import Error404Card from "@/components/errors/500/Error404Card";
import AnotherUserCard from "@/components/users/AnotherUserCard";

async function ProfilePage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;

  let result: UserData | null;
  let errorStatus: number | undefined | null;

  try {
    const response = await axios.get(
      `${process.env.NEXT_PUBLIC_API_URL}/api/v1/users/${id}`,
    );
    result = response.data as UserData;
    errorStatus = null;
  } catch (error) {
    result = null;
    errorStatus = (error as AxiosError).response?.status;
  }

  if (errorStatus === 404) {
    return <Error404Card />;
  }

  if (errorStatus) {
    return <Error500Card />;
  }

  return <>{result != null && <AnotherUserCard user={result} />}</>;
}

export default ProfilePage;
