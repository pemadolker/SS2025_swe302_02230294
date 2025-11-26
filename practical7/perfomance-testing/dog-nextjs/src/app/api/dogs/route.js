import { NextResponse } from "next/server";

export async function GET(request) {
  const { searchParams } = new URL(request.url);
  const breed = searchParams.get("breed");

  try {
    let url = "https://dog.ceo/api/breeds/image/random";
    
    if (breed) {
      url = `https://dog.ceo/api/breed/${breed}/images/random`;
    }

    const response = await fetch(url);
    const data = await response.json();

    return NextResponse.json(data);
  } catch (error) {
    return NextResponse.json(
      { status: "error", message: "Failed to fetch dog image" },
      { status: 500 }
    );
  }
}