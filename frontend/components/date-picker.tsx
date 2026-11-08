'use client'

import * as React from 'react'
import { format, setMonth, setYear, isBefore, isAfter, startOfDay } from 'date-fns'
import { Calendar, ChevronDown } from 'lucide-react'
import { DayPicker } from 'react-day-picker'
import 'react-day-picker/dist/style.css'

interface DatePickerProps {
  value?: string
  onChange: (date: string) => void
  placeholder?: string
  required?: boolean
  maxDate?: Date
  disableFutureDates?: boolean
}

const months = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December'
]

export function DatePicker({ 
  value, 
  onChange, 
  placeholder = 'Pick a date', 
  required = false,
  maxDate,
  disableFutureDates = false
}: DatePickerProps) {
  const [open, setOpen] = React.useState(false)
  const [selectedDate, setSelectedDate] = React.useState<Date | undefined>(
    value ? new Date(value) : undefined
  )
  const [displayDate, setDisplayDate] = React.useState<Date>(
    value ? new Date(value) : new Date()
  )

  const today = startOfDay(new Date())
  const effectiveMaxDate = disableFutureDates ? today : maxDate

  const handleDateSelect = (date: Date | undefined) => {
    if (date) {
      // Check if date is in the future
      if (disableFutureDates && isAfter(date, today)) {
        return
      }
      // Check if date exceeds maxDate
      if (effectiveMaxDate && isAfter(date, effectiveMaxDate)) {
        return
      }
      setSelectedDate(date)
      setDisplayDate(date)
      onChange(date.toISOString())
      setOpen(false)
    }
  }

  const handleMonthChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newDate = setMonth(displayDate, parseInt(e.target.value))
    setDisplayDate(newDate)
  }

  const handleYearChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newDate = setYear(displayDate, parseInt(e.target.value))
    setDisplayDate(newDate)
  }

  // Generate year range (100 years back, 10 years forward)
  const currentYear = new Date().getFullYear()
  const years = Array.from({ length: 111 }, (_, i) => currentYear - 100 + i)

  return (
    <div className="relative">
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 placeholder-gray-400 flex items-center justify-between bg-white hover:bg-gray-50"
      >
        <span className={selectedDate ? 'text-gray-900' : 'text-gray-400'}>
          {selectedDate ? format(selectedDate, 'MMM d, yyyy') : placeholder}
        </span>
        <Calendar className="w-4 h-4 text-gray-400" />
      </button>

      {open && (
        <div className="absolute top-full mt-1 left-0 z-50 bg-white border border-gray-300 rounded-lg shadow-lg p-4 w-80">
          {/* Month and Year Selectors */}
          <div className="flex gap-2 mb-4">
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-600 mb-1">Month</label>
              <select
                value={displayDate.getMonth()}
                onChange={handleMonthChange}
                className="w-full px-2 py-1 border border-gray-300 rounded text-sm text-gray-900 bg-white hover:bg-gray-50"
              >
                {months.map((month, index) => (
                  <option key={month} value={index}>{month}</option>
                ))}
              </select>
            </div>
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-600 mb-1">Year</label>
              <select
                value={displayDate.getFullYear()}
                onChange={handleYearChange}
                className="w-full px-2 py-1 border border-gray-300 rounded text-sm text-gray-900 bg-white hover:bg-gray-50"
              >
                {years.map(year => (
                  <option key={year} value={year}>{year}</option>
                ))}
              </select>
            </div>
          </div>

          {/* Calendar */}
          <DayPicker
            mode="single"
            month={displayDate}
            onMonthChange={setDisplayDate}
            selected={selectedDate}
            onSelect={handleDateSelect}
            disabled={(date) => {
              // Disable future dates if disableFutureDates is true
              if (disableFutureDates && isAfter(date, today)) {
                return true
              }
              // Disable dates after maxDate
              if (effectiveMaxDate && isAfter(date, effectiveMaxDate)) {
                return true
              }
              return false
            }}
            className="text-gray-900"
            classNames={{
              caption: 'text-sm font-semibold text-gray-900 mb-4 hidden',
              head_cell: 'text-xs font-medium text-gray-600 w-8 h-8',
              cell: 'text-xs text-center',
              day: 'w-8 h-8 rounded hover:bg-blue-50 cursor-pointer text-gray-900',
              day_selected: 'bg-blue-500 text-white hover:bg-blue-600',
              day_today: 'font-bold text-blue-600',
              day_outside: 'text-gray-400',
              day_disabled: 'text-gray-300 cursor-not-allowed',
              nav_button: 'p-1 text-gray-600 hover:text-gray-900',
              nav_button_previous: 'absolute left-2',
              nav_button_next: 'absolute right-2',
            }}
          />
        </div>
      )}

      {open && (
        <div className="fixed inset-0 z-40" onClick={() => setOpen(false)} />
      )}
    </div>
  )
}
