import React, { useState } from 'react'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { FileText, Link as LinkIcon, StickyNote, Download, ExternalLink, Copy, Check } from 'lucide-react'
import { resourcesApi } from '@/api/resources'
import type { Resource, ResourceType } from '@/types'
import { Link } from 'react-router-dom'

interface ResourceCardProps {
  resource: Resource
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

const ResourceCard: React.FC<ResourceCardProps> = ({ resource }) => {
  const [copied, setCopied] = useState(false)

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
            </div>
          </div>
          {getResourceTypeBadge(resource.ResourceType)}
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
      <CardFooter className="pt-0">
        <div className="text-sm text-muted-foreground">
          Uploaded by{' '}
          <Link 
            to={`/users/${resource.UserID}`} 
            className="text-primary hover:underline font-medium"
          >
            {resource.UserName}
          </Link>
        </div>
      </CardFooter>
    </Card>
  )
}

export default ResourceCard
