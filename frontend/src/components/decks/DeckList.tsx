import React from 'react'
import { Loader2, BookMarked } from 'lucide-react'
import DeckCard from './DeckCard'
import type { UserDeckCard } from '@/types'

interface DeckListProps {
  cards: UserDeckCard[]
  isLoading?: boolean
  currentUserId?: string
  onEdit?: (card: UserDeckCard) => void
  onDelete?: (cardId: string) => void
}

const DeckList: React.FC<DeckListProps> = ({ 
  cards, 
  isLoading = false,
  currentUserId,
  onEdit,
  onDelete
}) => {
  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (cards.length === 0) {
    return (
      <div className="text-center py-12">
        <BookMarked className="h-16 w-16 mx-auto text-muted-foreground mb-4" />
        <h3 className="text-lg font-semibold mb-2">No cards in your deck yet</h3>
        <p className="text-muted-foreground max-w-md mx-auto">
          Start by adding flashcards while studying or create your own custom cards!
        </p>
      </div>
    )
  }

  // Sort: custom cards first, then by creation date (newest first)
  const sortedCards = [...cards].sort((a, b) => {
    if (a.IsCustom && !b.IsCustom) return -1
    if (!a.IsCustom && b.IsCustom) return 1
    return new Date(b.CreatedAt).getTime() - new Date(a.CreatedAt).getTime()
  })

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {sortedCards.map((card) => (
        <DeckCard
          key={card.ID}
          card={card}
          currentUserId={currentUserId}
          onEdit={onEdit}
          onDelete={onDelete}
        />
      ))}
    </div>
  )
}

export default DeckList
