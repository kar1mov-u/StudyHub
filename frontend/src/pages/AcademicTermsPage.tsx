import React, { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Plus, Loader2 } from 'lucide-react'
import AcademicTermCard from '@/components/academic-terms/AcademicTermCard'
import AcademicTermForm from '@/components/academic-terms/AcademicTermForm'
import { academicTermsApi } from '@/api/academicTerms'
import { useToast } from '@/components/ui/toast'
import type { AcademicTerm } from '@/types'

const AcademicTermsPage: React.FC = () => {
  const { showToast } = useToast()
  const [currentTerm, setCurrentTerm] = useState<AcademicTerm | null>(null)
  const [loading, setLoading] = useState(true)
  const [formOpen, setFormOpen] = useState(false)

  const loadCurrentTerm = async () => {
    try {
      setLoading(true)
      const data = await academicTermsApi.getCurrentAcademicTerm()
      setCurrentTerm(data)
    } catch (error) {
      // If no current term exists, that's okay - show empty state
      if (error instanceof Error && error.message.includes('404')) {
        setCurrentTerm(null)
      } else {
        showToast(
          error instanceof Error ? error.message : 'Failed to load current academic term',
          'error'
        )
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadCurrentTerm()
  }, [])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Academic Terms</h1>
          <p className="text-muted-foreground mt-2">
            Manage the current academic term and create new semesters
          </p>
        </div>
        <Button onClick={() => setFormOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Start New Term
        </Button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : currentTerm ? (
        <div>
          <div className="mb-4 p-4 bg-primary/10 border border-primary/20 rounded-lg">
            <p className="text-sm font-medium text-primary">
              Current Academic Term
            </p>
            <p className="text-xs text-muted-foreground mt-1">
              This is the active term for all modules
            </p>
          </div>
          <AcademicTermCard term={currentTerm} />
        </div>
      ) : (
        <div className="text-center py-12 border-2 border-dashed rounded-lg">
          <p className="text-muted-foreground">No academic term found</p>
          <p className="text-sm text-muted-foreground mt-2">
            Create the first academic term to start tracking module runs
          </p>
          <Button onClick={() => setFormOpen(true)} className="mt-4" variant="outline">
            <Plus className="h-4 w-4 mr-2" />
            Create first term
          </Button>
        </div>
      )}

      <AcademicTermForm
        open={formOpen}
        onOpenChange={setFormOpen}
        onSuccess={loadCurrentTerm}
      />
    </div>
  )
}

export default AcademicTermsPage
