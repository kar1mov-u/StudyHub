import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ArrowLeft, Plus, Loader2, Trash2, Calendar } from 'lucide-react'
import ModuleRunForm from '@/components/modules/ModuleRunForm'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { modulesApi } from '@/api/modules'
import { moduleRunsApi } from '@/api/moduleRuns'
import { useToast } from '@/components/ui/toast'
import { useAuth } from '@/context/AuthContext'
import type { ModulePage } from '@/types'

const ModuleDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { showToast } = useToast()
  const { user } = useAuth()
  
  const [modulePage, setModulePage] = useState<ModulePage | null>(null)
  const [loading, setLoading] = useState(true)
  const [formOpen, setFormOpen] = useState(false)
  const [deleteRunId, setDeleteRunId] = useState<string | null>(null)

  const loadModuleData = async () => {
    if (!id) return
    
    try {
      setLoading(true)
      const data = await modulesApi.getModuleFull(id)
      setModulePage(data)
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to load module',
        'error'
      )
      navigate('/modules')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadModuleData()
  }, [id])

  const handleDeleteRun = async () => {
    if (!deleteRunId) return

    try {
      await moduleRunsApi.deleteModuleRun(deleteRunId)
      showToast('Module run deleted successfully', 'success')
      setDeleteRunId(null)
      loadModuleData()
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to delete module run',
        'error'
      )
    }
  }

  const handleWeekClick = (weekId: string) => {
    navigate(`/modules/${id}/weeks/${weekId}`)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!modulePage) {
    return null
  }

  const { Module: module, Run: activeRun, Weeks: weeks } = modulePage

  return (
    <div className="space-y-6">
      <div>
        <Button
          variant="ghost"
          onClick={() => navigate('/modules')}
          className="mb-4"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Modules
        </Button>
        
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">{module.Code}</h1>
            <p className="text-xl text-muted-foreground mt-2">{module.Name}</p>
            <Badge variant="outline" className="mt-2">
              {module.DepartmentName}
            </Badge>
          </div>
          {user?.IsAdmin && (
            <Button onClick={() => setFormOpen(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Add Run
            </Button>
          )}
        </div>
      </div>

      {activeRun && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Calendar className="h-5 w-5" />
              Active Run
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <p className="text-lg capitalize">
                {activeRun.Semester} {activeRun.Year}
              </p>
              <div className="flex items-center gap-2">
                <Badge variant="success">Active</Badge>
                {user?.IsAdmin && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setDeleteRunId(activeRun.ID)}
                  >
                    <Trash2 className="h-4 w-4 text-destructive" />
                  </Button>
                )}
              </div>
              {weeks && weeks.length > 0 && (
                <div className="mt-4">
                  <p className="text-sm font-medium text-muted-foreground mb-3">
                    Weeks: {weeks.length}
                  </p>
                  <div className="flex flex-wrap gap-2">
                    {weeks.map((week) => (
                      <Badge 
                        key={week.ID} 
                        variant="outline" 
                        className="cursor-pointer hover:bg-primary hover:text-primary-foreground transition-colors"
                        onClick={() => handleWeekClick(week.ID)}
                      >
                        Week {week.Number}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      )}

      {!activeRun && (
        <div className="text-center py-12 border-2 border-dashed rounded-lg">
          <p className="text-muted-foreground">No active run for this module</p>
          {user?.IsAdmin && (
            <Button onClick={() => setFormOpen(true)} className="mt-4" variant="outline">
              <Plus className="h-4 w-4 mr-2" />
              Create a run
            </Button>
          )}
        </div>
      )}

      {user?.IsAdmin && (
        <ModuleRunForm
          open={formOpen}
          onOpenChange={setFormOpen}
          onSuccess={loadModuleData}
          moduleId={id!}
        />
      )}

      <Dialog open={!!deleteRunId} onOpenChange={() => setDeleteRunId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Module Run</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this module run? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteRunId(null)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDeleteRun}>
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default ModuleDetailPage
