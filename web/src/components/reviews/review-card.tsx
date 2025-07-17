import React from "react"
import { ThumbsUp, ThumbsDown, Camera } from "lucide-react"
import { Card, CardContent, CardHeader } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { RatingStars } from "./rating-stars"
import { cn } from "@/lib/utils"
import { Review } from "@/types"

interface ReviewCardProps {
  review: Review
  onHelpfulClick?: (reviewId: string) => void
  onNotHelpfulClick?: (reviewId: string) => void
  className?: string
}

export const ReviewCard: React.FC<ReviewCardProps> = ({
  review,
  onHelpfulClick,
  onNotHelpfulClick,
  className
}) => {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  }

  return (
    <Card className={cn("w-full", className)}>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-3">
            {review.userAvatar && (
              <img
                src={review.userAvatar}
                alt={review.userName}
                className="w-8 h-8 rounded-full object-cover"
              />
            )}
            <div>
              <h4 className="font-semibold text-sm">{review.userName}</h4>
              <p className="text-xs text-muted-foreground">
                {formatDate(review.createdAt)}
              </p>
            </div>
          </div>
          <div className="flex flex-col items-end gap-1">
            <RatingStars rating={review.rating} size="sm" />
            <div className="flex items-center gap-1">
              <span className="text-xs text-muted-foreground">Solo-friendly:</span>
              <RatingStars rating={review.soloFriendlyRating} size="sm" />
            </div>
          </div>
        </div>
      </CardHeader>
      
      <CardContent className="pt-0">
        {review.comment && (
          <p className="text-sm text-foreground mb-3 leading-relaxed">
            {review.comment}
          </p>
        )}
        
        {review.tags && review.tags.length > 0 && (
          <div className="flex flex-wrap gap-1 mb-3">
            {review.tags.map((tag, index) => (
              <span
                key={index}
                className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-secondary text-secondary-foreground"
              >
                {tag}
              </span>
            ))}
          </div>
        )}
        
        {review.photos && review.photos.length > 0 && (
          <div className="grid grid-cols-2 gap-2 mb-3">
            {review.photos.slice(0, 4).map((photo, index) => (
              <div key={index} className="relative aspect-square rounded-lg overflow-hidden">
                <img
                  src={photo}
                  alt={`Review photo ${index + 1}`}
                  className="w-full h-full object-cover"
                />
                {index === 3 && review.photos!.length > 4 && (
                  <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                    <span className="text-white text-sm font-semibold">
                      +{review.photos!.length - 4}
                    </span>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
        
        <div className="flex items-center justify-between mt-3">
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onHelpfulClick?.(review.id)}
              className="h-8 text-xs"
            >
              <ThumbsUp className="w-3 h-3 mr-1" />
              Helpful ({review.helpful})
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onNotHelpfulClick?.(review.id)}
              className="h-8 text-xs"
            >
              <ThumbsDown className="w-3 h-3 mr-1" />
              Not helpful ({review.notHelpful})
            </Button>
          </div>
          
          {review.photos && review.photos.length > 0 && (
            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <Camera className="w-3 h-3" />
              {review.photos.length}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}