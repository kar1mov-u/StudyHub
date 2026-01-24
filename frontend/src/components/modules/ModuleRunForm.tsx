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
import { moduleRunsApi } from '@/api/moduleRuns'
import { useToast } from '@/components/ui/toast'
import type { CreateModuleRunRequest } from '@/types'

interface ModuleRunFormProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess: () => void
  moduleId: string
}

const ModuleRunForm: React.FC<ModuleRunFormProps> = ({
  open,
  onOpenChange,
  onSuccess,
  moduleId,
}) => {
  const { showToast } = useToast()
  const [loading, setLoading] = useState(false)
  const currentYear = new Date().getFullYear()
  
  const [formData, setFormData] = useState<CreateModuleRunRequest>({
    year: currentYear,
    semester: 'fall',
    is_active: false,
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      await moduleRunsApi.createModuleRun(moduleId, formData)
      showToast('Module run created successfully', 'success')
      onSuccess()
      onOpenChange(false)
      setFormData({ year: currentYear, semester: 'fall', is_active: false })
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to create module run',
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
          <DialogTitle>Create Module Run</DialogTitle>
          <DialogDescription>
            Add a new run for this module
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
            <div className="flex items-center space-x-2">
              <input
                type="checkbox"
                id="is_active"
                checked={formData.is_active}
                onChange={(e) =>
                  setFormData({ ...formData, is_active: e.target.checked })
                }
                className="h-4 w-4 rounded border-gray-300"
              />
              <Label htmlFor="is_active" className="cursor-pointer">
                Set as active run
              </Label>
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

export default ModuleRunForm
