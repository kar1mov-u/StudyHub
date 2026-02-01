import React from 'react'
import ResourceCard from './ResourceCard'
import type { Resource } from '@/types'
import { FileX } from 'lucide-react'

interface ResourceListProps {
  resources: Resource[]
  isLoading?: boolean
  emptyMessage?: string
}

const ResourceList: React.FC<ResourceListProps> = ({ 
  resources, 
  isLoading,
  emptyMessage = 'No resources available'
}) => {
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {[...Array(3)].map((_, i) => (
          <div
            key={i}
            className="h-48 bg-muted animate-pulse rounded-lg"
          />
        ))}
      </div>
    )
  }

  if (!resources || resources.length === 0) {
    return (
      <div className="text-center py-12 border-2 border-dashed rounded-lg">
        <FileX className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
        <p className="text-muted-foreground text-lg">{emptyMessage}</p>
        <p className="text-muted-foreground text-sm mt-2">
          Resources will appear here once they are added
        </p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {resources.map((resource) => (
        <ResourceCard key={resource.ID} resource={resource} />
      ))}
    </div>
  )
}

export default ResourceList
