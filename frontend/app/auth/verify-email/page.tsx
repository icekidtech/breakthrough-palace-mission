'use client'

import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { CheckCircle, XCircle, Loader } from 'lucide-react'
import Link from 'next/link'

export default function VerifyEmailPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const token = searchParams.get('token')
  
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading')
  const [message, setMessage] = useState('')

  useEffect(() => {
    const verifyEmail = async () => {
      if (!token) {
        setStatus('error')
        setMessage('Invalid verification link. No token provided.')
        return
      }

      try {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL
        const response = await fetch(`${apiUrl}/api/auth/verify-email`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ token }),
          credentials: 'include',
        })

        const data = await response.json()

        if (!response.ok) {
          setStatus('error')
          setMessage(data.error || 'Failed to verify email. The link may have expired.')
          return
        }

        setStatus('success')
        setMessage('Your email has been verified successfully!')

        // Redirect to dashboard after 3 seconds
        setTimeout(() => {
          router.push('/dashboard')
        }, 3000)
      } catch (err) {
        setStatus('error')
        setMessage(err instanceof Error ? err.message : 'Something went wrong during verification')
      }
    }

    verifyEmail()
  }, [token, router])

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-100">
      {/* Decorative elements */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-10 right-20 w-64 h-64 bg-blue-200/20 rounded-full blur-3xl" />
        <div className="absolute bottom-20 left-10 w-96 h-96 bg-indigo-200/20 rounded-full blur-3xl" />
      </div>

      <div className="relative z-10 flex items-center justify-center min-h-screen px-4 py-6 sm:px-6 lg:px-8">
        <div className="w-full max-w-md">
          {/* Card */}
          <div className="bg-white/80 backdrop-blur-xl rounded-3xl shadow-2xl p-8 border border-white/20">
            <div className="text-center">
              {status === 'loading' && (
                <>
                  <div className="flex justify-center mb-6">
                    <Loader className="w-16 h-16 text-blue-600 animate-spin" />
                  </div>
                  <h1 className="text-3xl font-bold text-gray-900 mb-2">Verifying Email</h1>
                  <p className="text-gray-600">Please wait while we verify your email address...</p>
                </>
              )}

              {status === 'success' && (
                <>
                  <div className="flex justify-center mb-6">
                    <CheckCircle className="w-16 h-16 text-green-600" />
                  </div>
                  <h1 className="text-3xl font-bold text-gray-900 mb-2">Email Verified!</h1>
                  <p className="text-gray-600 mb-6">{message}</p>
                  <p className="text-sm text-gray-500 mb-6">Redirecting to dashboard in 3 seconds...</p>
                  <Link
                    href="/dashboard"
                    className="inline-block bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-semibold py-3 px-6 rounded-xl transition-all duration-200 shadow-lg hover:shadow-xl"
                  >
                    Go to Dashboard Now
                  </Link>
                </>
              )}

              {status === 'error' && (
                <>
                  <div className="flex justify-center mb-6">
                    <XCircle className="w-16 h-16 text-red-600" />
                  </div>
                  <h1 className="text-3xl font-bold text-gray-900 mb-2">Verification Failed</h1>
                  <p className="text-gray-600 mb-6">{message}</p>
                  <div className="space-y-3">
                    <p className="text-sm text-gray-500">
                      If the link has expired, request a new verification email.
                    </p>
                    <Link
                      href="/auth/login"
                      className="inline-block bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-semibold py-3 px-6 rounded-xl transition-all duration-200 shadow-lg hover:shadow-xl"
                    >
                      Return to Login
                    </Link>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
