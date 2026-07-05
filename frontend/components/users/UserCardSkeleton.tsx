function UserCardSkeleton() {
  return (
    <div
      className="surface-container outline-variant border w-full max-w-2xl min-w-fit p-8"
      style={{ borderRadius: "25px" }}
    >
      <div className="flex items-start gap-5">
        {/* Avatar */}
        <div className="primary-container h-20 w-20 shrink-0 rounded-full" />

        <div className="flex-1">
          {/* User ID */}
          <div className="surface-container-high flex items-center justify-between rounded-[18px] px-5 py-3">
            <div className="surface-container-highest h-4 w-16 rounded" />
            <div className="surface-container-highest h-5 w-14 rounded" />
          </div>

          {/* Name */}
          <div className="mt-4">
            <div className="surface-container-high mb-2 h-4 w-12 rounded" />
            <div className="surface-container-high h-10 w-full rounded-xl" />
          </div>

          {/* Description */}
          <div className="mt-4">
            <div className="surface-container-high mb-2 h-4 w-24 rounded" />
            <div className="surface-container-high rounded-xl p-4">
              <div className="surface-container-highest h-4 w-full rounded" />
              <div className="surface-container-highest mt-3 h-4 w-11/12 rounded" />
              <div className="surface-container-highest mt-3 h-4 w-4/5 rounded" />
              <div className="surface-container-highest mt-3 h-4 w-2/3 rounded" />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default UserCardSkeleton;
