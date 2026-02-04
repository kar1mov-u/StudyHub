import React, { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Plus, BookOpen, PlayCircle, Users } from 'lucide-react'
import StatsCard from '@/components/admin/StatsCard'
import ModulesTable from '@/components/admin/ModulesTable'
import ModuleRunsTable from '@/components/admin/ModuleRunsTable'
import ModuleForm from '@/components/modules/ModuleForm'
import { modulesApi } from '@/api/modules'
import { moduleRunsApi } from '@/api/moduleRuns'
import { useToast } from '@/components/ui/toast'
import type { Module, ModuleRun } from '@/types'

const AdminDashboardPage: React.FC = () => {
  const { showToast } = useToast()
  const [activeTab, setActiveTab] = useState<'modules' | 'runs'>('modules')
  const [modules, setModules] = useState<Module[]>([])
  const [allRuns, setAllRuns] = useState<ModuleRun[]>([])
  const [loading, setLoading] = useState(true)
  const [createModuleOpen, setCreateModuleOpen] = useState(false)

  // Load initial data
  const loadData = async () => {
    try {
      setLoading(true)
      const modulesData = await modulesApi.listModules()
      setModules(modulesData)

      // Load all runs for all modules
      const allRunsPromises = modulesData.map((module) =>
        moduleRunsApi.listModuleRuns(module.ID).catch(() => [])
      )
      const runsArrays = await Promise.all(allRunsPromises)
      const flatRuns = runsArrays.flat()
      setAllRuns(flatRuns)
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
  }, [])

  const handleDeleteModule = async (id: string) => {
    try {
      await modulesApi.deleteModule(id)
      showToast('Module deleted successfully', 'success')
      loadData()
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to delete module',
        'error'
      )
    }
  }

  const handleDeleteRun = async (id: string) => {
    try {
      await moduleRunsApi.deleteModuleRun(id)
      showToast('Module run deleted successfully', 'success')
      loadData()
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to delete module run',
        'error'
      )
    }
  }

  // Calculate stats
  const totalModules = modules.length
  const totalRuns = allRuns.length

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Admin Dashboard</h1>
        <p className="text-muted-foreground mt-2">
          Manage modules, runs, and system settings
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <StatsCard
          title="Total Modules"
          value={totalModules}
          icon={BookOpen}
          description="All modules in the system"
        />
        <StatsCard
          title="Total Runs"
          value={totalRuns}
          icon={PlayCircle}
          description="All module runs"
        />
        <StatsCard
          title="Total Users"
          value="—"
          icon={Users}
          description="Registered users"
        />
      </div>

      {/* Management Tabs */}
      <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as 'modules' | 'runs')}>
        <div className="flex items-center justify-between mb-4">
          <TabsList>
            <TabsTrigger value="modules">Modules ({totalModules})</TabsTrigger>
            <TabsTrigger value="runs">Module Runs ({totalRuns})</TabsTrigger>
          </TabsList>

          {activeTab === 'modules' && (
            <Button onClick={() => setCreateModuleOpen(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Create Module
            </Button>
          )}
        </div>

        <TabsContent value="modules" className="mt-0">
          <ModulesTable
            modules={modules}
            isLoading={loading}
            onDelete={handleDeleteModule}
            onRefresh={loadData}
          />
        </TabsContent>

        <TabsContent value="runs" className="mt-0">
          <div className="mb-4 p-4 bg-muted rounded-lg">
            <p className="text-sm text-muted-foreground">
              Module runs are created automatically when a new academic term is started. 
              Go to Academic Terms page to create a new term.
            </p>
          </div>
          <ModuleRunsTable
            runs={allRuns}
            modules={modules}
            isLoading={loading}
            onDelete={handleDeleteRun}
            onRefresh={loadData}
          />
        </TabsContent>
      </Tabs>

      {/* Create Module Dialog */}
      <ModuleForm
        open={createModuleOpen}
        onOpenChange={setCreateModuleOpen}
        onSuccess={loadData}
      />
    </div>
  )
}

export default AdminDashboardPage
