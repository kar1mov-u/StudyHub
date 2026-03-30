import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Select } from '@/components/ui/select'
import { ArrowLeft, Loader2, Calendar, Upload, Link as LinkIcon, BookOpen, X, CheckSquare, Plus, BookMarked, GraduationCap } from 'lucide-react'
import ResourceList from '@/components/resources/ResourceList'
import ResourceUploadDialog from '@/components/resources/ResourceUploadDialog'
import ResourceLinkDialog from '@/components/resources/ResourceLinkDialog'
import CommentSection from '@/components/comments/CommentSection'
import DeckStats from '@/components/decks/DeckStats'
import DeckList from '@/components/decks/DeckList'
import CreateCustomCardDialog from '@/components/decks/CreateCustomCardDialog'
import EditCardDialog from '@/components/decks/EditCardDialog'
import { resourcesApi } from '@/api/resources'
import { modulesApi } from '@/api/modules'
import { decksApi } from '@/api/decks'
import { useAuth } from '@/context/AuthContext'
import { useToast } from '@/components/ui/toast'
import type { Resource, ModulePage, UserDeckCard, DeckStats as DeckStatsType } from '@/types'

const WeekDetailPage: React.FC = () => {
  const { moduleId, weekId } = useParams<{ moduleId: string; weekId: string }>()
  const navigate = useNavigate()
  const { user } = useAuth()
  const { showToast } = useToast()
  
  const [resources, setResources] = useState<Resource[]>([])
  const [modulePage, setModulePage] = useState<ModulePage | null>(null)
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<'files' | 'links' | 'deck'>('files')
  const [uploadFileDialogOpen, setUploadFileDialogOpen] = useState(false)
  const [uploadLinkDialogOpen, setUploadLinkDialogOpen] = useState(false)

  // Deck state
  const [deckCards, setDeckCards] = useState<UserDeckCard[]>([])
  const [deckStats, setDeckStats] = useState<DeckStatsType | null>(null)
  const [deckLoading, setDeckLoading] = useState(false)
  const [createCardDialogOpen, setCreateCardDialogOpen] = useState(false)
  const [editCardDialogOpen, setEditCardDialogOpen] = useState(false)
  const [selectedCard, setSelectedCard] = useState<UserDeckCard | null>(null)
  const [deckStudyFilter, setDeckStudyFilter] = useState<'all' | 'reviewed' | 'not-reviewed' | 'difficult'>('all')

  // Study selection mode state
  const [selectMode, setSelectMode] = useState(false)
  const [selectedResourceIds, setSelectedResourceIds] = useState<Set<string>>(new Set())

  const weekNumber = modulePage?.Weeks?.find(w => w.ID === weekId)?.Number

  // Filter resources by type
  const fileResources = resources.filter(r => r.ResourceType === 'file')
  const linkResources = resources.filter(r => r.ResourceType === 'link')
  
  // Get counts for tab badges
  const fileCount = fileResources.length
  const linkCount = linkResources.length
  const deckCount = deckCards.length

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

  const loadDeckData = async () => {
    if (!weekId || !user?.ID) return

    try {
      setDeckLoading(true)
      const [cards, stats] = await Promise.all([
        decksApi.getUserDeck(weekId),
        decksApi.getDeckStats(weekId),
      ])
      setDeckCards(cards)
      setDeckStats(stats)
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to load deck data',
        'error'
      )
    } finally {
      setDeckLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [moduleId, weekId, showToast])

  useEffect(() => {
    if (activeTab === 'deck') {
      loadDeckData()
    }
  }, [activeTab, weekId, user?.ID])

  const handleDeleteResource = (resourceId: string) => {
    setResources(prevResources => prevResources.filter(r => r.ID !== resourceId))
    showToast('Resource deleted successfully', 'success')
  }

  const handleEnterSelectMode = () => {
    setSelectMode(true)
    setSelectedResourceIds(new Set())
    // Switch to files tab since only files have ObjectIDs
    setActiveTab('files')
  }

  const handleCancelSelectMode = () => {
    setSelectMode(false)
    setSelectedResourceIds(new Set())
  }

  const handleToggleSelect = (resourceId: string) => {
    setSelectedResourceIds(prev => {
      const next = new Set(prev)
      if (next.has(resourceId)) {
        next.delete(resourceId)
      } else {
        next.add(resourceId)
      }
      return next
    })
  }

  const handleSelectAll = () => {
    if (selectedResourceIds.size === fileResources.length) {
      // Deselect all
      setSelectedResourceIds(new Set())
    } else {
      // Select all file resources
      setSelectedResourceIds(new Set(fileResources.map(r => r.ID)))
    }
  }

  const handleStartStudying = () => {
    // Get the ObjectIDs for the selected resources
    const selectedObjectIds = fileResources
      .filter(r => selectedResourceIds.has(r.ID) && r.ObjectID)
      .map(r => r.ObjectID)

    if (selectedObjectIds.length === 0) {
      showToast('Please select files that have content available', 'error')
      return
    }

    // Navigate to study page with the object IDs and weekId
    navigate('/study', {
      state: {
        objectIds: selectedObjectIds,
        weekId: weekId,
        weekNumber,
        moduleCode: modulePage?.Module.Code,
        moduleName: modulePage?.Module.Name,
      },
    })
  }

  const handleCreateCard = async (front: string, back: string) => {
    if (!weekId) return

    try {
      await decksApi.createCustomCard(weekId, { front, back })
      showToast('Custom card created successfully', 'success')
      setCreateCardDialogOpen(false)
      loadDeckData()
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to create card',
        'error'
      )
    }
  }

  const handleEditCard = (card: UserDeckCard) => {
    setSelectedCard(card)
    setEditCardDialogOpen(true)
  }

  const handleUpdateCard = async (cardId: string, front?: string, back?: string) => {
    try {
      await decksApi.updateCard(cardId, { front, back })
      showToast('Card updated successfully', 'success')
      setEditCardDialogOpen(false)
      setSelectedCard(null)
      loadDeckData()
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to update card',
        'error'
      )
    }
  }

  const handleDeleteCard = async (cardId: string) => {
    try {
      await decksApi.removeCard(cardId)
      showToast('Card removed from deck', 'success')
      loadDeckData()
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to delete card',
        'error'
      )
    }
  }

  const handleStudyDeck = () => {
    if (deckCards.length === 0) {
      showToast('No cards in your deck yet', 'error')
      return
    }

    // Filter cards based on selected filter
    let cardsToStudy = [...deckCards]
    
    switch (deckStudyFilter) {
      case 'reviewed':
        cardsToStudy = deckCards.filter(card => card.ReviewCount > 0)
        break
      case 'not-reviewed':
        cardsToStudy = deckCards.filter(card => card.ReviewCount === 0)
        break
      case 'difficult':
        // Cards rated 4 or 5 difficulty, or not rated yet
        cardsToStudy = deckCards.filter(card => 
          !card.DifficultyRating || card.DifficultyRating >= 4
        )
        break
      case 'all':
      default:
        // All cards
        break
    }

    if (cardsToStudy.length === 0) {
      showToast(`No cards match the filter: ${deckStudyFilter}`, 'error')
      return
    }

    // Navigate to study page with deck mode
    navigate('/study', {
      state: {
        deckMode: true,
        deckCards: cardsToStudy,
        weekId: weekId,
        weekNumber,
        moduleCode: modulePage?.Module.Code,
        moduleName: modulePage?.Module.Name,
      },
    })
  }

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
              {!selectMode && fileResources.length > 0 && (
                <Button onClick={handleEnterSelectMode} variant="outline">
                  <BookOpen className="h-4 w-4 mr-2" />
                  Study
                </Button>
              )}
            </div>
          </CardHeader>
        </Card>
      </div>

      <Tabs value={activeTab} onValueChange={(value) => { if (!selectMode) setActiveTab(value as 'files' | 'links' | 'deck') }}>
        <div className="flex items-center justify-between mb-4">
          <TabsList>
            <TabsTrigger value="files">
              Files ({fileCount})
            </TabsTrigger>
            <TabsTrigger value="links" className={selectMode ? 'opacity-50 cursor-not-allowed' : ''}>
              Links ({linkCount})
            </TabsTrigger>
            <TabsTrigger value="deck" className={selectMode ? 'opacity-50 cursor-not-allowed' : ''}>
              <BookMarked className="h-4 w-4 mr-2" />
              My Deck ({deckCount})
            </TabsTrigger>
          </TabsList>

          {!selectMode && !user?.IsAdmin && (
            <>
              {activeTab === 'files' ? (
                <Button onClick={() => setUploadFileDialogOpen(true)}>
                  <Upload className="h-4 w-4 mr-2" />
                  Upload File
                </Button>
              ) : activeTab === 'links' ? (
                <Button onClick={() => setUploadLinkDialogOpen(true)}>
                  <LinkIcon className="h-4 w-4 mr-2" />
                  Add Link
                </Button>
              ) : activeTab === 'deck' ? (
                <div className="flex items-center gap-2">
                  <Select 
                    value={deckStudyFilter} 
                    onChange={(e) => setDeckStudyFilter(e.target.value as any)}
                    className="w-[160px]"
                  >
                    <option value="all">All cards</option>
                    <option value="not-reviewed">Not reviewed</option>
                    <option value="reviewed">Reviewed</option>
                    <option value="difficult">Difficult (4-5★)</option>
                  </Select>
                  <Button variant="outline" onClick={handleStudyDeck} disabled={deckCount === 0}>
                    <GraduationCap className="h-4 w-4 mr-2" />
                    Study Deck
                  </Button>
                  <Button onClick={() => setCreateCardDialogOpen(true)}>
                    <Plus className="h-4 w-4 mr-2" />
                    Create Custom Card
                  </Button>
                </div>
              ) : null}
            </>
          )}

          {selectMode && (
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" onClick={handleSelectAll}>
                <CheckSquare className="h-4 w-4 mr-2" />
                {selectedResourceIds.size === fileResources.length ? 'Deselect All' : 'Select All'}
              </Button>
            </div>
          )}
        </div>

        <TabsContent value="files">
          <ResourceList 
            resources={fileResources} 
            isLoading={loading}
            emptyMessage="No files uploaded yet"
            currentUserId={user?.ID}
            onDelete={selectMode ? undefined : handleDeleteResource}
            selectable={selectMode}
            selectedIds={selectedResourceIds}
            onToggleSelect={handleToggleSelect}
          />
        </TabsContent>

        <TabsContent value="links">
          <ResourceList 
            resources={linkResources} 
            isLoading={loading}
            emptyMessage="No links added yet"
            currentUserId={user?.ID}
            onDelete={handleDeleteResource}
          />
        </TabsContent>

        <TabsContent value="deck">
          <div className="space-y-6">
            {deckStats && <DeckStats stats={deckStats} />}
            
            <DeckList
              cards={deckCards}
              isLoading={deckLoading}
              onEdit={handleEditCard}
              onDelete={handleDeleteCard}
            />
          </div>
        </TabsContent>
      </Tabs>

      {/* Comments / Discussion Section */}
      {weekId && <CommentSection weekId={weekId} />}

      {/* Selection mode bottom action bar */}
      {selectMode && (
        <div className="fixed bottom-0 left-0 right-0 bg-background border-t shadow-lg p-4 z-50">
          <div className="max-w-4xl mx-auto flex items-center justify-between">
            <div className="text-sm font-medium">
              {selectedResourceIds.size} file{selectedResourceIds.size !== 1 ? 's' : ''} selected
            </div>
            <div className="flex items-center gap-3">
              <Button variant="outline" onClick={handleCancelSelectMode}>
                <X className="h-4 w-4 mr-2" />
                Cancel
              </Button>
              <Button
                onClick={handleStartStudying}
                disabled={selectedResourceIds.size === 0}
              >
                <BookOpen className="h-4 w-4 mr-2" />
                Start Studying
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Add bottom padding when selection bar is visible */}
      {selectMode && <div className="h-20" />}

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
          <CreateCustomCardDialog
            open={createCardDialogOpen}
            onOpenChange={setCreateCardDialogOpen}
            onSubmit={handleCreateCard}
          />
          {selectedCard && (
            <EditCardDialog
              open={editCardDialogOpen}
              onOpenChange={(open) => {
                setEditCardDialogOpen(open)
                if (!open) setSelectedCard(null)
              }}
              card={selectedCard}
              onSubmit={handleUpdateCard}
            />
          )}
        </>
      )}
    </div>
  )
}

export default WeekDetailPage
