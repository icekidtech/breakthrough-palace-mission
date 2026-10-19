import BiodataForm from '@/components/biodata-form'

export default function OnboardingPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-50">
      <div className="container mx-auto">
        <BiodataForm isOnboarding={true} />
      </div>
    </div>
  )
}
