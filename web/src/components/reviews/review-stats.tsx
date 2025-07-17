import React from "react"
import { Star, Users } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { RatingStars } from "./rating-stars"
import { cn } from "@/lib/utils"
import { Review } from "@/types"

interface ReviewStatsProps {
  reviews: Review[]
  className?: string
}

export const ReviewStats: React.FC<ReviewStatsProps> = ({
  reviews,
  className
}) => {
  const totalReviews = reviews.length
  
  const averageRating = totalReviews > 0 
    ? reviews.reduce((sum, review) => sum + review.rating, 0) / totalReviews
    : 0
    
  const averageSoloFriendlyRating = totalReviews > 0
    ? reviews.reduce((sum, review) => sum + review.soloFriendlyRating, 0) / totalReviews
    : 0

  // Calculate rating distribution
  const ratingDistribution = [5, 4, 3, 2, 1].map(rating => {
    const count = reviews.filter(review => Math.floor(review.rating) === rating).length
    const percentage = totalReviews > 0 ? (count / totalReviews) * 100 : 0
    return { rating, count, percentage }
  })

  if (totalReviews === 0) {
    return (
      <Card className={cn("w-full", className)}>
        <CardContent className="py-6">
          <div className="text-center text-muted-foreground">
            <Star className="w-12 h-12 mx-auto mb-2 opacity-30" />
            <p>No reviews yet</p>
            <p className="text-sm">Be the first to share your experience!</p>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className={cn("w-full", className)}>
      <CardHeader className="pb-3">
        <CardTitle className="text-lg flex items-center gap-2">
          <Star className="w-5 h-5 fill-yellow-400 text-yellow-400" />
          Review Summary
        </CardTitle>
      </CardHeader>
      
      <CardContent className="space-y-4">
        {/* Overall Stats */}
        <div className="grid grid-cols-2 gap-4">
          <div className="text-center">
            <div className="text-2xl font-bold">{averageRating.toFixed(1)}</div>
            <RatingStars rating={averageRating} showCount={false} />
            <div className="text-sm text-muted-foreground mt-1">
              Overall Rating
            </div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold">{averageSoloFriendlyRating.toFixed(1)}</div>
            <RatingStars rating={averageSoloFriendlyRating} showCount={false} />
            <div className="text-sm text-muted-foreground mt-1">
              Solo-Friendly
            </div>
          </div>
        </div>

        {/* Total Reviews */}
        <div className="flex items-center justify-center gap-2 py-2 border-y">
          <Users className="w-4 h-4 text-muted-foreground" />
          <span className="text-sm font-medium">
            {totalReviews} review{totalReviews !== 1 ? "s" : ""}
          </span>
        </div>

        {/* Rating Distribution */}
        <div className="space-y-2">
          <h4 className="text-sm font-medium">Rating Distribution</h4>
          {ratingDistribution.map(({ rating, count, percentage }) => (
            <div key={rating} className="flex items-center gap-2">
              <div className="flex items-center gap-1 w-12">
                <span className="text-sm">{rating}</span>
                <Star className="w-3 h-3 fill-yellow-400 text-yellow-400" />
              </div>
              <div className="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
                <div
                  className="h-full bg-yellow-400 transition-all duration-500"
                  style={{ width: `${percentage}%` }}
                />
              </div>
              <span className="text-sm text-muted-foreground w-8">
                {count}
              </span>
            </div>
          ))}
        </div>

        {/* Additional Stats */}
        <div className="grid grid-cols-2 gap-4 pt-2 border-t">
          <div className="text-center">
            <div className="text-lg font-semibold">
              {reviews.filter(r => r.photos && r.photos.length > 0).length}
            </div>
            <div className="text-xs text-muted-foreground">
              Reviews with photos
            </div>
          </div>
          <div className="text-center">
            <div className="text-lg font-semibold">
              {reviews.reduce((sum, r) => sum + r.helpful, 0)}
            </div>
            <div className="text-xs text-muted-foreground">
              Helpful votes
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}