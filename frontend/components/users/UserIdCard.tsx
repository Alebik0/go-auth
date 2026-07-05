function UserIdCard({ id }: { id: number }) {
  return (
    <>
      {/* User ID */}
      <div className="surface-container-high flex items-center justify-between rounded-[18px] px-5 py-3">
        <span className="on-surface-variant text-sm">User ID</span>
        <span className="on-primary-container font-semibold">#{id}</span>
      </div>
    </>
  );
}

export default UserIdCard;
