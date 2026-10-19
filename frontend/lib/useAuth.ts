'use client'

import { useEffect, useState } from 'react'
import { useRouter, usePathname } from 'next/navigation'

export function useAuth() {
  const router = useRouter()
  const pathname = usePathname()
  const [isLoading, setIsLoading] = useState(true)
  const [isAuthenticated, setIsAuthenticated] = useState(false)

  useEffect(() => {
    // Check if token exists in localStorage
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null
    const publicRoutes = ['/auth/login', '/auth/signup', '/auth/forgot-password']
    const isPublicRoute = publicRoutes.some(route => pathname.startsWith(route))

    if (!token && !isPublicRoute) {
      // No token and on protected route - redirect to login
      router.push('/auth/login')
    } else if (token && isPublicRoute) {
      // Has token and on auth route - redirect to dashboard
      router.push('/')
    } else if (token) {
      setIsAuthenticated(true)
    }

    setIsLoading(false)
  }, [router, pathname])

  return { isLoading, isAuthenticated }
}
