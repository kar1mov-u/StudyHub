import React, { useState, useEffect, useCallback, useMemo } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { MessageSquare, Send, Loader2 } from 'lucide-react'
import { commentsApi } from '@/api/comments'
import { authApi } from '@/api/auth'
import { useAuth } from '@/context/AuthContext'
import { useToast } from '@/components/ui/toast'
import CommentItem from './CommentItem'
import type { Comment, User } from '@/types'

interface CommentSectionProps {
  weekId: string
}

const CommentSection: React.FC<CommentSectionProps> = ({ weekId }) => {
  const { user } = useAuth()
  const { showToast } = useToast()

  const [comments, setComments] = useState<Comment[]>([])
  const [usersMap, setUsersMap] = useState<Record<string, User>>({})
  const [loading, setLoading] = useState(true)
  const [newComment, setNewComment] = useState('')
  const [submitting, setSubmitting] = useState(false)

  // Fetch user details for unique user IDs
  const fetchUsers = useCallback(
    async (commentsList: Comment[]) => {
      const uniqueUserIds = [
        ...new Set(commentsList.map((c) => c.user_id)),
      ].filter((id) => id && !usersMap[id])

      if (uniqueUserIds.length === 0) return

      const results = await Promise.allSettled(
        uniqueUserIds.map((id) => authApi.getCurrentUser(id))
      )

      const newMap: Record<string, User> = { ...usersMap }
      results.forEach((result, index) => {
        if (result.status === 'fulfilled' && result.value) {
          newMap[uniqueUserIds[index]] = result.value
        }
      })
      setUsersMap(newMap)
    },
    [usersMap]
  )

  // Load comments
  const loadComments = useCallback(async () => {
    try {
      setLoading(true)
      const data = await commentsApi.getCommentsByWeek(weekId)
      const commentsList = Array.isArray(data) ? data : []
      setComments(commentsList)
      await fetchUsers(commentsList)
    } catch {
      // Silently handle - empty comment list is fine on error
      setComments([])
    } finally {
      setLoading(false)
    }
  }, [weekId])

  useEffect(() => {
    loadComments()
  }, [weekId])

  // Refetch user details when comments change
  useEffect(() => {
    if (comments.length > 0) {
      fetchUsers(comments)
    }
  }, [comments])

  // Separate top-level comments from replies
  const { topLevelComments, repliesByParent } = useMemo(() => {
    const top: Comment[] = []
    const replies: Record<string, Comment[]> = {}

    for (const comment of comments) {
      if (comment.reply_id) {
        if (!replies[comment.reply_id]) {
          replies[comment.reply_id] = []
        }
        replies[comment.reply_id].push(comment)
      } else {
        top.push(comment)
      }
    }

    // Sort top-level by newest first
    top.sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    )

    // Sort replies by oldest first (chronological within a thread)
    for (const parentId of Object.keys(replies)) {
      replies[parentId].sort(
        (a, b) =>
          new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
      )
    }

    return { topLevelComments: top, repliesByParent: replies }
  }, [comments])

  // Post a new comment
  const handlePostComment = async () => {
    if (!newComment.trim() || !user) return

    setSubmitting(true)
    try {
      await commentsApi.createComment({
        user_id: user.ID,
        week_id: weekId,
        content: newComment.trim(),
      })
      setNewComment('')
      await loadComments()
      showToast('Comment posted', 'success')
    } catch {
      showToast('Failed to post comment', 'error')
    } finally {
      setSubmitting(false)
    }
  }

  // Post a reply
  const handleReply = async (parentId: string, content: string) => {
    if (!user) return

    try {
      await commentsApi.createComment({
        user_id: user.ID,
        week_id: weekId,
        content,
        reply_id: parentId,
      })
      await loadComments()
      showToast('Reply posted', 'success')
    } catch {
      showToast('Failed to post reply', 'error')
    }
  }

  // Upvote
  const handleUpvote = async (commentId: string) => {
    try {
      await commentsApi.upvoteComment(commentId)
      // Optimistically update the local state
      setComments((prev) =>
        prev.map((c) =>
          c.id === commentId ? { ...c, upvote: (c.upvote || 0) + 1 } : c
        )
      )
    } catch {
      showToast('Failed to upvote', 'error')
    }
  }

  // Downvote
  const handleDownvote = async (commentId: string) => {
    try {
      await commentsApi.downvoteComment(commentId)
      setComments((prev) =>
        prev.map((c) =>
          c.id === commentId ? { ...c, downvote: (c.downvote || 0) + 1 } : c
        )
      )
    } catch {
      showToast('Failed to downvote', 'error')
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handlePostComment()
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-lg">
          <MessageSquare className="h-5 w-5" />
          Discussion
          {comments.length > 0 && (
            <span className="text-sm font-normal text-muted-foreground">
              ({comments.length})
            </span>
          )}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* New comment input */}
        {user && (
          <div className="flex gap-3">
            <div className="flex-shrink-0 h-8 w-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-medium text-xs">
              {user.FirstName.charAt(0)}
              {user.LastName.charAt(0)}
            </div>
            <div className="flex-1 space-y-2">
              <textarea
                value={newComment}
                onChange={(e) => setNewComment(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Share your thoughts or ask a question..."
                className="w-full min-h-[80px] max-h-[200px] rounded-md border border-input bg-background px-3 py-2 text-sm resize-none focus:outline-none focus:ring-1 focus:ring-ring placeholder:text-muted-foreground"
                disabled={submitting}
              />
              <div className="flex items-center justify-between">
                <p className="text-xs text-muted-foreground">
                  Press Enter to post, Shift+Enter for new line
                </p>
                <Button
                  size="sm"
                  onClick={handlePostComment}
                  disabled={!newComment.trim() || submitting}
                >
                  {submitting ? (
                    <Loader2 className="h-3.5 w-3.5 mr-1 animate-spin" />
                  ) : (
                    <Send className="h-3.5 w-3.5 mr-1" />
                  )}
                  Post
                </Button>
              </div>
            </div>
          </div>
        )}

        {/* Divider */}
        {user && comments.length > 0 && <div className="border-t" />}

        {/* Comments list */}
        {loading ? (
          <div className="flex items-center justify-center py-8">
            <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
          </div>
        ) : topLevelComments.length === 0 ? (
          <div className="text-center py-8">
            <MessageSquare className="h-8 w-8 text-muted-foreground mx-auto mb-2 opacity-50" />
            <p className="text-sm text-muted-foreground">
              No comments yet. Be the first to start a discussion!
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {topLevelComments.map((comment) => (
              <CommentItem
                key={comment.id}
                comment={comment}
                replies={repliesByParent[comment.id] || []}
                usersMap={usersMap}
                currentUserId={user?.ID}
                onUpvote={handleUpvote}
                onDownvote={handleDownvote}
                onReply={handleReply}
              />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}

export default CommentSection
