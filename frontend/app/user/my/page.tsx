"use client";

import { useEffect, useState } from "react";
import { UserData, usersApi } from "@/lib/api";
import { useRouter } from "next/navigation";
import UserCard from "@/components/UserCard";

function ProfilePage() {
  const router = useRouter();
  const [user, setUser] = useState<UserData | null>(null);

  useEffect(() => {
    usersApi
      .getMe()
      .then((response) => setUser(response.data as UserData))
      .catch((error) => console.log(error));
  }, [router]);

  return (
    <>
      {user && <UserCard user={user} />}
      {!user && <p>Loading...</p>}
    </>
  );
}

export default ProfilePage;
