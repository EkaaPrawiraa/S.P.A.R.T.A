'use client'

import { AppHeader } from '@/components/app-header'
import { AppNav } from '@/components/app-nav'
import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { isAuthenticated } from '@/lib/auth'

export default function AppLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const router = useRouter()
  const [isMounted, setIsMounted] = useState(false)

  useEffect(() => {
    setIsMounted(true)
    if (!isAuthenticated()) {
      router.push('/login')
    }
  }, [router])

  if (!isMounted) {
    return null
  }

  return (
    <div className="flex flex-col h-screen md:flex-row">
      {/* Desktop Sidebar */}
      <aside className="hidden md:flex md:flex-col md:w-64 border-r border-border/40 bg-background">
        <div className="p-6">
          <h1 className="text-2xl font-bold text-foreground">φ Gym</h1>
        </div>
        <nav className="flex-1 overflow-y-auto px-4 pb-4">
          <AppNav variant="sidebar" />
        </nav>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col">
        <AppHeader />
        <main className="flex-1 overflow-y-auto">
          {children}
        </main>

        {/* Mobile Bottom Nav */}
        <nav className="md:hidden border-t border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky bottom-0">
          <div className="px-4 py-3 overflow-x-auto">
            <AppNav variant="mobile" />
          </div>
        </nav>
      </div>
    </div>
  )
}
