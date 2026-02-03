import React, { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { Card, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Loader2, User as UserIcon, Mail } from 'lucide-react'
import UserResourceList from '@/components/resources/UserResourceList'
import { resourcesApi } from '@/api/resources'
import { authApi } from '@/api/auth'
import { useAuth } from '@/context/AuthContext'
import { useToast } from '@/components/ui/toast'
import type { UserResource, User } from '@/types'

const UserProfilePage: React.FC = () => {
  const { userId } = useParams<{ userId: string }>()
  const { user: currentUser } = useAuth()
  const { showToast } = useToast()
  
  const [resources, setResources] = useState<UserResource[]>([])
  const [profileUser, setProfileUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<'files' | 'links'>('files')

  // Determine if viewing own profile - check both userId param and if it matches current user
  const isOwnProfile = currentUser?.ID === userId

  // Filter resources by type
  const fileResources = resources.filter(r => r.ResourceType === 'file')
  const linkResources = resources.filter(r => r.ResourceType === 'link')
  
  // Get counts for tab badges
  const fileCount = fileResources.length
  const linkCount = linkResources.length

  const loadData = async () => {
    // Ensure we have a valid user ID before making any requests
    if (!userId) {
      setLoading(false)
      return
    }

    try {
      setLoading(true)
      const [resourcesData, userData] = await Promise.all([
        resourcesApi.getUserResources(userId),
        authApi.getCurrentUser(userId),
      ])
      setResources(resourcesData || [])
      setProfileUser(userData)
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to load profile',
        'error'
      )
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    // Only load data if userId exists
    if (userId) {
      loadData()
    }
  }, [userId])

  const handleDeleteResource = (resourceId: string) => {
    setResources(prevResources => prevResources.filter(r => r.ID !== resourceId))
    showToast('Resource deleted successfully', 'success')
  }

  // Show loading while waiting for userId or data
  if (!userId || loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!profileUser) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground">User not found</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <div className="flex items-start gap-4">
            <div className="h-16 w-16 rounded-full bg-primary/10 flex items-center justify-center">
              <UserIcon className="h-8 w-8 text-primary" />
            </div>
            <div className="flex-1">
              <CardTitle className="text-2xl">
                {isOwnProfile ? 'My Profile' : `${profileUser.FirstName} ${profileUser.LastName}`}
              </CardTitle>
              <CardDescription className="flex items-center gap-2 mt-2">
                <Mail className="h-4 w-4" />
                {profileUser.Email}
              </CardDescription>
              <div className="mt-3 text-sm text-muted-foreground">
                {resources.length} resource{resources.length !== 1 ? 's' : ''} uploaded
              </div>
            </div>
          </div>
        </CardHeader>
      </Card>

      <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as 'files' | 'links')}>
        <TabsList>
          <TabsTrigger value="files">
            Files ({fileCount})
          </TabsTrigger>
          <TabsTrigger value="links">
            Links ({linkCount})
          </TabsTrigger>
        </TabsList>

        <TabsContent value="files" className="mt-6">
          <UserResourceList 
            resources={fileResources} 
            isLoading={loading}
            emptyMessage="No files uploaded yet"
            showDelete={isOwnProfile}
            onDelete={handleDeleteResource}
          />
        </TabsContent>

        <TabsContent value="links" className="mt-6">
          <UserResourceList 
            resources={linkResources} 
            isLoading={loading}
            emptyMessage="No links added yet"
            showDelete={isOwnProfile}
            onDelete={handleDeleteResource}
          />
        </TabsContent>
      </Tabs>
    </div>
  )
}

export default UserProfilePage
