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
import { Upload, Loader2, FileIcon } from 'lucide-react'

interface ResourceUploadDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  weekId: string
  onSuccess: () => void
}

const ResourceUploadDialog: React.FC<ResourceUploadDialogProps> = ({
  open,
  onOpenChange,
  weekId,
  onSuccess,
}) => {
  const { showToast } = useToast()
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [uploading, setUploading] = useState(false)

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      // Simple size validation: 50MB max
      if (file.size > 50 * 1024 * 1024) {
        showToast('File size must be less than 50MB', 'error')
        return
      }
      setSelectedFile(file)
    }
  }

  const handleUpload = async () => {
    if (!selectedFile) return

    try {
      setUploading(true)
      const formData = new FormData()
      formData.append('file', selectedFile)

      await resourcesApi.uploadFile(weekId, formData)
      showToast('File uploaded successfully', 'success')
      onSuccess()
      handleClose()
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to upload file',
        'error'
      )
    } finally {
      setUploading(false)
    }
  }

  const handleClose = () => {
    setSelectedFile(null)
    onOpenChange(false)
  }

  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Upload Resource</DialogTitle>
          <DialogDescription>
            Select a file to upload to this week. Maximum file size: 50MB
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="file">Choose File</Label>
            <Input
              id="file"
              type="file"
              onChange={handleFileSelect}
              disabled={uploading}
              accept="*/*"
            />
          </div>

          {selectedFile && (
            <div className="flex items-center gap-3 p-3 bg-muted rounded-lg">
              <FileIcon className="h-8 w-8 text-primary" />
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium truncate">
                  {selectedFile.name}
                </p>
                <p className="text-xs text-muted-foreground">
                  {formatFileSize(selectedFile.size)}
                </p>
              </div>
            </div>
          )}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={handleClose} disabled={uploading}>
            Cancel
          </Button>
          <Button
            onClick={handleUpload}
            disabled={!selectedFile || uploading}
          >
            {uploading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Uploading...
              </>
            ) : (
              <>
                <Upload className="mr-2 h-4 w-4" />
                Upload
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export default ResourceUploadDialog
