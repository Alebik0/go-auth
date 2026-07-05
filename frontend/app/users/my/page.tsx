"use client";

import { useEffect } from "react";
import { usersApi } from "@/lib/api";
import { useRouter } from "next/navigation";
import { useUserStore } from "@/app/contexts/UserContext";

import UserCardSkeleton from "@/components/UserCardSkeleton";
import UserCard from "@/components/UserCard";

function ProfilePage() {
  const router = useRouter();
  const userState = useUserStore((s) => s.state);
  const setUser = useUserStore((s) => s.setUser);

  useEffect(() => {
    usersApi
      .getMe()
      .then((response) => setUser(response.data))
      .catch((error) => {
        switch (error.response?.status) {
          case 401:
            // Not logined
            setUser(null);
            break;
          default:
            // Internal error
            console.log("Internal server error", error);
            router.push("/error/500");
            break;
        }
      });
  }, [router, setUser]);

  if (userState.loading) {
    return <UserCardSkeleton />;
  }

  return <>{userState.user != null && <UserCard user={userState.user} />}</>;
}

export default ProfilePage;
