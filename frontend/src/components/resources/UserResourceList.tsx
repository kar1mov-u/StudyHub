import React from 'react'
import UserResourceCard from './UserResourceCard'
import { Loader2 } from 'lucide-react'
import type { UserResource } from '@/types'

interface UserResourceListProps {
  resources: UserResource[]
  isLoading?: boolean
  emptyMessage?: string
  showDelete?: boolean
  onDelete?: (resourceId: string) => void
}

const UserResourceList: React.FC<UserResourceListProps> = ({
  resources,
  isLoading = false,
  emptyMessage = 'No resources found',
  showDelete = false,
  onDelete,
}) => {
  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (resources.length === 0) {
    return (
      <div className="text-center py-12 text-muted-foreground">
        <p>{emptyMessage}</p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {resources.map((resource) => (
        <UserResourceCard 
          key={resource.ID} 
          resource={resource} 
          showDelete={showDelete}
          onDelete={onDelete}
        />
      ))}
    </div>
  )
}

export default UserResourceList
