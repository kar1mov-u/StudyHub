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
  const [terms, setTerms] = useState<AcademicTerm[]>([])
  const [loading, setLoading] = useState(true)
  const [formOpen, setFormOpen] = useState(false)

  const loadTerms = async () => {
    try {
      setLoading(true)
      const data = await academicTermsApi.listAcademicTerms()
      setTerms(data)
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to load academic terms',
        'error'
      )
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadTerms()
  }, [])

  const handleActivate = async (id: string) => {
    try {
      await academicTermsApi.activateAcademicTerm(id)
      showToast('Academic term activated successfully', 'success')
      loadTerms()
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to activate academic term',
        'error'
      )
    }
  }

  const handleDeactivate = async (id: string) => {
    try {
      await academicTermsApi.deactivateAcademicTerm(id)
      showToast('Academic term deactivated successfully', 'success')
      loadTerms()
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to deactivate academic term',
        'error'
      )
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Academic Terms</h1>
          <p className="text-muted-foreground mt-2">
            Manage academic terms and set the active term
          </p>
        </div>
        <Button onClick={() => setFormOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Create Term
        </Button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : terms.length === 0 ? (
        <div className="text-center py-12 border-2 border-dashed rounded-lg">
          <p className="text-muted-foreground">No academic terms found</p>
          <Button onClick={() => setFormOpen(true)} className="mt-4" variant="outline">
            <Plus className="h-4 w-4 mr-2" />
            Create your first term
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {terms.map((term) => (
            <AcademicTermCard
              key={term.ID}
              term={term}
              onActivate={handleActivate}
              onDeactivate={handleDeactivate}
            />
          ))}
        </div>
      )}

      <AcademicTermForm
        open={formOpen}
        onOpenChange={setFormOpen}
        onSuccess={loadTerms}
      />
    </div>
  )
}

export default AcademicTermsPage
