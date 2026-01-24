import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Calendar,
  BookOpen,
  Plus,
  Loader2,
  ArrowRight,
  CheckCircle2,
} from 'lucide-react'
import { academicTermsApi } from '@/api/academicTerms'
import { modulesApi } from '@/api/modules'
import { useToast } from '@/components/ui/toast'
import type { AcademicTerm, Module } from '@/types'

const HomePage: React.FC = () => {
  const navigate = useNavigate()
  const { showToast } = useToast()
  const [loading, setLoading] = useState(true)
  const [activeTerm, setActiveTerm] = useState<AcademicTerm | null>(null)
  const [modules, setModules] = useState<Module[]>([])

  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true)
        const [termData, modulesData] = await Promise.all([
          academicTermsApi.getActiveAcademicTerm().catch(() => null),
          modulesApi.listModules(),
        ])
        setActiveTerm(termData)
        setModules(modulesData)
      } catch (error) {
        showToast(
          error instanceof Error ? error.message : 'Failed to load data',
          'error'
        )
      } finally {
        setLoading(false)
      }
    }

    loadData()
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground mt-2">
          Welcome to StudyHub - Manage your academic modules
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Active Academic Term
            </CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {activeTerm ? (
              <div>
                <div className="text-2xl font-bold capitalize">
                  {activeTerm.Semester} {activeTerm.Year}
                </div>
                <div className="flex items-center gap-1 mt-2">
                  <CheckCircle2 className="h-4 w-4 text-green-500" />
                  <p className="text-xs text-muted-foreground">Active</p>
                </div>
              </div>
            ) : (
              <div>
                <p className="text-sm text-muted-foreground">No active term</p>
                <Button
                  variant="link"
                  size="sm"
                  className="px-0 mt-2"
                  onClick={() => navigate('/academic-terms')}
                >
                  Create one
                  <ArrowRight className="h-4 w-4 ml-1" />
                </Button>
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Modules</CardTitle>
            <BookOpen className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{modules.length}</div>
            <Button
              variant="link"
              size="sm"
              className="px-0 mt-2"
              onClick={() => navigate('/modules')}
            >
              View all
              <ArrowRight className="h-4 w-4 ml-1" />
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Quick Actions</CardTitle>
            <Plus className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent className="space-y-2">
            <Button
              variant="outline"
              size="sm"
              className="w-full justify-start"
              onClick={() => navigate('/modules')}
            >
              <Plus className="h-4 w-4 mr-2" />
              New Module
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="w-full justify-start"
              onClick={() => navigate('/academic-terms')}
            >
              <Plus className="h-4 w-4 mr-2" />
              New Term
            </Button>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Recent Modules</CardTitle>
            <Button variant="link" onClick={() => navigate('/modules')}>
              View all
              <ArrowRight className="h-4 w-4 ml-1" />
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {modules.length === 0 ? (
            <div className="text-center py-8 border-2 border-dashed rounded-lg">
              <p className="text-muted-foreground">No modules yet</p>
              <Button
                variant="outline"
                size="sm"
                className="mt-4"
                onClick={() => navigate('/modules')}
              >
                <Plus className="h-4 w-4 mr-2" />
                Create your first module
              </Button>
            </div>
          ) : (
            <div className="space-y-3">
              {modules.slice(0, 5).map((module) => (
                <div
                  key={module.ID}
                  className="flex items-center justify-between p-3 border rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
                  onClick={() => navigate(`/modules/${module.ID}`)}
                >
                  <div className="flex items-center gap-3">
                    <BookOpen className="h-5 w-5 text-primary" />
                    <div>
                      <p className="font-medium">{module.Code}</p>
                      <p className="text-sm text-muted-foreground">
                        {module.Name}
                      </p>
                    </div>
                  </div>
                  <Badge variant="outline">{module.DepartmentName}</Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

export default HomePage
