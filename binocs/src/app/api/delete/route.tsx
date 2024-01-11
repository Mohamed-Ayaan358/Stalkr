import { NextRequest, NextResponse } from "next/server";

// Define the POST function for your API route
export async function POST(req: NextRequest, res: NextResponse) {
  try {
    const data = await req.json();
    console.log(data);

    const apiResponse = await fetch("http://localhost:8080/delete", {
      method: "POST",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin": "*",
      },
      body: JSON.stringify(data),
    });

    if (apiResponse.ok) {
      const responseData = await apiResponse.json();
      console.log(responseData);
      return NextResponse.json(
        { data: responseData },
        { status: apiResponse.status }
      );
    } else {
      console.log(apiResponse);
      console.log("Invalid HTTP status code:", apiResponse.status);
      return NextResponse.error();
    }
  } catch (error) {
    console.error("Error making DELETE request to API:", error);
    return NextResponse.error();
  }
}
