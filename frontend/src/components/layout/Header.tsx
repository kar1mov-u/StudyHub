import React from 'react'
import { BookOpen } from 'lucide-react'

const Header: React.FC = () => {
  return (
    <header className="border-b bg-white">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center gap-2">
          <BookOpen className="h-8 w-8 text-primary" />
          <h1 className="text-2xl font-bold text-gray-900">StudyHub</h1>
        </div>
      </div>
    </header>
  )
}

export default Header
