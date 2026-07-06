"use server";

async function Error500Icon() {
  return (
    <div className="error-container mb-6 inline-flex h-24 w-24 items-center justify-center rounded-full">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 64 64"
        width="64"
        height="64"
      >
        <path
          d="
  M32 8
  Q34 8 35 10
  L57 49
  Q58 51 56 52
  L8 52
  Q6 51 7 49
  L29 10
  Q30 8 32 8
  Z
"
          fill="#E53935"
          stroke="#B71C1C"
          strokeWidth="2"
        />

        <rect x="29.5" y="20" width="5" height="20" rx="2.5" fill="white" />

        <circle cx="32" cy="47" r="3" fill="white" />
      </svg>
    </div>
  );
}

export default Error500Icon;
