'use client'

import { useState, useEffect } from 'react'
import { AlertCircle, CheckCircle2, Loader2 } from 'lucide-react'

interface Biodata {
  id: string
  member_id: string
  full_name: string
  date_of_birth?: string
  gender: string
  marital_status: string
  nationality: string
  phone_number: string
  email_address: string
  home_address: string
  emergency_contact: string
  emergency_phone: string
  date_joined_church?: string
  previous_church?: string
  department: string
  family_name: string
  occupation: string
  company_business_name?: string
  educational_qualification: string
  birthday?: string
  wedding_anniversary?: string
  special_occasions?: string
  welfare_support_needs?: string
  photograph?: string
  social_media_handles?: string
  special_interests?: string
  contact_consent?: boolean
  is_complete: boolean
  created_at: string
  updated_at: string
}

interface BiodataFormProps {
  isOnboarding?: boolean
  onComplete?: () => void
}

export default function BiodataForm({ isOnboarding = false, onComplete }: BiodataFormProps) {
  const [biodata, setBiodata] = useState<Partial<Biodata>>({})
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const [completionPercent, setCompletionPercent] = useState(0)

  // Load biodata
  useEffect(() => {
    const fetchBiodata = async () => {
      try {
        const response = await fetch('/api/biodata', {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
        })
        if (!response.ok) throw new Error('Failed to load biodata')
        const data = await response.json()
        setBiodata(data.biodata)
        calculateCompletion(data.biodata)
      } catch (err) {
        setError('Failed to load biodata')
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    fetchBiodata()
  }, [])

  const calculateCompletion = (data: Partial<Biodata>) => {
    const requiredFields = [
      'full_name',
      'phone_number',
      'email_address',
      'home_address',
      'emergency_contact',
      'emergency_phone',
      'date_of_birth',
      'gender',
      'marital_status',
      'date_joined_church',
      'department',
      'family_name',
      'occupation',
      'educational_qualification',
    ]
    const filledFields = requiredFields.filter(field => data[field as keyof Biodata])
    setCompletionPercent(Math.round((filledFields.length / requiredFields.length) * 100))
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value, type } = e.target
    const newValue = type === 'checkbox' ? (e.target as HTMLInputElement).checked : value
    setBiodata(prev => ({ ...prev, [name]: newValue }))
  }

  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setBiodata(prev => ({ ...prev, [name]: value ? new Date(value).toISOString() : undefined }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    setError('')
    setSuccess(false)

    try {
      const response = await fetch('/api/biodata', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify(biodata),
      })

      if (!response.ok) throw new Error('Failed to save biodata')
      const data = await response.json()
      setBiodata(data.biodata)
      calculateCompletion(data.biodata)
      setSuccess(true)
      if (onComplete) onComplete()
      setTimeout(() => setSuccess(false), 3000)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save biodata')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center p-8">
        <Loader2 className="w-6 h-6 animate-spin text-blue-500" />
      </div>
    )
  }

  return (
    <div className="w-full max-w-4xl mx-auto p-4 sm:p-6">
      {isOnboarding && (
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Complete Your Profile</h1>
          <p className="text-gray-600 mb-4">Help us know you better. Complete these details to get started.</p>
          
          {/* Completion Progress */}
          <div className="bg-white rounded-lg border border-gray-200 p-4">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-medium text-gray-700">Profile Completion</span>
              <span className="text-sm font-bold text-blue-600">{completionPercent}%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div 
                className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                style={{ width: `${completionPercent}%` }}
              />
            </div>
          </div>
        </div>
      )}

      {error && (
        <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg flex gap-3">
          <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
          <p className="text-red-700 text-sm">{error}</p>
        </div>
      )}

      {success && (
        <div className="mb-4 p-4 bg-green-50 border border-green-200 rounded-lg flex gap-3">
          <CheckCircle2 className="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" />
          <p className="text-green-700 text-sm">Profile saved successfully!</p>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-8">
        {/* Personal Information Section */}
        <section className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Personal Information</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Full Name <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="full_name"
                value={biodata.full_name || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Date of Birth <span className="text-red-500">*</span>
              </label>
              <input
                type="date"
                name="date_of_birth"
                value={biodata.date_of_birth ? biodata.date_of_birth.split('T')[0] : ''}
                onChange={handleDateChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Gender <span className="text-red-500">*</span>
              </label>
              <select
                name="gender"
                value={biodata.gender || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="">Select gender</option>
                <option value="Male">Male</option>
                <option value="Female">Female</option>
                <option value="Other">Other</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Marital Status <span className="text-red-500">*</span>
              </label>
              <select
                name="marital_status"
                value={biodata.marital_status || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="">Select status</option>
                <option value="Single">Single</option>
                <option value="Married">Married</option>
                <option value="Divorced">Divorced</option>
                <option value="Widowed">Widowed</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Nationality
              </label>
              <input
                type="text"
                name="nationality"
                value={biodata.nationality || ''}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
        </section>

        {/* Contact Information Section */}
        <section className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Contact Information</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Phone Number <span className="text-red-500">*</span>
              </label>
              <input
                type="tel"
                name="phone_number"
                value={biodata.phone_number || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Email Address <span className="text-red-500">*</span>
              </label>
              <input
                type="email"
                name="email_address"
                value={biodata.email_address || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Home Address <span className="text-red-500">*</span>
              </label>
              <textarea
                name="home_address"
                value={biodata.home_address || ''}
                onChange={handleChange}
                required
                rows={2}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Emergency Contact <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="emergency_contact"
                value={biodata.emergency_contact || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Emergency Phone <span className="text-red-500">*</span>
              </label>
              <input
                type="tel"
                name="emergency_phone"
                value={biodata.emergency_phone || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
        </section>

        {/* Church Information Section */}
        <section className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Church Information</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Date Joined Church <span className="text-red-500">*</span>
              </label>
              <input
                type="date"
                name="date_joined_church"
                value={biodata.date_joined_church ? biodata.date_joined_church.split('T')[0] : ''}
                onChange={handleDateChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Previous Church
              </label>
              <input
                type="text"
                name="previous_church"
                value={biodata.previous_church || ''}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Department/Unit <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="department"
                value={biodata.department || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="e.g., Choir, Ushering, etc."
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Family Name <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="family_name"
                value={biodata.family_name || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
        </section>

        {/* Skills & Profession Section */}
        <section className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Skills & Profession</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Occupation <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="occupation"
                value={biodata.occupation || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Company/Business Name
              </label>
              <input
                type="text"
                name="company_business_name"
                value={biodata.company_business_name || ''}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Educational Qualification <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="educational_qualification"
                value={biodata.educational_qualification || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="e.g., BSc, HND, etc."
              />
            </div>
          </div>
        </section>

        {/* Welfare & Support Section */}
        <section className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Welfare & Support</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Birthday
              </label>
              <input
                type="date"
                name="birthday"
                value={biodata.birthday ? biodata.birthday.split('T')[0] : ''}
                onChange={handleDateChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Wedding Anniversary
              </label>
              <input
                type="date"
                name="wedding_anniversary"
                value={biodata.wedding_anniversary ? biodata.wedding_anniversary.split('T')[0] : ''}
                onChange={handleDateChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Special Occasions (e.g., anniversaries)
              </label>
              <textarea
                name="special_occasions"
                value={biodata.special_occasions || ''}
                onChange={handleChange}
                rows={2}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Areas Where You Need Church Support
              </label>
              <textarea
                name="welfare_support_needs"
                value={biodata.welfare_support_needs || ''}
                onChange={handleChange}
                rows={2}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="e.g., prayer, counseling, financial assistance, etc."
              />
            </div>
          </div>
        </section>

        {/* Optional Information Section */}
        <section className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Optional Information</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Photograph URL
              </label>
              <input
                type="url"
                name="photograph"
                value={biodata.photograph || ''}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="https://..."
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Social Media Handles
              </label>
              <input
                type="text"
                name="social_media_handles"
                value={biodata.social_media_handles || ''}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="e.g., @username, etc."
              />
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Special Interests/Hobbies
              </label>
              <textarea
                name="special_interests"
                value={biodata.special_interests || ''}
                onChange={handleChange}
                rows={2}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div className="md:col-span-2 flex items-center">
              <input
                type="checkbox"
                id="contact_consent"
                name="contact_consent"
                checked={biodata.contact_consent || false}
                onChange={handleChange}
                className="w-4 h-4 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
              />
              <label htmlFor="contact_consent" className="ml-2 text-sm text-gray-700">
                I consent to be contacted for church programs and activities
              </label>
            </div>
          </div>
        </section>

        {/* Submit Button */}
        <div className="flex justify-end gap-4">
          <button
            type="submit"
            disabled={saving}
            className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed flex items-center gap-2"
          >
            {saving && <Loader2 className="w-4 h-4 animate-spin" />}
            {saving ? 'Saving...' : 'Save Profile'}
          </button>
        </div>
      </form>
    </div>
  )
}
