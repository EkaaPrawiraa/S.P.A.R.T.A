import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  const token = request.cookies.get('auth_token')?.value

  // If accessing app routes without token, redirect to login
  if (request.nextUrl.pathname.startsWith('/app') && !token) {
    return NextResponse.redirect(new URL('/login', request.url))
  }

  // If accessing auth routes with token, redirect to app
  if (request.nextUrl.pathname.startsWith('/login') && token) {
    return NextResponse.redirect(new URL('/app', request.url))
  }

  if (request.nextUrl.pathname.startsWith('/register') && token) {
    return NextResponse.redirect(new URL('/app', request.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/app/:path*', '/login', '/register'],
}
