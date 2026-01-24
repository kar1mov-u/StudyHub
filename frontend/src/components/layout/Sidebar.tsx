import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import { Home, BookOpen, Calendar } from 'lucide-react'
import { cn } from '@/lib/utils'

const Sidebar: React.FC = () => {
  const location = useLocation()

  const navItems = [
    { path: '/', label: 'Dashboard', icon: Home },
    { path: '/modules', label: 'Modules', icon: BookOpen },
    { path: '/academic-terms', label: 'Academic Terms', icon: Calendar },
  ]

  return (
    <aside className="w-64 border-r bg-white min-h-[calc(100vh-73px)]">
      <nav className="p-4 space-y-2">
        {navItems.map((item) => {
          const Icon = item.icon
          const isActive = location.pathname === item.path
          
          return (
            <Link
              key={item.path}
              to={item.path}
              className={cn(
                "flex items-center gap-3 px-4 py-3 rounded-lg transition-colors",
                isActive
                  ? "bg-primary text-primary-foreground"
                  : "text-gray-700 hover:bg-gray-100"
              )}
            >
              <Icon className="h-5 w-5" />
              <span className="font-medium">{item.label}</span>
            </Link>
          )
        })}
      </nav>
    </aside>
  )
}

export default Sidebar
