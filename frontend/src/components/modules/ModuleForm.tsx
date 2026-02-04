import React, { useState, useEffect } from 'react'
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
import { modulesApi } from '@/api/modules'
import { useToast } from '@/components/ui/toast'
import type { CreateModuleRequest, Module } from '@/types'

interface ModuleFormProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess: () => void
  mode?: 'create' | 'edit'
  initialData?: Module
}

const ModuleForm: React.FC<ModuleFormProps> = ({
  open,
  onOpenChange,
  onSuccess,
  mode = 'create',
  initialData,
}) => {
  const { showToast } = useToast()
  const [loading, setLoading] = useState(false)
  
  const [formData, setFormData] = useState<CreateModuleRequest>({
    code: '',
    name: '',
    department_name: '',
  })

  // Initialize form with initial data when in edit mode
  useEffect(() => {
    if (mode === 'edit' && initialData) {
      setFormData({
        code: initialData.Code,
        name: initialData.Name,
        department_name: initialData.DepartmentName,
      })
    } else {
      setFormData({ code: '', name: '', department_name: '' })
    }
  }, [mode, initialData, open])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      if (mode === 'edit' && initialData) {
        await modulesApi.updateModule(initialData.ID, formData)
        showToast('Module updated successfully', 'success')
      } else {
        await modulesApi.createModule(formData)
        showToast('Module created successfully', 'success')
      }
      onSuccess()
      onOpenChange(false)
      setFormData({ code: '', name: '', department_name: '' })
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : `Failed to ${mode} module`,
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
          <DialogTitle>{mode === 'edit' ? 'Edit Module' : 'Create Module'}</DialogTitle>
          <DialogDescription>
            {mode === 'edit' ? 'Update module information' : 'Add a new module to the system'}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="code">Module Code</Label>
              <Input
                id="code"
                placeholder="e.g., CS101"
                value={formData.code}
                onChange={(e) =>
                  setFormData({ ...formData, code: e.target.value })
                }
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="name">Module Name</Label>
              <Input
                id="name"
                placeholder="e.g., Introduction to Computer Science"
                value={formData.name}
                onChange={(e) =>
                  setFormData({ ...formData, name: e.target.value })
                }
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="department">Department Name</Label>
              <Input
                id="department"
                placeholder="e.g., Computer Science"
                value={formData.department_name}
                onChange={(e) =>
                  setFormData({ ...formData, department_name: e.target.value })
                }
                required
              />
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
              {loading ? (mode === 'edit' ? 'Updating...' : 'Creating...') : (mode === 'edit' ? 'Update' : 'Create')}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

export default ModuleForm
