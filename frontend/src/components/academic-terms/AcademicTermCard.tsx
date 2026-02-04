import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Calendar } from 'lucide-react'
import type { AcademicTerm } from '@/types'

interface AcademicTermCardProps {
  term: AcademicTerm
}

const AcademicTermCard: React.FC<AcademicTermCardProps> = ({ term }) => {
  return (
    <Card className="border-primary border-2">
      <CardHeader>
        <div className="flex items-center gap-3">
          <div className="h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center">
            <Calendar className="h-6 w-6 text-primary" />
          </div>
          <div>
            <CardTitle className="text-2xl capitalize">
              {term.Semester} {term.Year}
            </CardTitle>
            <p className="text-sm text-muted-foreground mt-1">
              Active academic term
            </p>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Semester:</span>
            <span className="font-medium capitalize">{term.Semester}</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Year:</span>
            <span className="font-medium">{term.Year}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

export default AcademicTermCard
