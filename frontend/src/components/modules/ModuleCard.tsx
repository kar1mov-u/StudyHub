import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { BookOpen, Trash2 } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useAuth } from '@/context/AuthContext'
import type { Module } from '@/types'

interface ModuleCardProps {
  module: Module
  onDelete: (id: string) => void
}

const ModuleCard: React.FC<ModuleCardProps> = ({ module, onDelete }) => {
  const navigate = useNavigate()
  const { user } = useAuth()
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const handleDelete = () => {
    onDelete(module.ID)
    setDeleteDialogOpen(false)
  }

  return (
    <>
      <Card className="hover:shadow-lg transition-shadow cursor-pointer">
        <CardHeader onClick={() => navigate(`/modules/${module.ID}`)}>
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-2">
              <BookOpen className="h-5 w-5 text-primary" />
              <CardTitle className="text-lg">{module.Code}</CardTitle>
            </div>
          </div>
          <h3 className="text-sm font-semibold text-gray-900 mt-2">
            {module.Name}
          </h3>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <Badge variant="outline">{module.DepartmentName}</Badge>
            {user?.IsAdmin && (
              <Button
                variant="ghost"
                size="sm"
                onClick={(e) => {
                  e.stopPropagation()
                  setDeleteDialogOpen(true)
                }}
              >
                <Trash2 className="h-4 w-4 text-destructive" />
              </Button>
            )}
          </div>
        </CardContent>
      </Card>

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Module</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {module.Code} - {module.Name}? This
              action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteDialogOpen(false)}>
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

export default ModuleCard
