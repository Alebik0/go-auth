import { UserData } from "@/lib/api";
import { ChangeEvent, useState } from "react";

type UserCardProps = {
  user: UserData;
};

function UserCard({ user }: UserCardProps) {
  const [name, setName] = useState(user.name);
  const [description, setDescription] = useState(user.description);

  return (
    <main className='min-h-screen bg-[#0F1115] flex items-center justify-center p-6'>
      <div
        className='w-full max-w-md bg-[#181B20] border border-[#2A2F36] p-8 shadow-2xl'
        style={{ borderRadius: "25px" }}
      >
        <div className='flex items-start gap-5'>
          {/* Avatar */}
          <div className='flex h-20 w-20 shrink-0 items-center justify-center rounded-full bg-[#A0D49B]/15 border border-[#A0D49B]/30'>
            <span className='text-3xl font-bold text-[#A0D49B]'>
              {user.name.charAt(0).toUpperCase()}
            </span>
          </div>

          <div className='flex-1'>
            {/* User ID */}
            <div className='flex items-center justify-between rounded-[18px] bg-[#111418] px-5 py-3'>
              <span className='text-sm text-gray-400'>User ID</span>
              <span className='font-semibold text-[#A0D49B]'>#{user.id}</span>
            </div>

            {/* Name */}
            <div className='mt-4'>
              <label className='mb-1 block text-sm text-gray-400'>Name</label>
              <input
                type='text'
                value={name}
                onChange={(e) => setName(e.target.value)}
                className='w-full rounded-xl border border-[#2A2F36] bg-[#111418] px-4 py-2 text-white outline-none focus:border-[#A0D49B]'
              />
            </div>

            {/* Description */}
            <div className='mt-4'>
              <label className='mb-1 block text-sm text-gray-400'>
                Description
              </label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={4}
                className='w-full rounded-xl border border-[#2A2F36] bg-[#111418] px-4 py-2 text-white outline-none resize-none focus:border-[#A0D49B]'
              />
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}

export default UserCard;
