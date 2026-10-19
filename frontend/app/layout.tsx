import { Analytics } from '@vercel/analytics/next'
import type { Metadata, Viewport } from 'next'
import './globals.css'
import { ProtectedRoute } from '@/components/protected-route'

export const metadata: Metadata = {
  title: 'Breakthrough Palace Mission · Family Records',
  description: 'Breakthrough Palace Mission - Manage family records, meetings, attendance, and contributions.',
  generator: 'v0.app',
  icons: {
    icon: [
      {
        url: '/icon-light-32x32.png',
        media: '(prefers-color-scheme: light)',
      },
      {
        url: '/icon-dark-32x32.png',
        media: '(prefers-color-scheme: dark)',
      },
      {
        url: '/icon.svg',
        type: 'image/svg+xml',
      },
    ],
    apple: '/apple-icon.png',
  },
}

export const viewport: Viewport = {
  colorScheme: 'light',
  themeColor: [{ color: '#f4f8fc' }],
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en">
      <body className="antialiased">
        <ProtectedRoute>
          {children}
        </ProtectedRoute>
        {process.env.NODE_ENV === 'production' && <Analytics />}
      </body>
    </html>
  )
}
