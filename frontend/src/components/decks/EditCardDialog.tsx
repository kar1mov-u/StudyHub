import React, { useState, useEffect } from 'react'
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
import type { UserDeckCard } from '@/types'

interface EditCardDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  card: UserDeckCard | null
  onSubmit: (cardId: string, front?: string, back?: string) => Promise<void>
}

const EditCardDialog: React.FC<EditCardDialogProps> = ({
  open,
  onOpenChange,
  card,
  onSubmit,
}) => {
  const [front, setFront] = useState('')
  const [back, setBack] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState('')

  // Populate form when card changes
  useEffect(() => {
    if (card) {
      setFront(card.Front)
      setBack(card.Back)
    }
  }, [card])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!card) return

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

    // Check if anything changed
    const frontChanged = front.trim() !== card.Front
    const backChanged = back.trim() !== card.Back

    if (!frontChanged && !backChanged) {
      onOpenChange(false)
      return
    }

    try {
      setIsSubmitting(true)
      await onSubmit(
        card.ID,
        frontChanged ? front.trim() : undefined,
        backChanged ? back.trim() : undefined
      )
      onOpenChange(false)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update card')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCancel = () => {
    if (card) {
      setFront(card.Front)
      setBack(card.Back)
    }
    setError('')
    onOpenChange(false)
  }

  if (!card) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Edit Flashcard</DialogTitle>
            <DialogDescription>
              Update the question and answer for this flashcard.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="edit-front">Question (Front)</Label>
              <Input
                id="edit-front"
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
              <Label htmlFor="edit-back">Answer (Back)</Label>
              <Input
                id="edit-back"
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
              Save Changes
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

export default EditCardDialog
