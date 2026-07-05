"use client";

import { useEffect, useState } from "react";
import { UserData, usersApi } from "@/lib/api";
import { useRouter } from "next/navigation";

import UserCardSkeleton from "@/components/UserCardSkeleton";
import UserCard from "@/components/UserCard";
import Error500 from "@/components/Error500";

function ProfilePage() {
  const router = useRouter();
  const [user, setUser] = useState<UserData | null>(null);
  const [isInternalError, setIsInternalError] = useState(false);

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
            setIsInternalError(true);
        }
      });
  }, [router]);

  return (
    <>
      {isInternalError && <Error500 />}
      {user && <UserCard user={user} />}
      {!user && <UserCardSkeleton />}
    </>
  );
}

export default ProfilePage;
