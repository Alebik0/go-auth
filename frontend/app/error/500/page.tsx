"use server";

import Error500Card from "@/components/errors/500/Error500Card";

async function Page500() {
  return (
    <>
      <Error500Card />
    </>
  );
}

export default Page500;
