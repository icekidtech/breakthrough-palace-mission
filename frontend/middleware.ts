import { NextRequest, NextResponse } from 'next/server'

export function middleware(request: NextRequest) {
  const token = request.cookies.get('token')?.value
  
  // Public routes that don't require authentication
  const publicRoutes = ['/auth/login', '/auth/signup', '/auth/forgot-password']
  
  // Check if the current path is a public route
  const isPublicRoute = publicRoutes.some(route => request.nextUrl.pathname.startsWith(route))
  
  // If no token and trying to access protected route, redirect to login
  if (!token && !isPublicRoute) {
    return NextResponse.redirect(new URL('/auth/login', request.url))
  }
  
  // If has token and trying to access auth routes, redirect to dashboard
  if (token && isPublicRoute) {
    return NextResponse.redirect(new URL('/', request.url))
  }
  
  return NextResponse.next()
}

export const config = {
  matcher: [
    '/((?!_next/static|_next/image|favicon.ico|public).*)',
  ],
}
