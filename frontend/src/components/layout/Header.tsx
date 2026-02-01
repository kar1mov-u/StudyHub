import React from 'react'
import { BookOpen, LogOut, User, Shield } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/context/AuthContext'

const Header: React.FC = () => {
  const { user, logout } = useAuth()

  return (
    <header className="border-b bg-white">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <BookOpen className="h-8 w-8 text-primary" />
            <h1 className="text-2xl font-bold text-gray-900">StudyHub</h1>
          </div>
          
          {user && (
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <User className="h-5 w-5 text-muted-foreground" />
                <div className="text-sm">
                  <p className="font-medium">
                    {user.FirstName} {user.LastName}
                  </p>
                  <p className="text-muted-foreground text-xs">{user.Email}</p>
                </div>
                {user.IsAdmin && (
                  <Badge variant="default" className="ml-2">
                    <Shield className="h-3 w-3 mr-1" />
                    Admin
                  </Badge>
                )}
              </div>
              <Button variant="outline" size="sm" onClick={logout}>
                <LogOut className="h-4 w-4 mr-2" />
                Logout
              </Button>
            </div>
          )}
        </div>
      </div>
    </header>
  )
}

export default Header
