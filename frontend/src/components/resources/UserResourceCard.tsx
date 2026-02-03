import React, { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { FileText, Link as LinkIcon, StickyNote, Download, ExternalLink, Copy, Check, Trash2 } from 'lucide-react'
import { resourcesApi } from '@/api/resources'
import type { UserResource, ResourceType } from '@/types'

interface UserResourceCardProps {
  resource: UserResource
  showDelete?: boolean
  onDelete?: (resourceId: string) => void
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

const UserResourceCard: React.FC<UserResourceCardProps> = ({ resource, showDelete = false, onDelete }) => {
  const [copied, setCopied] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)

  const handleDownload = () => {
    if (resource.ObjectID) {
      resourcesApi.downloadResource(resource.ObjectID)
    }
  }

  const handleOpenLink = () => {
    if (resource.ExternalLink) {
      window.open(resource.ExternalLink, '_blank', 'noopener,noreferrer')
    }
  }

  const handleCopyLink = async () => {
    if (resource.ExternalLink) {
      try {
        await navigator.clipboard.writeText(resource.ExternalLink)
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

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-4">
          <div className="flex items-start gap-3 flex-1 min-w-0">
            <div className="mt-1 text-primary">
              {getResourceIcon(resource.ResourceType)}
            </div>
            <div className="flex-1 min-w-0">
              <CardTitle className="text-lg truncate">
                {resource.Name || 'Untitled Resource'}
              </CardTitle>
              <div className="text-sm text-muted-foreground mt-1">
                <div>{resource.ModuleName}</div>
                <div>
                  {resource.Semester} {resource.Year}, Week {resource.WeekNumber}
                </div>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-2">
            {showDelete && (
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
      <CardContent>
        <div className="space-y-3">
          {resource.ResourceType === 'link' ? (
            <div className="flex gap-2">
              <Button
                onClick={handleOpenLink}
                className="flex-1"
                variant="default"
                disabled={!resource.ExternalLink}
              >
                <ExternalLink className="h-4 w-4 mr-2" />
                Open Link
              </Button>
              <Button
                onClick={handleCopyLink}
                variant="outline"
                disabled={!resource.ExternalLink}
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
    </Card>
  )
}

export default UserResourceCard
