import { NextResponse, NextRequest } from "next/server";

export function proxy(request: NextRequest) {
  const hasRefresh = request.cookies.has("refresh_token");

  if (!hasRefresh) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("returnTo", request.nextUrl.pathname);
    return NextResponse.redirect(loginUrl);
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    { source: "/((?:room/.*|)$)" },
  ],
}