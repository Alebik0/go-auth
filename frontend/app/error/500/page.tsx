"use server";

import Error500 from "@/components/Error500";

async function Page500() {
  return (
    <>
      <Error500 />
    </>
  );
}

export default Page500;
