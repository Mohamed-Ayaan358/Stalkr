import { NextRequest, NextResponse } from "next/server";

export async function GET(req: NextRequest) {
  const res = await fetch("http://localhost:8080/websites", {
    cache: "no-cache",
  });

  if (res.ok) {
    const responseData = await res.json();
    return NextResponse.json({ data: responseData }, { status: res.status });
  } else {
    console.log("Invalid HTTP status code:", res.status);
    return NextResponse.error();
  }
}
