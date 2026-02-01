import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { ArrowLeft, Loader2, Calendar, Upload, Link as LinkIcon } from 'lucide-react'
import ResourceList from '@/components/resources/ResourceList'
import ResourceUploadDialog from '@/components/resources/ResourceUploadDialog'
import ResourceLinkDialog from '@/components/resources/ResourceLinkDialog'
import { resourcesApi } from '@/api/resources'
import { modulesApi } from '@/api/modules'
import { useToast } from '@/components/ui/toast'
import type { Resource, ModulePage } from '@/types'

const WeekDetailPage: React.FC = () => {
  const { moduleId, weekId } = useParams<{ moduleId: string; weekId: string }>()
  const navigate = useNavigate()
  const { showToast } = useToast()
  
  const [resources, setResources] = useState<Resource[]>([])
  const [modulePage, setModulePage] = useState<ModulePage | null>(null)
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<'files' | 'links'>('files')
  const [uploadFileDialogOpen, setUploadFileDialogOpen] = useState(false)
  const [uploadLinkDialogOpen, setUploadLinkDialogOpen] = useState(false)

  const weekNumber = modulePage?.Weeks?.find(w => w.ID === weekId)?.Number

  // Filter resources by type
  const fileResources = resources.filter(r => r.ResourceType === 'file')
  const linkResources = resources.filter(r => r.ResourceType === 'link')
  
  // Get counts for tab badges
  const fileCount = fileResources.length
  const linkCount = linkResources.length

  const loadData = async () => {
    if (!moduleId || !weekId) return

    try {
      setLoading(true)
      const [moduleData, resourcesData] = await Promise.all([
        modulesApi.getModuleFull(moduleId),
        resourcesApi.getResourcesByWeek(weekId),
      ])
      setModulePage(moduleData)
      setResources(resourcesData || [])
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to load data',
        'error'
      )
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [moduleId, weekId, showToast])

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!modulePage) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground">Module not found</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <Button
          variant="ghost"
          onClick={() => navigate(`/modules/${moduleId}`)}
          className="mb-4"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Module
        </Button>
        
        <Card>
          <CardHeader>
            <div className="flex items-start justify-between">
              <div>
                <CardTitle className="text-2xl flex items-center gap-2">
                  <Calendar className="h-6 w-6" />
                  Week {weekNumber || '?'}
                </CardTitle>
                <p className="text-muted-foreground mt-2">
                  {modulePage.Module.Code} - {modulePage.Module.Name}
                </p>
                <p className="text-sm text-muted-foreground">
                  {modulePage.Run.Semester} {modulePage.Run.Year}
                </p>
              </div>
            </div>
          </CardHeader>
        </Card>
      </div>

      <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as 'files' | 'links')}>
        <div className="flex items-center justify-between mb-4">
          <TabsList>
            <TabsTrigger value="files">
              Files ({fileCount})
            </TabsTrigger>
            <TabsTrigger value="links">
              Links ({linkCount})
            </TabsTrigger>
          </TabsList>

          {activeTab === 'files' ? (
            <Button onClick={() => setUploadFileDialogOpen(true)}>
              <Upload className="h-4 w-4 mr-2" />
              Upload File
            </Button>
          ) : (
            <Button onClick={() => setUploadLinkDialogOpen(true)}>
              <LinkIcon className="h-4 w-4 mr-2" />
              Add Link
            </Button>
          )}
        </div>

        <TabsContent value="files">
          <ResourceList 
            resources={fileResources} 
            isLoading={loading}
            emptyMessage="No files uploaded yet"
          />
        </TabsContent>

        <TabsContent value="links">
          <ResourceList 
            resources={linkResources} 
            isLoading={loading}
            emptyMessage="No links added yet"
          />
        </TabsContent>
      </Tabs>

      {weekId && (
        <>
          <ResourceUploadDialog
            open={uploadFileDialogOpen}
            onOpenChange={setUploadFileDialogOpen}
            weekId={weekId}
            onSuccess={loadData}
          />
          <ResourceLinkDialog
            open={uploadLinkDialogOpen}
            onOpenChange={setUploadLinkDialogOpen}
            weekId={weekId}
            onSuccess={loadData}
          />
        </>
      )}
    </div>
  )
}

export default WeekDetailPage
