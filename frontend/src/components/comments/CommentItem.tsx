import React, { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  ThumbsUp,
  ThumbsDown,
  Reply,
  Send,
  User as UserIcon,
} from 'lucide-react'
import type { Comment, User } from '@/types'

function timeAgo(dateStr: string): string {
  const now = Date.now()
  const time = new Date(dateStr).getTime()
  const seconds = Math.floor((now - time) / 1000)

  if (seconds < 0 || isNaN(seconds)) return 'just now'
  if (seconds < 60) return 'just now'
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}d ago`
  const weeks = Math.floor(days / 7)
  if (weeks < 4) return `${weeks}w ago`
  const months = Math.floor(days / 30)
  if (months < 12) return `${months}mo ago`
  const years = Math.floor(days / 365)
  return `${years}y ago`
}

interface CommentItemProps {
  comment: Comment
  replies: Comment[]
  usersMap: Record<string, User>
  currentUserId?: string
  onUpvote: (commentId: string) => void
  onDownvote: (commentId: string) => void
  onReply: (parentId: string, content: string) => void
  isNested?: boolean
}

const CommentItem: React.FC<CommentItemProps> = ({
  comment,
  replies,
  usersMap,
  currentUserId,
  onUpvote,
  onDownvote,
  onReply,
  isNested = false,
}) => {
  const [showReplyInput, setShowReplyInput] = useState(false)
  const [replyContent, setReplyContent] = useState('')
  const [submittingReply, setSubmittingReply] = useState(false)

  const author = usersMap[comment.user_id]
  const authorName = author
    ? `${author.FirstName} ${author.LastName}`
    : 'Unknown User'
  const authorInitials = author
    ? `${author.FirstName.charAt(0)}${author.LastName.charAt(0)}`
    : '?'
  const isOwnComment = currentUserId === comment.user_id

  const handleSubmitReply = async () => {
    if (!replyContent.trim()) return
    setSubmittingReply(true)
    try {
      await onReply(comment.id, replyContent.trim())
      setReplyContent('')
      setShowReplyInput(false)
    } finally {
      setSubmittingReply(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSubmitReply()
    }
  }

  return (
    <div className={isNested ? '' : ''}>
      <div
        className={`flex gap-3 ${
          isNested ? 'ml-10 pl-4 border-l-2 border-muted' : ''
        }`}
      >
        {/* Avatar */}
        <div
          className={`flex-shrink-0 rounded-full flex items-center justify-center font-medium text-xs ${
            isOwnComment
              ? 'bg-primary text-primary-foreground'
              : 'bg-muted text-muted-foreground'
          } ${isNested ? 'h-7 w-7' : 'h-8 w-8'}`}
        >
          {author ? authorInitials : <UserIcon className="h-3.5 w-3.5" />}
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          {/* Header */}
          <div className="flex items-center gap-2 flex-wrap">
            <span className="text-sm font-medium">{authorName}</span>
            {isOwnComment && (
              <Badge variant="secondary" className="text-[10px] px-1.5 py-0">
                you
              </Badge>
            )}
            <span className="text-xs text-muted-foreground">
              {timeAgo(comment.created_at)}
            </span>
          </div>

          {/* Comment body */}
          <p className="text-sm mt-1 whitespace-pre-wrap break-words">
            {comment.content}
          </p>

          {/* Actions */}
          <div className="flex items-center gap-1 mt-2">
            <Button
              variant="ghost"
              size="sm"
              className="h-7 px-2 text-xs text-muted-foreground hover:text-foreground"
              onClick={() => onUpvote(comment.id)}
            >
              <ThumbsUp className="h-3.5 w-3.5 mr-1" />
              {comment.upvote || 0}
            </Button>
            <Button
              variant="ghost"
              size="sm"
              className="h-7 px-2 text-xs text-muted-foreground hover:text-foreground"
              onClick={() => onDownvote(comment.id)}
            >
              <ThumbsDown className="h-3.5 w-3.5 mr-1" />
              {comment.downvote || 0}
            </Button>
            {!isNested && (
              <Button
                variant="ghost"
                size="sm"
                className="h-7 px-2 text-xs text-muted-foreground hover:text-foreground"
                onClick={() => setShowReplyInput(!showReplyInput)}
              >
                <Reply className="h-3.5 w-3.5 mr-1" />
                Reply
              </Button>
            )}
          </div>

          {/* Reply input */}
          {showReplyInput && (
            <div className="mt-2 flex gap-2">
              <textarea
                value={replyContent}
                onChange={(e) => setReplyContent(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Write a reply..."
                className="flex-1 min-h-[60px] max-h-[120px] rounded-md border border-input bg-background px-3 py-2 text-sm resize-none focus:outline-none focus:ring-1 focus:ring-ring"
                disabled={submittingReply}
              />
              <Button
                size="sm"
                className="h-auto self-end"
                onClick={handleSubmitReply}
                disabled={!replyContent.trim() || submittingReply}
              >
                <Send className="h-3.5 w-3.5" />
              </Button>
            </div>
          )}

          {/* Nested replies */}
          {replies.length > 0 && (
            <div className="mt-3 space-y-3">
              {replies.map((reply) => (
                <CommentItem
                  key={reply.id}
                  comment={reply}
                  replies={[]}
                  usersMap={usersMap}
                  currentUserId={currentUserId}
                  onUpvote={onUpvote}
                  onDownvote={onDownvote}
                  onReply={onReply}
                  isNested
                />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default CommentItem
