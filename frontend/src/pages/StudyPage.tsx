import React, { useState, useEffect, useCallback } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ArrowLeft, ArrowRight, Loader2, RotateCcw, BookOpen, ChevronLeft, Plus, Check, Star } from 'lucide-react'
import { contentsApi } from '@/api/contents'
import { decksApi } from '@/api/decks'
import { useToast } from '@/components/ui/toast'
import type { Flashcard, UserDeckCard } from '@/types'

interface StudyPageState {
  objectIds?: string[]
  weekId?: string
  weekNumber?: number
  moduleCode?: string
  moduleName?: string
  // Deck mode
  deckMode?: boolean
  deckCards?: UserDeckCard[]
}

const StudyPage: React.FC = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const { showToast } = useToast()

  const state = location.state as StudyPageState | null

  const [flashcards, setFlashcards] = useState<Flashcard[]>([])
  const [loading, setLoading] = useState(true)
  const [currentIndex, setCurrentIndex] = useState(0)
  const [flipped, setFlipped] = useState(false)
  
  // Deck management state
  const [deckCardIds, setDeckCardIds] = useState<Set<string>>(new Set())
  const [addingToDeck, setAddingToDeck] = useState(false)
  const [deckCardMap, setDeckCardMap] = useState<Map<string, string>>(new Map()) // flashcardId -> deckCardId
  
  // Deck mode: studying user's personal deck cards
  const isDeckMode = state?.deckMode === true
  const [studyingDeckCards, setStudyingDeckCards] = useState<UserDeckCard[]>([])

  const loadFlashcards = useCallback(async () => {
    if (isDeckMode) {
      // In deck mode, use the deck cards passed from state
      if (state?.deckCards) {
        setStudyingDeckCards(state.deckCards)
        setLoading(false)
      }
      return
    }

    // Normal mode: load flashcards from objects
    if (!state?.objectIds || state.objectIds.length === 0) return

    try {
      setLoading(true)
      const cards = await contentsApi.getFlashcards(state.objectIds)
      setFlashcards(cards || [])
    } catch (error) {
      showToast(
        error instanceof Error ? error.message : 'Failed to load flashcards',
        'error'
      )
    } finally {
      setLoading(false)
    }
  }, [state?.objectIds, state?.deckCards, isDeckMode, showToast])

  // Load user's deck to check which cards are already added
  const loadDeck = useCallback(async () => {
    if (!state?.weekId || isDeckMode) return // Don't load in deck mode

    try {
      const deck = await decksApi.getUserDeck(state.weekId)
      const cardIds = new Set<string>()
      const cardMap = new Map<string, string>()
      
      deck.forEach((card: UserDeckCard) => {
        if (card.SourceFlashcardID) {
          cardIds.add(card.SourceFlashcardID)
          cardMap.set(card.SourceFlashcardID, card.ID)
        }
      })
      
      setDeckCardIds(cardIds)
      setDeckCardMap(cardMap)
    } catch (error) {
      // Silent fail - deck loading is optional
      console.error('Failed to load deck:', error)
    }
  }, [state?.weekId, isDeckMode])

  useEffect(() => {
    if (!state?.objectIds && !state?.deckMode) {
      navigate('/modules', { replace: true })
      return
    }
    loadFlashcards()
    loadDeck()
  }, [state, navigate, loadFlashcards, loadDeck])

  // Use deck cards if in deck mode, otherwise use flashcards
  const cards = isDeckMode ? studyingDeckCards : flashcards
  const currentCard = cards[currentIndex]
  const totalCards = cards.length
  
  // In deck mode, we're already in the deck so no need for "add to deck" logic
  const isInDeck = isDeckMode ? true : (currentCard && 'ID' in currentCard ? deckCardIds.has(currentCard.ID) : false)

  const handleFlip = () => {
    setFlipped(prev => !prev)
  }

  const handleNext = () => {
    if (currentIndex < totalCards - 1) {
      setCurrentIndex(prev => prev + 1)
      setFlipped(false)
    }
  }

  const handlePrevious = () => {
    if (currentIndex > 0) {
      setCurrentIndex(prev => prev - 1)
      setFlipped(false)
    }
  }

  const handleRestart = () => {
    setCurrentIndex(0)
    setFlipped(false)
  }

  const handleAddToDeck = async () => {
    if (!currentCard || !state?.weekId || isDeckMode) return

    try {
      setAddingToDeck(true)
      await decksApi.addCardToDeck(state.weekId, { flashcard_id: (currentCard as Flashcard).ID })
      setDeckCardIds(prev => new Set(prev).add((currentCard as Flashcard).ID))
      showToast('Added to your deck!', 'success')
      
      // Reload deck to get the new card ID
      await loadDeck()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to add card'
      if (message.includes('already')) {
        showToast('Card is already in your deck', 'error')
      } else {
        showToast(message, 'error')
      }
    } finally {
      setAddingToDeck(false)
    }
  }

  const handleRateDifficulty = async (rating: number) => {
    if (!currentCard) return

    // Get the deck card ID
    const deckCardId = isDeckMode 
      ? (currentCard as UserDeckCard).ID 
      : deckCardMap.get((currentCard as Flashcard).ID)
    
    if (!deckCardId) return

    try {
      await decksApi.recordReview(deckCardId, { difficulty_rating: rating })
      showToast(`Rated: ${rating}/5`, 'success')
      
      // Update local state if in deck mode
      if (isDeckMode) {
        setStudyingDeckCards(prev => 
          prev.map(card => 
            card.ID === deckCardId 
              ? { ...card, DifficultyRating: rating, ReviewCount: card.ReviewCount + 1 }
              : card
          )
        )
      }
    } catch (error) {
      showToast('Failed to record rating', 'error')
    }
  }

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      switch (e.key) {
        case ' ':
        case 'Enter':
          e.preventDefault()
          handleFlip()
          break
        case 'ArrowRight':
          e.preventDefault()
          handleNext()
          break
        case 'ArrowLeft':
          e.preventDefault()
          handlePrevious()
          break
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [currentIndex, totalCards])

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center py-24">
        <Loader2 className="h-10 w-10 animate-spin text-primary mb-4" />
        <p className="text-muted-foreground">Loading flashcards...</p>
      </div>
    )
  }

  if (!state?.objectIds && !state?.deckMode) {
    return null
  }

  if (cards.length === 0) {
    return (
      <div className="space-y-6">
        <Button variant="ghost" onClick={() => navigate(-1)}>
          <ChevronLeft className="h-4 w-4 mr-2" />
          Back
        </Button>
        <div className="text-center py-24">
          <BookOpen className="h-16 w-16 mx-auto text-muted-foreground mb-4" />
          <h2 className="text-xl font-semibold mb-2">
            {isDeckMode ? 'No Cards in Deck' : 'No Flashcards Available'}
          </h2>
          <p className="text-muted-foreground max-w-md mx-auto">
            {isDeckMode 
              ? 'Your deck is empty or no cards match the selected filter. Add cards to your deck from the study page.'
              : 'No flashcards have been generated for the selected files yet. Flashcards are generated automatically -- please check back later.'}
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6 max-w-3xl mx-auto">
      {/* Header */}
      <div className="flex items-center justify-between">
        <Button variant="ghost" onClick={() => navigate(-1)}>
          <ChevronLeft className="h-4 w-4 mr-2" />
          Back
        </Button>
        <div className="text-center">
          {state.moduleCode && (
            <p className="text-sm text-muted-foreground">
              {state.moduleCode}{state.weekNumber ? ` - Week ${state.weekNumber}` : ''}
            </p>
          )}
        </div>
        <div className="w-20" /> {/* Spacer for centering */}
      </div>

      {/* Progress */}
      <div className="space-y-2">
        <div className="flex items-center justify-between text-sm text-muted-foreground">
          <span>Card {currentIndex + 1} of {totalCards}</span>
          <Badge variant="outline">
            {Math.round(((currentIndex + 1) / totalCards) * 100)}%
          </Badge>
        </div>
        <div className="w-full bg-muted rounded-full h-2">
          <div
            className="bg-primary h-2 rounded-full transition-all duration-300"
            style={{ width: `${((currentIndex + 1) / totalCards) * 100}%` }}
          />
        </div>
      </div>

      {/* Flashcard */}
      <div
        className="perspective-1000 cursor-pointer"
        style={{ perspective: '1000px' }}
        onClick={handleFlip}
      >
        <div
          className="relative w-full transition-transform duration-500"
          style={{
            transformStyle: 'preserve-3d',
            transform: flipped ? 'rotateY(180deg)' : 'rotateY(0deg)',
          }}
        >
          {/* Front face */}
          <Card
            className="w-full min-h-[300px] flex flex-col items-center justify-center p-8"
            style={{
              backfaceVisibility: 'hidden',
            }}
          >
            <Badge variant="secondary" className="mb-4">Front</Badge>
            <p className="text-xl text-center leading-relaxed whitespace-pre-wrap">
              {currentCard.Front}
            </p>
            <p className="text-sm text-muted-foreground mt-6">
              Click to flip
            </p>
            {isDeckMode && (
              <div className="mt-4 text-xs text-muted-foreground">
                <p>Reviews: {(currentCard as UserDeckCard).ReviewCount || 0}</p>
                {(currentCard as UserDeckCard).DifficultyRating && (
                  <p>Last rating: {(currentCard as UserDeckCard).DifficultyRating}/5</p>
                )}
              </div>
            )}
          </Card>

          {/* Back face */}
          <Card
            className="w-full min-h-[300px] flex flex-col items-center justify-center p-8 absolute top-0 left-0 bg-primary/5 space-y-4"
            style={{
              backfaceVisibility: 'hidden',
              transform: 'rotateY(180deg)',
            }}
          >
            <Badge variant="default" className="mb-2">Back</Badge>
            <p className="text-xl text-center leading-relaxed whitespace-pre-wrap">
              {currentCard.Back}
            </p>
            
            {/* Add to Deck Button or Difficulty Rating - only show if weekId is available */}
            {state?.weekId && (
              <div className="mt-4 pt-4 border-t w-full flex flex-col items-center gap-3">
                {!isDeckMode && !isInDeck ? (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={(e) => {
                      e.stopPropagation()
                      handleAddToDeck()
                    }}
                    disabled={addingToDeck}
                  >
                    {addingToDeck ? (
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    ) : (
                      <Plus className="h-4 w-4 mr-2" />
                    )}
                    Add to My Deck
                  </Button>
                ) : (
                  <>
                    {!isDeckMode && (
                      <Badge variant="secondary" className="text-xs">
                        <Check className="h-3 w-3 mr-1" />
                        In Your Deck
                      </Badge>
                    )}
                    
                    {/* Difficulty Rating - available in both modes if card is in deck */}
                    <div className="flex flex-col items-center gap-2">
                      <p className="text-xs text-muted-foreground">Rate difficulty:</p>
                      <div className="flex gap-1">
                        {[1, 2, 3, 4, 5].map((rating) => (
                          <button
                            key={rating}
                            onClick={(e) => {
                              e.stopPropagation()
                              handleRateDifficulty(rating)
                            }}
                            className="p-1 hover:scale-110 transition-transform"
                            title={`Rate ${rating}/5`}
                          >
                            <Star className="h-5 w-5 text-yellow-500 hover:fill-yellow-500" />
                          </button>
                        ))}
                      </div>
                    </div>
                  </>
                )}
              </div>
            )}

            <p className="text-sm text-muted-foreground mt-2">
              Click to flip back
            </p>
          </Card>
        </div>
      </div>

      {/* Navigation controls */}
      <div className="flex items-center justify-between">
        <Button
          variant="outline"
          onClick={handlePrevious}
          disabled={currentIndex === 0}
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Previous
        </Button>

        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={handleFlip}>
            <RotateCcw className="h-4 w-4 mr-2" />
            Flip
          </Button>
          {currentIndex === totalCards - 1 && (
            <Button variant="outline" size="sm" onClick={handleRestart}>
              Restart
            </Button>
          )}
        </div>

        <Button
          variant="outline"
          onClick={handleNext}
          disabled={currentIndex === totalCards - 1}
        >
          Next
          <ArrowRight className="h-4 w-4 ml-2" />
        </Button>
      </div>

      {/* Keyboard hints */}
      <div className="text-center text-xs text-muted-foreground">
        <span className="inline-flex items-center gap-4">
          <span><kbd className="px-1.5 py-0.5 bg-muted rounded text-xs">Space</kbd> flip</span>
          <span><kbd className="px-1.5 py-0.5 bg-muted rounded text-xs">&larr;</kbd> previous</span>
          <span><kbd className="px-1.5 py-0.5 bg-muted rounded text-xs">&rarr;</kbd> next</span>
        </span>
      </div>
    </div>
  )
}

export default StudyPage
