// app/api/dogs/breeds/route.js
import { NextResponse } from "next/server";

export async function GET() {
  try {
    const response = await fetch("https://dog.ceo/api/breeds/list/all");
    const data = await response.json();

    return NextResponse.json(data);
  } catch (error) {
    return NextResponse.json(
      { status: "error", message: "Failed to fetch breeds list" },
      { status: 500 }
    );
  }
}