"use server";

import { UserData } from "@/lib/api";

import UserAvatar from "./UserAvatar";
import UserIdCard from "./UserIdCard";

async function AnotherUserCard({ user }: { user: UserData }) {
  return (
    <div
      className="surface-container outline-variant border w-full max-w-2xl min-w-fit mx-auto p-8"
      style={{ borderRadius: "25px" }}
    >
      <div className="flex items-start gap-5">
        <UserAvatar>{user.name.charAt(0).toUpperCase()}</UserAvatar>

        <div className="flex-1">
          <UserIdCard id={user.id} />

          {/* Name */}
          <div className="mt-4">
            <label className="on-surface mb-1 block text-sm">Name</label>
            <div className="surface-container-high on-surface-variant focus-surface-container-highest w-full rounded-xl px-4 py-2 text-white outline-none">
              {user.name}
            </div>
          </div>

          {/* Description */}
          <div className="mt-4">
            <label className="on-surface mb-1 block text-sm">Description</label>
            <textarea
              disabled
              value={user.description}
              rows={4}
              className="surface-container-high on-surface-variant focus-surface-container-highest w-full rounded-xl px-4 py-2 text-white outline-none"
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default AnotherUserCard;
