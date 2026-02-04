import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Select } from '@/components/ui/select'
import { Label } from '@/components/ui/label'
import { ArrowLeft, Loader2, Trash2, Calendar } from 'lucide-react'
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
import type { Module, ModuleRun, Week } from '@/types'

const ModuleDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { showToast } = useToast()
  const { user } = useAuth()
  
  const [module, setModule] = useState<Module | null>(null)
  const [allRuns, setAllRuns] = useState<ModuleRun[]>([])
  const [selectedRunId, setSelectedRunId] = useState<string | null>(null)
  const [weeks, setWeeks] = useState<Week[]>([])
  const [loading, setLoading] = useState(true)
  const [weeksLoading, setWeeksLoading] = useState(false)
  const [deleteRunId, setDeleteRunId] = useState<string | null>(null)

  // Get the selected run object
  const selectedRun = allRuns.find((run) => run.ID === selectedRunId)

  const loadModuleData = async () => {
    if (!id) return
    
    try {
      setLoading(true)
      // Load module info and active run (default)
      const modulePage = await modulesApi.getModuleFull(id)
      setModule(modulePage.Module)
      
      // Load all runs for this module
      const runs = await moduleRunsApi.listModuleRuns(id)
      setAllRuns(runs)
      
      // Set default to active run if it exists, otherwise first run
      if (modulePage.Run) {
        setSelectedRunId(modulePage.Run.ID)
        setWeeks(modulePage.Weeks || [])
      } else if (runs.length > 0) {
        setSelectedRunId(runs[0].ID)
        // Load weeks for first run
        await loadWeeksForRun(runs[0].ID)
      }
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

  const loadWeeksForRun = async (runId: string) => {
    try {
      setWeeksLoading(true)
      const runPage = await moduleRunsApi.getModuleRun(runId)
      setWeeks(runPage.Weeks || [])
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to load weeks',
        'error'
      )
      setWeeks([])
    } finally {
      setWeeksLoading(false)
    }
  }

  useEffect(() => {
    loadModuleData()
  }, [id])

  const handleRunChange = async (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newRunId = e.target.value
    setSelectedRunId(newRunId)
    await loadWeeksForRun(newRunId)
  }

  const handleDeleteRun = async () => {
    if (!deleteRunId) return

    try {
      await moduleRunsApi.deleteModuleRun(deleteRunId)
      showToast('Module run deleted successfully', 'success')
      setDeleteRunId(null)
      // If we deleted the selected run, reset selection
      if (deleteRunId === selectedRunId) {
        setSelectedRunId(null)
        setWeeks([])
      }
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

  if (!module) {
    return null
  }

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
        </div>
      </div>

      {allRuns.length > 0 ? (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Calendar className="h-5 w-5" />
              Module Run
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {/* Run Selection Dropdown */}
              <div className="space-y-2">
                <Label htmlFor="run-select">Select Module Run</Label>
                <div className="flex items-center gap-2">
                  <Select
                    id="run-select"
                    value={selectedRunId || ''}
                    onChange={handleRunChange}
                    className="flex-1"
                  >
                    {allRuns.map((run) => (
                      <option key={run.ID} value={run.ID}>
                        {run.Semester.charAt(0).toUpperCase() + run.Semester.slice(1)} {run.Year}
                      </option>
                    ))}
                  </Select>
                  {user?.IsAdmin && selectedRun && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setDeleteRunId(selectedRun.ID)}
                    >
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  )}
                </div>
              </div>

              {/* Selected Run Info */}
              {selectedRun && (
                <div className="flex items-center gap-2">
                  <span className="text-sm text-muted-foreground capitalize">
                    {selectedRun.Semester} {selectedRun.Year}
                  </span>
                </div>
              )}

              {/* Weeks Display */}
              {weeksLoading ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-primary" />
                </div>
              ) : weeks.length > 0 ? (
                <div>
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
              ) : (
                <p className="text-sm text-muted-foreground">No weeks found for this run</p>
              )}
            </div>
          </CardContent>
        </Card>
      ) : (
        <div className="text-center py-12 border-2 border-dashed rounded-lg">
          <p className="text-muted-foreground">No runs found for this module</p>
          <p className="text-sm text-muted-foreground mt-2">
            Module runs are created automatically when a new academic term is started
          </p>
        </div>
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
