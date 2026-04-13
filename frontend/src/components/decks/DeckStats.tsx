import React from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { BarChart3, BookMarked, Star, Calendar } from 'lucide-react'
import type { DeckStats as DeckStatsType } from '@/types'

interface DeckStatsProps {
  stats: DeckStatsType
}

const DeckStats: React.FC<DeckStatsProps> = ({ stats }) => {
  const progressPercentage = stats.total_cards > 0 
    ? Math.round((stats.reviewed_cards / stats.total_cards) * 100) 
    : 0

  const formatDate = (dateStr: string | null) => {
    if (!dateStr) return 'Never'
    const date = new Date(dateStr)
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  }

  return (
    <Card>
      <CardContent className="pt-6">
        <div className="flex items-center gap-2 mb-4">
          <BarChart3 className="h-5 w-5 text-primary" />
          <h3 className="font-semibold text-lg">Deck Statistics</h3>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {/* Total Cards */}
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <BookMarked className="h-4 w-4 text-muted-foreground" />
              <p className="text-xs text-muted-foreground">Total Cards</p>
            </div>
            <p className="text-2xl font-bold">{stats.total_cards}</p>
          </div>

          {/* Reviewed Cards */}
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <Badge variant="outline" className="text-xs">Progress</Badge>
            </div>
            <p className="text-2xl font-bold">
              {stats.reviewed_cards}
              <span className="text-sm text-muted-foreground ml-1">
                / {stats.total_cards}
              </span>
            </p>
            <div className="w-full bg-muted rounded-full h-1.5 mt-1">
              <div
                className="bg-primary h-1.5 rounded-full transition-all"
                style={{ width: `${progressPercentage}%` }}
              />
            </div>
          </div>

          {/* Average Rating */}
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <Star className="h-4 w-4 text-yellow-500 fill-yellow-500" />
              <p className="text-xs text-muted-foreground">Avg Rating</p>
            </div>
            <p className="text-2xl font-bold">
              {stats.average_rating > 0 ? stats.average_rating.toFixed(1) : '-'}
              <span className="text-sm text-muted-foreground ml-1">/ 5</span>
            </p>
          </div>

          {/* Last Reviewed */}
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <p className="text-xs text-muted-foreground">Last Reviewed</p>
            </div>
            <p className="text-lg font-semibold">
              {formatDate(stats.last_reviewed_at)}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

export default DeckStats
