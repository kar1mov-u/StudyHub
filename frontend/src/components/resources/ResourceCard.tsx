import React, { useState } from 'react'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { FileText, Link as LinkIcon, StickyNote, Download, ExternalLink, Copy, Check, Trash2 } from 'lucide-react'
import { resourcesApi } from '@/api/resources'
import type { Resource, ResourceType } from '@/types'
import { Link } from 'react-router-dom'

interface ResourceCardProps {
  resource: Resource
  currentUserId?: string
  onDelete?: (resourceId: string) => void
  selectable?: boolean
  selected?: boolean
  onToggleSelect?: (resourceId: string) => void
}

const getResourceIcon = (type: ResourceType) => {
  switch (type) {
    case 'file':
      return <FileText className="h-5 w-5" />
    case 'link':
      return <LinkIcon className="h-5 w-5" />
    case 'note':
      return <StickyNote className="h-5 w-5" />
    default:
      return <FileText className="h-5 w-5" />
  }
}

const getResourceTypeBadge = (type: ResourceType) => {
  switch (type) {
    case 'file':
      return <Badge variant="default">File</Badge>
    case 'link':
      return <Badge variant="outline">Link</Badge>
    case 'note':
      return <Badge variant="secondary">Note</Badge>
    default:
      return <Badge variant="default">{type}</Badge>
  }
}

const ResourceCard: React.FC<ResourceCardProps> = ({ resource, currentUserId, onDelete, selectable, selected, onToggleSelect }) => {
  const [copied, setCopied] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  
  const isOwnResource = currentUserId && currentUserId === resource.UserID

  const handleDownload = () => {
    if (resource.ObjectID) {
      resourcesApi.downloadResource(resource.ObjectID)
    }
  }

  const handleOpenLink = () => {
    if (resource.Url) {
      window.open(resource.Url, '_blank', 'noopener,noreferrer')
    }
  }

  const handleCopyLink = async () => {
    if (resource.Url) {
      try {
        await navigator.clipboard.writeText(resource.Url)
        setCopied(true)
        setTimeout(() => setCopied(false), 2000)
      } catch (err) {
        console.error('Failed to copy link:', err)
      }
    }
  }

  const handleDelete = async () => {
    if (!onDelete) return
    
    const confirmDelete = window.confirm(
      `Are you sure you want to delete "${resource.Name}"? This action cannot be undone.`
    )
    
    if (confirmDelete) {
      setIsDeleting(true)
      try {
        await resourcesApi.deleteResource(resource.ID)
        onDelete(resource.ID)
      } catch (err) {
        console.error('Failed to delete resource:', err)
        alert('Failed to delete resource. Please try again.')
      } finally {
        setIsDeleting(false)
      }
    }
  }

  const handleCardClick = () => {
    if (selectable && onToggleSelect) {
      onToggleSelect(resource.ID)
    }
  }

  return (
    <Card
      className={`hover:shadow-md transition-shadow ${selectable ? 'cursor-pointer' : ''} ${selected ? 'ring-2 ring-primary bg-primary/5' : ''}`}
      onClick={selectable ? handleCardClick : undefined}
    >
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-4">
          <div className="flex items-start gap-3 flex-1 min-w-0">
            {selectable && (
              <div className="mt-1 flex-shrink-0">
                <div
                  className={`h-5 w-5 rounded border-2 flex items-center justify-center transition-colors ${
                    selected
                      ? 'bg-primary border-primary text-primary-foreground'
                      : 'border-muted-foreground/40'
                  }`}
                >
                  {selected && <Check className="h-3 w-3" />}
                </div>
              </div>
            )}
            <div className="mt-1 text-primary">
              {getResourceIcon(resource.ResourceType)}
            </div>
            <div className="flex-1 min-w-0">
              <CardTitle className="text-lg truncate">
                {resource.Name || 'Untitled Resource'}
              </CardTitle>
            </div>
          </div>
          <div className="flex items-center gap-2">
            {!selectable && isOwnResource && onDelete && (
              <Button
                onClick={handleDelete}
                variant="ghost"
                size="sm"
                disabled={isDeleting}
                className="h-8 w-8 p-0 text-destructive hover:text-destructive hover:bg-destructive/10"
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            )}
            {getResourceTypeBadge(resource.ResourceType)}
          </div>
        </div>
      </CardHeader>
      {!selectable && (
        <CardContent>
          <div className="space-y-3">
            {resource.ResourceType === 'link' ? (
              <div className="flex gap-2">
                <Button
                  onClick={handleOpenLink}
                  className="flex-1"
                  variant="default"
                  disabled={!resource.Url}
                >
                  <ExternalLink className="h-4 w-4 mr-2" />
                  Open Link
                </Button>
                <Button
                  onClick={handleCopyLink}
                  variant="outline"
                  disabled={!resource.Url}
                >
                  {copied ? (
                    <Check className="h-4 w-4" />
                  ) : (
                    <Copy className="h-4 w-4" />
                  )}
                </Button>
              </div>
            ) : (
              <Button
                onClick={handleDownload}
                className="w-full"
                variant="default"
                disabled={!resource.ObjectID}
              >
                <Download className="h-4 w-4 mr-2" />
                Download
              </Button>
            )}
          </div>
        </CardContent>
      )}
      <CardFooter className="pt-0">
        <div className="text-sm text-muted-foreground">
          Uploaded by{' '}
          {selectable ? (
            <span className="text-primary font-medium">{resource.UserName}</span>
          ) : (
            <Link 
              to={`/users/${resource.UserID}`} 
              className="text-primary hover:underline font-medium"
            >
              {resource.UserName}
            </Link>
          )}
        </div>
      </CardFooter>
    </Card>
  )
}

export default ResourceCard
