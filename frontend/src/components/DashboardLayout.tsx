import type { ReactNode } from 'react'
import { Link, useLocation } from '@tanstack/react-router'
import {
  Home,
  PiggyBank,
  Wallet,
  Target,
  BookOpen,
  BarChart3,
  User,
  Settings,
  Menu,
  X,
} from 'lucide-react'
import { useState } from 'react'
import { useAuth } from '../lib/auth/context'
import { NotificationsDropdown } from './NotificationsDropdown'

interface DashboardLayoutProps {
  children: ReactNode
}

// Requirement 18.1: Main dashboard layout with sidebar navigation
export function DashboardLayout({ children }: DashboardLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const location = useLocation()
  const { user, logout } = useAuth()

  const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: Home },
    { name: 'Savings', href: '/savings', icon: PiggyBank },
    { name: 'Budget', href: '/budget', icon: Wallet },
    { name: 'Goals', href: '/goals', icon: Target },
    { name: 'Education', href: '/education', icon: BookOpen },
    { name: 'Analytics', href: '/analytics', icon: BarChart3 },
    { name: 'Profile', href: '/profile', icon: User },
    { name: 'Settings', href: '/settings', icon: Settings },
  ]

  const isActive = (href: string) => location.pathname === href

  const handleLogout = async () => {
    await logout()
  }

  return (
    <div className="min-h-screen bg-[var(--bg)]">
      {/* Mobile sidebar backdrop */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={`
          fixed top-0 left-0 z-50 h-full w-64 bg-[var(--island-bg)] border-r border-[var(--line)]
          transform transition-transform duration-200 ease-in-out
          lg:translate-x-0
          ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        `}
      >
        <div className="flex flex-col h-full">
          {/* Logo and close button */}
          <div className="flex items-center justify-between p-6 border-b border-[var(--line)]">
            <h1 className="text-2xl font-bold text-[var(--sea-ink)]">InSavein</h1>
            <button
              onClick={() => setSidebarOpen(false)}
              className="lg:hidden p-2 hover:bg-[var(--link-bg-hover)] rounded-lg"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          {/* Navigation */}
          <nav className="flex-1 p-4 space-y-1">
            {navigation.map((item) => {
              const Icon = item.icon
              const active = isActive(item.href)
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`
                    flex items-center gap-3 px-4 py-3 rounded-lg transition-colors
                    ${
                      active
                        ? 'bg-[var(--link-bg-hover)] text-[var(--sea-ink)] font-semibold'
                        : 'text-[var(--sea-ink-soft)] hover:bg-[var(--link-bg-hover)] hover:text-[var(--sea-ink)]'
                    }
                  `}
                  onClick={() => setSidebarOpen(false)}
                >
                  <Icon className="w-5 h-5" />
                  <span>{item.name}</span>
                </Link>
              )
            })}
          </nav>

          {/* User section */}
          <div className="p-4 border-t border-[var(--line)]">
            <div className="flex items-center gap-3 mb-3">
              <div className="w-10 h-10 rounded-full bg-[var(--link-bg-hover)] flex items-center justify-center">
                <User className="w-5 h-5 text-[var(--sea-ink-soft)]" />
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-semibold text-[var(--sea-ink)] truncate">
                  {user?.first_name} {user?.last_name}
                </p>
                <p className="text-xs text-[var(--sea-ink-soft)] truncate">
                  {user?.email}
                </p>
              </div>
            </div>
            <button
              onClick={handleLogout}
              className="w-full px-4 py-2 text-sm border border-[var(--line)] rounded-lg hover:bg-[var(--link-bg-hover)] transition-colors"
            >
              Logout
            </button>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <div className="lg:pl-64">
        {/* Desktop header */}
        <header className="hidden lg:block sticky top-0 z-30 bg-[var(--island-bg)] border-b border-[var(--line)] px-8 py-4">
          <div className="flex items-center justify-end">
            <NotificationsDropdown />
          </div>
        </header>

        {/* Mobile header */}
        <header className="lg:hidden sticky top-0 z-30 bg-[var(--island-bg)] border-b border-[var(--line)] px-4 py-3">
          <div className="flex items-center justify-between">
            <button
              onClick={() => setSidebarOpen(true)}
              className="p-2 hover:bg-[var(--link-bg-hover)] rounded-lg"
            >
              <Menu className="w-6 h-6" />
            </button>
            <h1 className="text-xl font-bold">InSavein</h1>
            <NotificationsDropdown />
          </div>
        </header>

        {/* Page content */}
        <main className="p-4 lg:p-8">{children}</main>
      </div>
    </div>
  )
}
