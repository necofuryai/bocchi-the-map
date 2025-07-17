import React, { useState } from "react"
import { Camera, X, Plus } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { RatingStars } from "./rating-stars"
import { cn } from "@/lib/utils"

interface ReviewFormProps {
  spotId: string
  onSubmit: (reviewData: ReviewFormData) => void
  onCancel?: () => void
  className?: string
}

export interface ReviewFormData {
  spotId: string
  rating: number
  soloFriendlyRating: number
  comment: string
  tags: string[]
  photos: File[]
}

export const ReviewForm: React.FC<ReviewFormProps> = ({
  spotId,
  onSubmit,
  onCancel,
  className
}) => {
  const [rating, setRating] = useState(0)
  const [soloFriendlyRating, setSoloFriendlyRating] = useState(0)
  const [comment, setComment] = useState("")
  const [tags, setTags] = useState<string[]>([])
  const [tagInput, setTagInput] = useState("")
  const [photos, setPhotos] = useState<File[]>([])
  const [isSubmitting, setIsSubmitting] = useState(false)

  const commonTags = [
    "Solo-friendly",
    "Quiet",
    "Scenic",
    "Crowded",
    "Peaceful",
    "Instagram-worthy",
    "Hidden gem",
    "Tourist spot",
    "Good for photos",
    "Relaxing"
  ]

  const handleTagAdd = (tag: string) => {
    if (tag.trim() && !tags.includes(tag.trim())) {
      setTags([...tags, tag.trim()])
      setTagInput("")
    }
  }

  const handleTagRemove = (tagToRemove: string) => {
    setTags(tags.filter(tag => tag !== tagToRemove))
  }

  const handlePhotoUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (files) {
      const newPhotos = Array.from(files).slice(0, 5 - photos.length)
      setPhotos([...photos, ...newPhotos])
    }
  }

  const handlePhotoRemove = (index: number) => {
    setPhotos(photos.filter((_, i) => i !== index))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (rating === 0 || soloFriendlyRating === 0) {
      return
    }

    setIsSubmitting(true)
    try {
      await onSubmit({
        spotId,
        rating,
        soloFriendlyRating,
        comment,
        tags,
        photos
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Card className={cn("w-full", className)}>
      <CardHeader>
        <CardTitle className="text-lg">Write a Review</CardTitle>
      </CardHeader>
      
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Overall Rating */}
          <div>
            <label className="block text-sm font-medium mb-2">
              Overall Rating *
            </label>
            <RatingStars
              rating={rating}
              interactive
              onRatingChange={setRating}
              size="lg"
            />
          </div>

          {/* Solo-Friendly Rating */}
          <div>
            <label className="block text-sm font-medium mb-2">
              Solo-Friendly Rating *
            </label>
            <RatingStars
              rating={soloFriendlyRating}
              interactive
              onRatingChange={setSoloFriendlyRating}
              size="lg"
            />
          </div>

          {/* Comment */}
          <div>
            <label className="block text-sm font-medium mb-2">
              Your Review
            </label>
            <textarea
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              placeholder="Share your experience..."
              className="w-full min-h-[100px] p-3 border rounded-md resize-none focus:outline-none focus:ring-2 focus:ring-primary"
              maxLength={500}
            />
            <div className="text-xs text-muted-foreground mt-1">
              {comment.length}/500 characters
            </div>
          </div>

          {/* Tags */}
          <div>
            <label className="block text-sm font-medium mb-2">
              Tags
            </label>
            <div className="flex flex-wrap gap-2 mb-2">
              {commonTags.map((tag) => (
                <button
                  key={tag}
                  type="button"
                  onClick={() => handleTagAdd(tag)}
                  className={cn(
                    "px-3 py-1 rounded-full text-xs border transition-colors",
                    tags.includes(tag)
                      ? "bg-primary text-primary-foreground border-primary"
                      : "bg-background border-input hover:bg-accent"
                  )}
                >
                  {tag}
                </button>
              ))}
            </div>
            <div className="flex gap-2">
              <input
                type="text"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyPress={(e) => {
                  if (e.key === "Enter") {
                    e.preventDefault()
                    handleTagAdd(tagInput)
                  }
                }}
                placeholder="Add custom tag..."
                className="flex-1 px-3 py-2 border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => handleTagAdd(tagInput)}
              >
                <Plus className="w-4 h-4" />
              </Button>
            </div>
            {tags.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-2">
                {tags.map((tag) => (
                  <span
                    key={tag}
                    className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs bg-secondary text-secondary-foreground"
                  >
                    {tag}
                    <button
                      type="button"
                      onClick={() => handleTagRemove(tag)}
                      className="hover:text-destructive"
                    >
                      <X className="w-3 h-3" />
                    </button>
                  </span>
                ))}
              </div>
            )}
          </div>

          {/* Photos */}
          <div>
            <label className="block text-sm font-medium mb-2">
              Photos (up to 5)
            </label>
            <div className="grid grid-cols-3 gap-2">
              {photos.map((photo, index) => (
                <div key={index} className="relative aspect-square rounded-lg overflow-hidden">
                  <img
                    src={URL.createObjectURL(photo)}
                    alt={`Upload ${index + 1}`}
                    className="w-full h-full object-cover"
                  />
                  <button
                    type="button"
                    onClick={() => handlePhotoRemove(index)}
                    className="absolute top-1 right-1 w-6 h-6 bg-black/50 rounded-full flex items-center justify-center"
                  >
                    <X className="w-3 h-3 text-white" />
                  </button>
                </div>
              ))}
              {photos.length < 5 && (
                <label className="aspect-square border-2 border-dashed border-input rounded-lg flex flex-col items-center justify-center cursor-pointer hover:bg-accent transition-colors">
                  <Camera className="w-6 h-6 text-muted-foreground mb-1" />
                  <span className="text-xs text-muted-foreground">Add Photo</span>
                  <input
                    type="file"
                    accept="image/*"
                    multiple
                    onChange={handlePhotoUpload}
                    className="hidden"
                  />
                </label>
              )}
            </div>
          </div>

          {/* Submit Buttons */}
          <div className="flex gap-2 pt-4">
            <Button
              type="submit"
              disabled={rating === 0 || soloFriendlyRating === 0 || isSubmitting}
              className="flex-1"
            >
              {isSubmitting ? "Submitting..." : "Submit Review"}
            </Button>
            {onCancel && (
              <Button
                type="button"
                variant="outline"
                onClick={onCancel}
                disabled={isSubmitting}
              >
                Cancel
              </Button>
            )}
          </div>
        </form>
      </CardContent>
    </Card>
  )
}