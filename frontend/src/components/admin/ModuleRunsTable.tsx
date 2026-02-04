import React, { useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Trash2, Loader2 } from 'lucide-react'
import type { ModuleRun, Module } from '@/types'

interface ModuleRunWithModule extends ModuleRun {
  ModuleName?: string
  ModuleCode?: string
}

interface ModuleRunsTableProps {
  runs: ModuleRunWithModule[]
  modules: Module[]
  isLoading: boolean
  onDelete: (id: string) => void
  onRefresh: () => void
}

const ModuleRunsTable: React.FC<ModuleRunsTableProps> = ({
  runs,
  modules,
  isLoading,
  onDelete,
}) => {
  const [deleteRunId, setDeleteRunId] = useState<string | null>(null)

  const getModuleName = (moduleId: string) => {
    const module = modules.find((m) => m.ID === moduleId)
    return module ? `${module.Code} - ${module.Name}` : 'Unknown Module'
  }

  const handleDelete = () => {
    if (deleteRunId) {
      onDelete(deleteRunId)
      setDeleteRunId(null)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (runs.length === 0) {
    return (
      <div className="text-center py-12 border-2 border-dashed rounded-lg">
        <p className="text-muted-foreground">No module runs found</p>
      </div>
    )
  }

  return (
    <>
      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Module</TableHead>
              <TableHead>Year</TableHead>
              <TableHead>Semester</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {runs.map((run) => (
              <TableRow key={run.ID}>
                <TableCell className="font-medium">
                  {getModuleName(run.ModuleID)}
                </TableCell>
                <TableCell>{run.Year}</TableCell>
                <TableCell className="capitalize">{run.Semester}</TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setDeleteRunId(run.ID)}
                    >
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* Delete Confirmation Dialog */}
      <Dialog open={!!deleteRunId} onOpenChange={() => setDeleteRunId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Module Run</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this module run? This action cannot be
              undone and will affect all associated weeks and resources.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteRunId(null)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}

export default ModuleRunsTable
