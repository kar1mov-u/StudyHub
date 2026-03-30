import React, { useState } from 'react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Loader2 } from 'lucide-react'

interface CreateCustomCardDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSubmit: (front: string, back: string) => Promise<void>
}

const CreateCustomCardDialog: React.FC<CreateCustomCardDialogProps> = ({
  open,
  onOpenChange,
  onSubmit,
}) => {
  const [front, setFront] = useState('')
  const [back, setBack] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    // Validation
    if (!front.trim()) {
      setError('Question cannot be empty')
      return
    }
    if (!back.trim()) {
      setError('Answer cannot be empty')
      return
    }

    try {
      setIsSubmitting(true)
      await onSubmit(front.trim(), back.trim())
      // Reset form on success
      setFront('')
      setBack('')
      onOpenChange(false)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create card')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCancel = () => {
    setFront('')
    setBack('')
    setError('')
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Create Custom Flashcard</DialogTitle>
            <DialogDescription>
              Create your own flashcard with a question and answer.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="front">Question (Front)</Label>
              <Input
                id="front"
                placeholder="Enter the question..."
                value={front}
                onChange={(e) => setFront(e.target.value)}
                disabled={isSubmitting}
                autoFocus
              />
              <p className="text-xs text-muted-foreground">
                {front.length} characters
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="back">Answer (Back)</Label>
              <Input
                id="back"
                placeholder="Enter the answer..."
                value={back}
                onChange={(e) => setBack(e.target.value)}
                disabled={isSubmitting}
              />
              <p className="text-xs text-muted-foreground">
                {back.length} characters
              </p>
            </div>

            {error && (
              <p className="text-sm text-destructive">{error}</p>
            )}
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={handleCancel}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
              Create Card
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

export default CreateCustomCardDialog
