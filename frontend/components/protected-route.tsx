'use client'

import { useEffect, useState } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { Loader2 } from 'lucide-react'

function getTokenFromCookie() {
  if (typeof document === 'undefined') return null
  
  const cookies = document.cookie.split(';')
  for (let cookie of cookies) {
    const [name, value] = cookie.trim().split('=')
    if (name === 'auth_token') {
      return decodeURIComponent(value)
    }
  }
  return null
}

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const router = useRouter()
  const pathname = usePathname()
  const [isLoading, setIsLoading] = useState(true)
  const [isAuthenticated, setIsAuthenticated] = useState(false)

  useEffect(() => {
    const token = getTokenFromCookie()
    const publicRoutes = ['/auth/login', '/auth/signup', '/auth/forgot-password', '/auth/verify-email', '/auth/reset-password']
    const isPublicRoute = publicRoutes.some(route => pathname.startsWith(route))
    // Routes where authenticated users should be redirected away
    const authRedirectRoutes = ['/auth/login', '/auth/signup', '/auth/forgot-password']
    const isAuthRedirectRoute = authRedirectRoutes.some(route => pathname.startsWith(route))

    if (!token && !isPublicRoute) {
      router.push('/auth/login')
    } else if (token && isAuthRedirectRoute) {
      // Only redirect away from login/signup/forgot-password, allow verify-email/reset-password
      router.push('/')
    } else if (token) {
      setIsAuthenticated(true)
    } else if (isPublicRoute) {
      setIsAuthenticated(true) // Allow public routes without token
    }

    setIsLoading(false)
  }, [router, pathname])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
        <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
      </div>
    )
  }

  if (!isAuthenticated) {
    return null
  }

  return <>{children}</>
}
