'use client'

import { useRouter } from 'next/navigation'
import BiodataForm from '@/components/biodata-form'

export default function OnboardingPage() {
  const router = useRouter()

  const handleSkip = () => {
    router.push('/')
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-50">
      <div className="container mx-auto">
        <div className="flex justify-end p-6">
          <button
            onClick={handleSkip}
            className="text-sm font-medium text-gray-600 hover:text-gray-900 underline"
          >
            Skip for now
          </button>
        </div>
        <BiodataForm 
          isOnboarding={true}
          onComplete={() => router.push('/')}
        />
      </div>
    </div>
  )
}
