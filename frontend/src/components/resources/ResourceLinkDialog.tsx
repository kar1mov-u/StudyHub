import React, { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useToast } from '@/components/ui/toast'
import { resourcesApi } from '@/api/resources'
import { Link as LinkIcon, Loader2 } from 'lucide-react'

interface ResourceLinkDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  weekId: string
  onSuccess: () => void
}

const ResourceLinkDialog: React.FC<ResourceLinkDialogProps> = ({
  open,
  onOpenChange,
  weekId,
  onSuccess,
}) => {
  const { showToast } = useToast()
  const [name, setName] = useState('')
  const [url, setUrl] = useState('')
  const [uploading, setUploading] = useState(false)

  const validateUrl = (urlString: string): boolean => {
    try {
      new URL(urlString)
      return true
    } catch {
      return false
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!name.trim()) {
      showToast('Please enter a name for the link', 'error')
      return
    }

    if (!url.trim()) {
      showToast('Please enter a URL', 'error')
      return
    }

    if (!validateUrl(url)) {
      showToast('Please enter a valid URL (e.g., https://example.com)', 'error')
      return
    }

    try {
      setUploading(true)
      await resourcesApi.uploadLink(weekId, name.trim(), url.trim())
      showToast('Link added successfully', 'success')
      onSuccess()
      handleClose()
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to add link',
        'error'
      )
    } finally {
      setUploading(false)
    }
  }

  const handleClose = () => {
    setName('')
    setUrl('')
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Link</DialogTitle>
          <DialogDescription>
            Add a link resource to this week. You can share links to external resources, documents, or websites.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="link-name">
                Name <span className="text-destructive">*</span>
              </Label>
              <Input
                id="link-name"
                type="text"
                placeholder="e.g., Course Textbook PDF"
                value={name}
                onChange={(e) => setName(e.target.value)}
                disabled={uploading}
                maxLength={200}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="link-url">
                URL <span className="text-destructive">*</span>
              </Label>
              <Input
                id="link-url"
                type="url"
                placeholder="https://example.com"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                disabled={uploading}
              />
              <p className="text-xs text-muted-foreground">
                Include the full URL starting with http:// or https://
              </p>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" type="button" onClick={handleClose} disabled={uploading}>
              Cancel
            </Button>
            <Button type="submit" disabled={uploading}>
              {uploading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Adding...
                </>
              ) : (
                <>
                  <LinkIcon className="mr-2 h-4 w-4" />
                  Add Link
                </>
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

export default ResourceLinkDialog
