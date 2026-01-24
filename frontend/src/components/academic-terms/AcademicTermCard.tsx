import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Calendar, CheckCircle2, Circle } from 'lucide-react'
import type { AcademicTerm } from '@/types'

interface AcademicTermCardProps {
  term: AcademicTerm
  onActivate: (id: string) => void
  onDeactivate: (id: string) => void
}

const AcademicTermCard: React.FC<AcademicTermCardProps> = ({
  term,
  onActivate,
  onDeactivate,
}) => {
  return (
    <Card className={term.IsActive ? 'border-primary border-2' : ''}>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-2">
            <Calendar className="h-5 w-5 text-muted-foreground" />
            <CardTitle className="text-xl capitalize">
              {term.Semester} {term.Year}
            </CardTitle>
          </div>
          {term.IsActive ? (
            <Badge variant="success">
              <CheckCircle2 className="h-3 w-3 mr-1" />
              Active
            </Badge>
          ) : (
            <Badge variant="outline">
              <Circle className="h-3 w-3 mr-1" />
              Inactive
            </Badge>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex gap-2">
          {term.IsActive ? (
            <Button
              variant="outline"
              size="sm"
              onClick={() => onDeactivate(term.ID)}
            >
              Deactivate
            </Button>
          ) : (
            <Button
              variant="default"
              size="sm"
              onClick={() => onActivate(term.ID)}
            >
              Activate
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  )
}

export default AcademicTermCard
