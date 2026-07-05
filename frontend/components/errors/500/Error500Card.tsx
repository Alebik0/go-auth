"use server";

import Link from "next/link";
import Error500Icon from "./Error500Icon";

async function Error500Card() {
  return (
    <div className="flex items-center justify-center">
      <div className="max-w-lg text-center">
        <Error500Icon />

        <p className="on-error-container text-sm font-semibold uppercase tracking-wider">
          Error 500
        </p>

        <h1 className="on-surface mt-2 text-5xl font-bold tracking-tight">
          Internal Server Error
        </h1>

        <p className="on-surface-variant mt-6 text-lg">
          Sorry, something went wrong on our end. Please try again in a few
          moments.
        </p>

        <div className="mt-10 flex flex-col justify-center gap-4 sm:flex-row">
          <Link
            href="/"
            className="primary-container on-primary-container rounded-lg px-6 py-3 font-medium transition hover-secondary-container hover-on-secondary-container"
          >
            Go Home
          </Link>
        </div>

        <p className="on-surface-variant mt-8 text-sm">
          If the problem persists, please contact support.
        </p>
      </div>
    </div>
  );
}

export default Error500Card;
