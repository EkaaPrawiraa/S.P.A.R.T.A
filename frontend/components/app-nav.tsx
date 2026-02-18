'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { cn } from '@/lib/utils'
import {
  LayoutDashboard,
  Dumbbell,
  Split,
  Zap,
  Apple,
  BarChart3,
  Brain,
  Calendar,
  Settings,
} from 'lucide-react'

const navItems = [
  { href: '/app', label: 'Dashboard', icon: LayoutDashboard },
  { href: '/app/workouts', label: 'Workouts', icon: Dumbbell },
  { href: '/app/splits', label: 'Splits', icon: Split },
  { href: '/app/exercises', label: 'Exercises', icon: Zap },
  { href: '/app/nutrition', label: 'Nutrition', icon: Apple },
  { href: '/app/analytics', label: 'Analytics', icon: BarChart3 },
  { href: '/app/ai-tools', label: 'AI Tools', icon: Brain },
  { href: '/app/planner', label: 'Planner', icon: Calendar },
  { href: '/app/settings', label: 'Settings', icon: Settings },
]

interface AppNavProps {
  className?: string
  variant?: 'sidebar' | 'mobile'
}

export function AppNav({ className, variant = 'sidebar' }: AppNavProps) {
  const pathname = usePathname()

  if (variant === 'mobile') {
    return (
      <nav className={cn('flex gap-2 overflow-x-auto pb-2', className)}>
        {navItems.map((item) => {
          const isActive =
            pathname === item.href ||
            (item.href !== '/app' && pathname.startsWith(item.href))
          const Icon = item.icon

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                'flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium whitespace-nowrap transition-colors',
                isActive
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:text-foreground'
              )}
            >
              <Icon className="h-4 w-4" />
              <span>{item.label}</span>
            </Link>
          )
        })}
      </nav>
    )
  }

  return (
    <nav className={cn('space-y-1', className)}>
      {navItems.map((item) => {
        const isActive =
          pathname === item.href ||
          (item.href !== '/app' && pathname.startsWith(item.href))
        const Icon = item.icon

        return (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              'flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-colors',
              isActive
                ? 'bg-primary text-primary-foreground shadow-lg ring-1 ring-primary/30'
                : 'text-muted-foreground hover:text-foreground hover:bg-accent'
            )}
          >
            <Icon className="h-5 w-5" />
            <span>{item.label}</span>
          </Link>
        )
      })}
    </nav>
  )
}
