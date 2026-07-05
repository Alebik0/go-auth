function UserAvatar({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <>
      {/* Avatar */}
      <div className="primary-container flex h-20 w-20 shrink-0 items-center justify-center rounded-full">
        <span className="on-primary-container text-5xl font-bold">
          {children}
        </span>
      </div>
    </>
  );
}

export default UserAvatar;
