import { NextRequest, NextResponse } from 'next/server'

export function middleware(request: NextRequest) {
  const token = request.cookies.get('auth_token')?.value
  const pathname = request.nextUrl.pathname
  
  // Public routes that don't require authentication
  const publicRoutes = ['/auth/login', '/auth/signup', '/auth/forgot-password', '/auth/verify-email', '/auth/reset-password']
  
  // Routes that are allowed without redirecting (onboarding for incomplete biodata)
  const allowedRoutes = ['/onboarding']
  
  // Check if the current path is a public route
  const isPublicRoute = publicRoutes.some(route => pathname.startsWith(route))
  
  // Check if the current path is an allowed route
  const isAllowedRoute = allowedRoutes.some(route => pathname.startsWith(route))
  
  // If no token and trying to access protected route, redirect to login
  if (!token && !isPublicRoute && !isAllowedRoute) {
    return NextResponse.redirect(new URL('/auth/login', request.url))
  }
  
  // If has token and trying to access auth routes (except password reset/verification), redirect to dashboard
  if (token && publicRoutes.slice(0, 3).some(route => pathname.startsWith(route))) {
    return NextResponse.redirect(new URL('/', request.url))
  }
  
  return NextResponse.next()
}

export const config = {
  matcher: [
    '/((?!_next/static|_next/image|favicon.ico|public).*)',
  ],
}
