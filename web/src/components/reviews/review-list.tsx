"use client"

import { Star, Plus, MessageCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { ReviewCard } from './review-card'
import type { Review } from '@/types'

interface ReviewListProps {
  reviews: Review[]
  averageRating: number | null
  isLoading: boolean
  error: string | null
  onAddReview: () => void
  className?: string
}

// Helper function to render star rating
const StarRating = ({ rating, className = "" }: { rating: number; className?: string }) => {
  return (
    <div className={`flex items-center gap-1 ${className}`}>
      {[1, 2, 3, 4, 5].map((star) => (
        <Star 
          key={star} 
          className={`w-4 h-4 ${
            star <= rating 
              ? 'fill-yellow-400 text-yellow-400' 
              : 'text-gray-300'
          }`}
        />
      ))}
    </div>
  )
}

export const ReviewList = ({ 
  reviews, 
  averageRating, 
  isLoading, 
  error,
  onAddReview,
  className = "" 
}: ReviewListProps) => {
  // Calculate average solo-friendly rating
  const averageSoloFriendlyRating = reviews.length > 0 
    ? reviews.reduce((acc, review) => acc + review.soloFriendlyRating, 0) / reviews.length 
    : null

  if (isLoading) {
    return (
      <div className={`space-y-4 ${className}`}>
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-24 mb-2"></div>
          <div className="h-3 bg-gray-200 rounded w-32"></div>
        </div>
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="animate-pulse">
              <div className="h-20 bg-gray-200 rounded"></div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className={`text-center py-8 ${className}`}>
        <MessageCircle className="w-8 h-8 text-gray-400 mx-auto mb-2" />
        <p className="text-sm text-gray-500 mb-4">{error}</p>
        <Button onClick={onAddReview} size="sm" className="bg-blue-600 hover:bg-blue-700">
          <Plus className="w-4 h-4 mr-2" />
          Add Review
        </Button>
      </div>
    )
  }

  if (reviews.length === 0) {
    return (
      <div className={`text-center py-8 ${className}`}>
        <MessageCircle className="w-8 h-8 text-gray-400 mx-auto mb-2" />
        <p className="text-sm text-gray-500 mb-4">No reviews yet. Be the first to share your experience!</p>
        <Button onClick={onAddReview} size="sm" className="bg-blue-600 hover:bg-blue-700">
          <Plus className="w-4 h-4 mr-2" />
          Write First Review
        </Button>
      </div>
    )
  }

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Reviews Summary */}
      <Card className="p-4">
        <div className="flex items-center justify-between mb-4">
          <h3 className="font-semibold text-gray-900">Reviews ({reviews.length})</h3>
          <Button onClick={onAddReview} size="sm" className="bg-blue-600 hover:bg-blue-700">
            <Plus className="w-4 h-4 mr-2" />
            Add Review
          </Button>
        </div>

        {/* Average ratings */}
        <div className="grid grid-cols-2 gap-4 mb-4">
          <div className="text-center">
            <div className="text-2xl font-bold text-gray-900 mb-1">
              {averageRating ? averageRating.toFixed(1) : '-'}
            </div>
            <StarRating rating={averageRating || 0} className="justify-center mb-1" />
            <div className="text-sm text-gray-500">Overall Rating</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-gray-900 mb-1">
              {averageSoloFriendlyRating ? averageSoloFriendlyRating.toFixed(1) : '-'}
            </div>
            <StarRating rating={averageSoloFriendlyRating || 0} className="justify-center mb-1" />
            <div className="text-sm text-gray-500">Solo-friendly</div>
          </div>
        </div>

        {/* Rating distribution */}
        <div className="space-y-2">
          {[5, 4, 3, 2, 1].map((rating) => {
            const count = reviews.filter(r => r.rating === rating).length
            const percentage = reviews.length > 0 ? (count / reviews.length) * 100 : 0
            
            return (
              <div key={rating} className="flex items-center gap-2 text-sm">
                <div className="flex items-center gap-1 w-8">
                  <span>{rating}</span>
                  <Star className="w-3 h-3 fill-yellow-400 text-yellow-400" />
                </div>
                <div className="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
                  <div 
                    className="h-full bg-yellow-400 transition-all duration-300"
                    style={{ width: `${percentage}%` }}
                  />
                </div>
                <div className="w-8 text-right text-gray-500">{count}</div>
              </div>
            )
          })}
        </div>
      </Card>

      {/* Reviews List */}
      <div className="space-y-3">
        {reviews.map((review) => (
          <ReviewCard key={review.id} review={review} />
        ))}
      </div>
    </div>
  )
}