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
  const currentYear = new Date().getFullYear()
  
  const [formData, setFormData] = useState<CreateAcademicTermRequest>({
    year: currentYear,
    semester: 'fall',
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      await academicTermsApi.createAcademicTerm(formData)
      showToast('Academic term created successfully', 'success')
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

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Academic Term</DialogTitle>
          <DialogDescription>
            Add a new academic term to the system
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
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
            <Button type="submit" disabled={loading}>
              {loading ? 'Creating...' : 'Create'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

export default AcademicTermForm
