import { UserData } from "@/lib/api";
import { useState } from "react";

type UserCardProps = {
  user: UserData;
};

function UserCard({ user }: UserCardProps) {
  const [name, setName] = useState(user.name);
  const [description, setDescription] = useState(user.description);

  return (
    <div
      className="surface-container outline-variant border w-full max-w-2xl min-w-fit mx-auto p-8"
      style={{ borderRadius: "25px" }}
    >
      <div className="flex items-start gap-5">
        {/* Avatar */}
        <div className="primary-container flex h-20 w-20 shrink-0 items-center justify-center rounded-full">
          <span className="on-primary-container text-5xl font-bold">
            {user.name.charAt(0).toUpperCase()}
          </span>
        </div>

        <div className="flex-1">
          {/* User ID */}
          <div className="surface-container-high flex items-center justify-between rounded-[18px] px-5 py-3">
            <span className="on-surface-variant text-sm">User ID</span>
            <span className="on-primary-container font-semibold">
              #{user.id}
            </span>
          </div>

          {/* Name */}
          <div className="mt-4">
            <label className="on-surface mb-1 block text-sm">Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="surface-container-high on-surface-variant w-full rounded-xl px-4 py-2 text-white outline-none focus-surface-container-highest"
            />
          </div>

          {/* Description */}
          <div className="mt-4">
            <label className="on-surface mb-1 block text-sm">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={4}
              className="surface-container-high on-surface-variant w-full rounded-xl px-4 py-2 outline-none resize-none focus-surface-container-highest"
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default UserCard;
