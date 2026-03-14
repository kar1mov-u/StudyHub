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
  FileText,
  Link2,
  Upload,
  Clock,
} from 'lucide-react'
import { academicTermsApi } from '@/api/academicTerms'
import { modulesApi } from '@/api/modules'
import { resourcesApi } from '@/api/resources'
import { useAuth } from '@/context/AuthContext'
import { useToast } from '@/components/ui/toast'
import type { AcademicTerm, Module, UserResource } from '@/types'

function timeAgo(dateString: string): string {
  const now = new Date()
  const date = new Date(dateString)
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  if (seconds < 60) return 'just now'
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}d ago`
  const weeks = Math.floor(days / 7)
  if (weeks < 4) return `${weeks}w ago`
  const months = Math.floor(days / 30)
  if (months < 12) return `${months}mo ago`
  const years = Math.floor(days / 365)
  return `${years}y ago`
}

const HomePage: React.FC = () => {
  const navigate = useNavigate()
  const { user } = useAuth()
  const { showToast } = useToast()
  const [loading, setLoading] = useState(true)
  const [activeTerm, setActiveTerm] = useState<AcademicTerm | null>(null)
  const [modules, setModules] = useState<Module[]>([])
  const [myResources, setMyResources] = useState<UserResource[]>([])

  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true)
        const promises: [
          Promise<AcademicTerm | null>,
          Promise<Module[]>,
          Promise<UserResource[]>,
        ] = [
          academicTermsApi.getCurrentAcademicTerm().catch(() => null),
          modulesApi.listModules(),
          user
            ? resourcesApi.getUserResources(user.ID).catch(() => [])
            : Promise.resolve([]),
        ]
        const [termData, modulesData, resourcesData] = await Promise.all(promises)
        setActiveTerm(termData)
        setModules(modulesData)
        setMyResources(resourcesData)
      } catch (error) {
        showToast(
          error instanceof Error ? error.message : 'Failed to load dashboard data',
          'error'
        )
      } finally {
        setLoading(false)
      }
    }

    loadData()
  }, [user])

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  const fileCount = myResources.filter((r) => r.ResourceType === 'file').length
  const linkCount = myResources.filter((r) => r.ResourceType === 'link').length

  const recentResources = [...myResources]
    .sort(
      (a, b) =>
        new Date(b.CreatedAt).getTime() - new Date(a.CreatedAt).getTime()
    )
    .slice(0, 5)

  return (
    <div className="space-y-6">
      {/* Greeting */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">
          Welcome back, {user?.FirstName || 'Student'}!
        </h1>
        <p className="text-muted-foreground mt-1">
          {activeTerm ? (
            <span className="flex items-center gap-1.5">
              <CheckCircle2 className="h-4 w-4 text-green-500" />
              <span className="capitalize">
                {activeTerm.Semester} {activeTerm.Year}
              </span>
              <span>- Active Term</span>
            </span>
          ) : (
            'No active academic term'
          )}
        </p>
      </div>

      {/* Stats Row */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* My Uploads */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">My Uploads</CardTitle>
            <Upload className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{myResources.length}</div>
            <p className="text-xs text-muted-foreground mt-1">
              {fileCount} {fileCount === 1 ? 'file' : 'files'}, {linkCount}{' '}
              {linkCount === 1 ? 'link' : 'links'}
            </p>
          </CardContent>
        </Card>

        {/* Total Modules */}
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
              className="px-0 mt-1 h-auto"
              onClick={() => navigate('/modules')}
            >
              View all
              <ArrowRight className="h-3 w-3 ml-1" />
            </Button>
          </CardContent>
        </Card>

        {/* Active Term */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Term</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {activeTerm ? (
              <div>
                <div className="text-2xl font-bold capitalize">
                  {activeTerm.Semester} {activeTerm.Year}
                </div>
                <div className="flex items-center gap-1 mt-1">
                  <CheckCircle2 className="h-3 w-3 text-green-500" />
                  <p className="text-xs text-muted-foreground">Active</p>
                </div>
              </div>
            ) : (
              <div>
                <p className="text-sm text-muted-foreground">No active term</p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* My Recent Uploads */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <Clock className="h-5 w-5" />
              My Recent Uploads
            </CardTitle>
            {myResources.length > 0 && user && (
              <Button
                variant="link"
                size="sm"
                onClick={() => navigate(`/users/${user.ID}`)}
              >
                View all
                <ArrowRight className="h-4 w-4 ml-1" />
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent>
          {recentResources.length === 0 ? (
            <div className="text-center py-8 border-2 border-dashed rounded-lg">
              <Upload className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
              <p className="text-muted-foreground">No uploads yet</p>
              <p className="text-xs text-muted-foreground mt-1">
                Head to a module week to upload files or add links
              </p>
              <Button
                variant="outline"
                size="sm"
                className="mt-4"
                onClick={() => navigate('/modules')}
              >
                Browse Modules
              </Button>
            </div>
          ) : (
            <div className="space-y-3">
              {recentResources.map((resource) => (
                <div
                  key={resource.ID}
                  className="flex items-center justify-between p-3 border rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
                  onClick={() => navigate(`/modules`)}
                >
                  <div className="flex items-center gap-3">
                    {resource.ResourceType === 'file' ? (
                      <div className="h-9 w-9 rounded-lg bg-blue-50 flex items-center justify-center">
                        <FileText className="h-4 w-4 text-blue-600" />
                      </div>
                    ) : (
                      <div className="h-9 w-9 rounded-lg bg-purple-50 flex items-center justify-center">
                        <Link2 className="h-4 w-4 text-purple-600" />
                      </div>
                    )}
                    <div>
                      <p className="font-medium text-sm">{resource.Name}</p>
                      <p className="text-xs text-muted-foreground">
                        {resource.ModuleName} - Week {resource.WeekNumber}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge
                      variant={
                        resource.ResourceType === 'file' ? 'default' : 'secondary'
                      }
                      className="text-xs"
                    >
                      {resource.ResourceType}
                    </Badge>
                    <span className="text-xs text-muted-foreground whitespace-nowrap">
                      {timeAgo(resource.CreatedAt)}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Browse Modules */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <BookOpen className="h-5 w-5" />
              Browse Modules
            </CardTitle>
            {modules.length > 6 && (
              <Button
                variant="link"
                size="sm"
                onClick={() => navigate('/modules')}
              >
                View all
                <ArrowRight className="h-4 w-4 ml-1" />
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent>
          {modules.length === 0 ? (
            <div className="text-center py-8 border-2 border-dashed rounded-lg">
              <BookOpen className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
              <p className="text-muted-foreground">No modules yet</p>
              {user?.IsAdmin && (
                <Button
                  variant="outline"
                  size="sm"
                  className="mt-4"
                  onClick={() => navigate('/modules')}
                >
                  <Plus className="h-4 w-4 mr-2" />
                  Create your first module
                </Button>
              )}
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
              {modules.slice(0, 6).map((module) => (
                <div
                  key={module.ID}
                  className="flex flex-col p-4 border rounded-lg hover:bg-gray-50 cursor-pointer transition-colors hover:shadow-sm"
                  onClick={() => navigate(`/modules/${module.ID}`)}
                >
                  <div className="flex items-center gap-2 mb-2">
                    <BookOpen className="h-4 w-4 text-primary flex-shrink-0" />
                    <span className="font-semibold text-sm">{module.Code}</span>
                  </div>
                  <p className="text-sm text-muted-foreground line-clamp-1 mb-2">
                    {module.Name}
                  </p>
                  <Badge variant="outline" className="w-fit text-xs mt-auto">
                    {module.DepartmentName}
                  </Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Admin Quick Actions */}
      {user?.IsAdmin && (
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Quick Actions</CardTitle>
            <Plus className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate('/modules')}
            >
              <Plus className="h-4 w-4 mr-2" />
              New Module
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate('/academic-terms')}
            >
              <Plus className="h-4 w-4 mr-2" />
              New Term
            </Button>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

export default HomePage
