import React, { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
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
import { Pencil, Trash2, Loader2 } from 'lucide-react'
import ModuleForm from '@/components/modules/ModuleForm'
import type { Module } from '@/types'

interface ModulesTableProps {
  modules: Module[]
  isLoading: boolean
  onDelete: (id: string) => void
  onRefresh: () => void
}

const ModulesTable: React.FC<ModulesTableProps> = ({
  modules,
  isLoading,
  onDelete,
  onRefresh,
}) => {
  const [editModule, setEditModule] = useState<Module | null>(null)
  const [deleteModuleId, setDeleteModuleId] = useState<string | null>(null)

  const handleDelete = () => {
    if (deleteModuleId) {
      onDelete(deleteModuleId)
      setDeleteModuleId(null)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (modules.length === 0) {
    return (
      <div className="text-center py-12 border-2 border-dashed rounded-lg">
        <p className="text-muted-foreground">No modules found</p>
      </div>
    )
  }

  return (
    <>
      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Code</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Department</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {modules.map((module) => (
              <TableRow key={module.ID}>
                <TableCell className="font-medium">{module.Code}</TableCell>
                <TableCell>{module.Name}</TableCell>
                <TableCell>
                  <Badge variant="outline">{module.DepartmentName}</Badge>
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setEditModule(module)}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setDeleteModuleId(module.ID)}
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

      {/* Edit Dialog */}
      {editModule && (
        <ModuleForm
          open={!!editModule}
          onOpenChange={(open) => !open && setEditModule(null)}
          onSuccess={() => {
            onRefresh()
            setEditModule(null)
          }}
          mode="edit"
          initialData={editModule}
        />
      )}

      {/* Delete Confirmation Dialog */}
      <Dialog open={!!deleteModuleId} onOpenChange={() => setDeleteModuleId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Module</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this module? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteModuleId(null)}>
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

export default ModulesTable
