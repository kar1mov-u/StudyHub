import React, { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { AlertTriangle } from 'lucide-react'
import { academicTermsApi } from '@/api/academicTerms'
import { useToast } from '@/components/ui/toast'
import type { CreateAcademicTermRequest } from '@/types'

interface AcademicTermFormProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess: () => void
}

const AcademicTermForm: React.FC<AcademicTermFormProps> = ({
  open,
  onOpenChange,
  onSuccess,
}) => {
  const { showToast } = useToast()
  const [loading, setLoading] = useState(false)
  const [showConfirmation, setShowConfirmation] = useState(false)
  const currentYear = new Date().getFullYear()
  
  const [formData, setFormData] = useState<CreateAcademicTermRequest>({
    year: currentYear,
    semester: 'fall',
  })

  const handleFormSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    // Show confirmation dialog before actually creating the term
    setShowConfirmation(true)
  }

  const handleConfirmedSubmit = async () => {
    setLoading(true)

    try {
      await academicTermsApi.createNewTerm(formData)
      showToast('New academic term created successfully! Module runs have been generated for all modules.', 'success')
      setShowConfirmation(false)
      onSuccess()
      onOpenChange(false)
      setFormData({ year: currentYear, semester: 'fall' })
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to create academic term',
        'error'
      )
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = () => {
    setShowConfirmation(false)
    setFormData({ year: currentYear, semester: 'fall' })
  }

  return (
    <>
      {/* Main Form Dialog */}
      <Dialog open={open && !showConfirmation} onOpenChange={onOpenChange}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Start New Academic Term</DialogTitle>
            <DialogDescription>
              Create a new semester and automatically generate module runs for all existing modules
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleFormSubmit}>
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="year">Year</Label>
                <Input
                  id="year"
                  type="number"
                  value={formData.year}
                  onChange={(e) =>
                    setFormData({ ...formData, year: parseInt(e.target.value) })
                  }
                  required
                  min={2000}
                  max={2100}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="semester">Semester</Label>
                <Select
                  id="semester"
                  value={formData.semester}
                  onChange={(e) =>
                    setFormData({ ...formData, semester: e.target.value })
                  }
                  required
                >
                  <option value="spring">Spring</option>
                  <option value="fall">Fall</option>
                </Select>
              </div>
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
              >
                Cancel
              </Button>
              <Button type="submit">
                Continue
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Confirmation Dialog */}
      <Dialog open={showConfirmation} onOpenChange={setShowConfirmation}>
        <DialogContent>
          <DialogHeader>
            <div className="flex items-center gap-3">
              <div className="h-10 w-10 rounded-full bg-amber-500/10 flex items-center justify-center">
                <AlertTriangle className="h-5 w-5 text-amber-500" />
              </div>
              <div>
                <DialogTitle>Confirm New Academic Term</DialogTitle>
              </div>
            </div>
            <DialogDescription className="pt-4 space-y-3">
              <p>
                You are about to create <strong className="capitalize">{formData.semester} {formData.year}</strong> as the new academic term.
              </p>
              <div className="bg-muted p-3 rounded-md space-y-2">
                <p className="font-medium text-sm">This will:</p>
                <ul className="list-disc list-inside text-sm space-y-1 text-muted-foreground">
                  <li>Set this term as the current active term</li>
                  <li>Create module runs for all existing modules</li>
                  <li>Generate weeks for each module run</li>
                </ul>
              </div>
              <p className="text-sm">
                Are you sure you want to proceed?
              </p>
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={handleCancel}
              disabled={loading}
            >
              Cancel
            </Button>
            <Button 
              onClick={handleConfirmedSubmit} 
              disabled={loading}
            >
              {loading ? 'Creating...' : 'Yes, Create Term'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}

export default AcademicTermForm
