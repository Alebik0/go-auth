"use client";

import { useEffect, useState } from "react";
import { UserData, usersApi } from "@/lib/api";
import { useRouter } from "next/navigation";

import UserCardSkeleton from "@/components/UserCardSkeleton";
import UserCard from "@/components/UserCard";

function ProfilePage() {
  const router = useRouter();
  const [user, setUser] = useState<UserData | null>(null);

  useEffect(() => {
    usersApi
      .getMe()
      .then((response) => setUser(response.data as UserData))
      .catch((error) => {
        switch (error.response?.status) {
          case 401:
            // Not logined
            router.push("/register");
            break;
          default:
            // Internal error
            console.log("Internal server error", error);
            router.push("/error/500");
            break;
        }
      });
  }, [router]);

  return (
    <>
      {!user && <UserCardSkeleton />}
      {user && <UserCard user={user} />}
    </>
  );
}

export default ProfilePage;
