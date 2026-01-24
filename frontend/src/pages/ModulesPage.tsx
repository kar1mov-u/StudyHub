import React, { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Plus, Loader2 } from 'lucide-react'
import ModuleCard from '@/components/modules/ModuleCard'
import ModuleForm from '@/components/modules/ModuleForm'
import { modulesApi } from '@/api/modules'
import { useToast } from '@/components/ui/toast'
import type { Module } from '@/types'

const ModulesPage: React.FC = () => {
  const { showToast } = useToast()
  const [modules, setModules] = useState<Module[]>([])
  const [loading, setLoading] = useState(true)
  const [formOpen, setFormOpen] = useState(false)

  const loadModules = async () => {
    try {
      setLoading(true)
      const data = await modulesApi.listModules()
      setModules(data)
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to load modules',
        'error'
      )
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadModules()
  }, [])

  const handleDelete = async (id: string) => {
    try {
      await modulesApi.deleteModule(id)
      showToast('Module deleted successfully', 'success')
      loadModules()
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to delete module',
        'error'
      )
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Modules</h1>
          <p className="text-muted-foreground mt-2">
            Manage your course modules and their runs
          </p>
        </div>
        <Button onClick={() => setFormOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Create Module
        </Button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : modules.length === 0 ? (
        <div className="text-center py-12 border-2 border-dashed rounded-lg">
          <p className="text-muted-foreground">No modules found</p>
          <Button onClick={() => setFormOpen(true)} className="mt-4" variant="outline">
            <Plus className="h-4 w-4 mr-2" />
            Create your first module
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {modules.map((module) => (
            <ModuleCard
              key={module.ID}
              module={module}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}

      <ModuleForm
        open={formOpen}
        onOpenChange={setFormOpen}
        onSuccess={loadModules}
      />
    </div>
  )
}

export default ModulesPage
