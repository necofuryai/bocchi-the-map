import React from "react"
import { Star } from "lucide-react"
import { cn } from "@/lib/utils"

interface RatingStarsProps {
  rating: number
  maxStars?: number
  size?: "sm" | "md" | "lg"
  interactive?: boolean
  onRatingChange?: (rating: number) => void
  showCount?: boolean
  className?: string
}

export const RatingStars: React.FC<RatingStarsProps> = ({
  rating,
  maxStars = 5,
  size = "md",
  interactive = false,
  onRatingChange,
  showCount = false,
  className
}) => {
  const sizeClasses = {
    sm: "h-3 w-3",
    md: "h-4 w-4",
    lg: "h-5 w-5"
  }

  const handleStarClick = (starIndex: number) => {
    if (interactive && onRatingChange) {
      onRatingChange(starIndex + 1)
    }
  }

  const stars = Array.from({ length: maxStars }, (_, index) => {
    const isFilled = index < Math.floor(rating)
    const isHalfFilled = index < rating && index >= Math.floor(rating)

    return (
      <button
        key={index}
        type="button"
        onClick={() => handleStarClick(index)}
        disabled={!interactive}
        className={cn(
          "relative",
          interactive && "cursor-pointer hover:scale-110 transition-transform",
          !interactive && "cursor-default"
        )}
      >
        <Star
          className={cn(
            sizeClasses[size],
            isFilled || isHalfFilled
              ? "fill-yellow-400 text-yellow-400"
              : "fill-gray-200 text-gray-200"
          )}
        />
        {isHalfFilled && (
          <Star
            className={cn(
              sizeClasses[size],
              "absolute inset-0 fill-yellow-400 text-yellow-400"
            )}
            style={{
              clipPath: `polygon(0 0, ${((rating - index) * 100)}% 0, ${((rating - index) * 100)}% 100%, 0 100%)`
            }}
          />
        )}
      </button>
    )
  })

  return (
    <div className={cn("flex items-center gap-1", className)}>
      <div className="flex items-center gap-0.5">
        {stars}
      </div>
      {showCount && (
        <span className="text-sm text-muted-foreground ml-1">
          ({rating.toFixed(1)})
        </span>
      )}
    </div>
  )
}