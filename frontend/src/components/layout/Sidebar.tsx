import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import { Home, BookOpen, Calendar, User, Shield } from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAuth } from '@/context/AuthContext'

const Sidebar: React.FC = () => {
  const location = useLocation()
  const { user } = useAuth()

  return (
    <aside className="w-64 border-r bg-white min-h-[calc(100vh-73px)]">
      <nav className="p-4 space-y-2">
        <Link
          to="/"
          className={cn(
            "flex items-center gap-3 px-4 py-3 rounded-lg transition-colors",
            location.pathname === '/'
              ? "bg-primary text-primary-foreground"
              : "text-gray-700 hover:bg-gray-100"
          )}
        >
          <Home className="h-5 w-5" />
          <span className="font-medium">Dashboard</span>
        </Link>
        
        <Link
          to="/modules"
          className={cn(
            "flex items-center gap-3 px-4 py-3 rounded-lg transition-colors",
            location.pathname.startsWith('/modules')
              ? "bg-primary text-primary-foreground"
              : "text-gray-700 hover:bg-gray-100"
          )}
        >
          <BookOpen className="h-5 w-5" />
          <span className="font-medium">Modules</span>
        </Link>
        
        <Link
          to="/academic-terms"
          className={cn(
            "flex items-center gap-3 px-4 py-3 rounded-lg transition-colors",
            location.pathname === '/academic-terms'
              ? "bg-primary text-primary-foreground"
              : "text-gray-700 hover:bg-gray-100"
          )}
        >
          <Calendar className="h-5 w-5" />
          <span className="font-medium">Academic Terms</span>
        </Link>

        {user?.IsAdmin && (
          <Link
            to="/admin"
            className={cn(
              "flex items-center gap-3 px-4 py-3 rounded-lg transition-colors",
              location.pathname === '/admin'
                ? "bg-primary text-primary-foreground"
                : "text-gray-700 hover:bg-gray-100"
            )}
          >
            <Shield className="h-5 w-5" />
            <span className="font-medium">Admin Panel</span>
          </Link>
        )}

        {user && (
          <Link
            to={`/users/${user.ID}`}
            className={cn(
              "flex items-center gap-3 px-4 py-3 rounded-lg transition-colors",
              location.pathname.startsWith('/users/')
                ? "bg-primary text-primary-foreground"
                : "text-gray-700 hover:bg-gray-100"
            )}
          >
            <User className="h-5 w-5" />
            <span className="font-medium">My Profile</span>
          </Link>
        )}
      </nav>
    </aside>
  )
}

export default Sidebar
