import { NextRequest, NextResponse } from "next/server";

// Define the POST function for your API route
export async function POST(req: NextRequest, res: NextResponse) {
  try {
    // Parse the JSON data from the request body
    const data = await req.json();
    console.log(JSON.stringify(data));

    // Send a POST request with JSON data to "http://localhost:8080/add"
    const apiResponse = await fetch("http://localhost:8080/add", {
      method: "POST", // Specify the HTTP method
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin": "*",
      },
      body: JSON.stringify(data), // Include the serialized JSON data in the request body
    });

    if (apiResponse.ok) {
      const responseData = await apiResponse.json();
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
    console.error("Error making POST request to API:", error);
    return NextResponse.error();
  }
}
