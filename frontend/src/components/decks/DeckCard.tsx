import React from 'react'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Edit, Trash2, Star, RotateCcw } from 'lucide-react'
import type { UserDeckCard } from '@/types'

interface DeckCardProps {
  card: UserDeckCard
  currentUserId?: string
  onEdit?: (card: UserDeckCard) => void
  onDelete?: (cardId: string) => void
}

const DeckCard: React.FC<DeckCardProps> = ({ 
  card, 
  currentUserId,
  onEdit, 
  onDelete 
}) => {
  const [isFlipped, setIsFlipped] = React.useState(false)
  const isOwner = currentUserId === card.UserID

  const renderStars = (rating: number | null) => {
    if (!rating) return <span className="text-xs text-muted-foreground">Not rated</span>
    
    return (
      <div className="flex items-center gap-0.5">
        {[1, 2, 3, 4, 5].map((star) => (
          <Star
            key={star}
            className={`h-3 w-3 ${
              star <= rating
                ? 'fill-yellow-500 text-yellow-500'
                : 'text-muted-foreground'
            }`}
          />
        ))}
      </div>
    )
  }

  return (
    <Card className="group hover:shadow-md transition-shadow">
      <CardContent className="pt-6">
        <div className="flex items-start justify-between mb-3">
          <Badge variant={card.IsCustom ? 'secondary' : 'outline'} className="text-xs">
            {card.IsCustom ? 'Custom' : 'From File'}
          </Badge>
          {card.ReviewCount > 0 && (
            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <RotateCcw className="h-3 w-3" />
              <span>{card.ReviewCount}</span>
            </div>
          )}
        </div>

        <div 
          className="cursor-pointer min-h-[120px] flex items-center justify-center"
          onClick={() => setIsFlipped(!isFlipped)}
        >
          {!isFlipped ? (
            <div className="text-center space-y-2">
              <Badge variant="outline" className="mb-2 text-xs">Front</Badge>
              <p className="text-sm font-medium line-clamp-3">
                {card.Front}
              </p>
              <p className="text-xs text-muted-foreground mt-2">
                Click to see answer
              </p>
            </div>
          ) : (
            <div className="text-center space-y-2">
              <Badge variant="default" className="mb-2 text-xs">Back</Badge>
              <p className="text-sm line-clamp-3">
                {card.Back}
              </p>
              <p className="text-xs text-muted-foreground mt-2">
                Click to see question
              </p>
            </div>
          )}
        </div>

        {/* Rating display */}
        <div className="mt-4 pt-3 border-t flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground">Difficulty:</span>
            {renderStars(card.DifficultyRating)}
          </div>
        </div>
      </CardContent>

      {isOwner && (onEdit || onDelete) && (
        <CardFooter className="pt-0 pb-4 gap-2">
          {onEdit && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => onEdit(card)}
              className="flex-1"
            >
              <Edit className="h-3 w-3 mr-1" />
              Edit
            </Button>
          )}
          {onDelete && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                if (confirm('Remove this card from your deck?')) {
                  onDelete(card.ID)
                }
              }}
              className="flex-1 text-destructive hover:text-destructive"
            >
              <Trash2 className="h-3 w-3 mr-1" />
              Delete
            </Button>
          )}
        </CardFooter>
      )}
    </Card>
  )
}

export default DeckCard
