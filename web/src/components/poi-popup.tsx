"use client"

import React, { useState, useEffect, useCallback } from "react"
import { X, MapPin, Star, Plus, AlertCircle, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ReviewForm, ReviewStats, type ReviewFormData } from "@/components/reviews"
import { useMapStore } from "@/stores/use-map-store"
import { api, handleAPIError } from "@/lib/api-client"
import { cn } from "@/lib/utils"
import { useUser } from "@clerk/nextjs"
// import type { Spot } from "@/types"

interface POIPopupProps {
  className?: string
}

export const POIPopup: React.FC<POIPopupProps> = ({ className }) => {
  const {
    popup,
    closePopup,
    setPopupReviews,
    setPopupReviewsLoading,
    setPopupReviewsError,
  } = useMapStore()

  const { user } = useUser()
  const [showReviewForm, setShowReviewForm] = useState(false)
  const [isSubmittingReview, setIsSubmittingReview] = useState(false)

  const fetchReviews = useCallback(async (spotId: string) => {
    setPopupReviewsLoading(true)
    setPopupReviewsError(null)
    
    try {
      const response = await api.reviews.list({ spot_id: spotId })
      if (response.error) {
        setPopupReviewsError(handleAPIError(response.error))
      } else {
        setPopupReviews(response.data || [])
      }
    } catch {
      setPopupReviewsError("Failed to load reviews")
    }
  }, [setPopupReviewsLoading, setPopupReviewsError, setPopupReviews])

  // Fetch reviews when popup opens
  useEffect(() => {
    if (popup.isOpen && popup.spot) {
      fetchReviews(popup.spot.id)
    }
  }, [popup.isOpen, popup.spot, fetchReviews])

  const handleSubmitReview = async (reviewData: ReviewFormData) => {
    if (!popup.spot) return

    setIsSubmittingReview(true)
    try {
      // Get current user from Clerk
      if (!user) {
        setPopupReviewsError("You must be logged in to submit a review")
        return
      }

      const currentUserId = user.id
      const currentUserName = user.fullName || user.primaryEmailAddress?.emailAddress.split('@')[0] || "Anonymous User"

      const reviewPayload = {
        spotId: reviewData.spotId,
        userId: currentUserId,
        userName: currentUserName,
        rating: reviewData.rating,
        soloFriendlyRating: reviewData.soloFriendlyRating,
        comment: reviewData.comment,
        tags: reviewData.tags,
        photos: [], // TODO: Handle photo uploads
        helpful: 0,
        notHelpful: 0,
      }

      const response = await api.reviews.create(reviewPayload)
      if (response.error) {
        setPopupReviewsError(handleAPIError(response.error))
      } else {
        // Refresh reviews after successful submission
        await fetchReviews(popup.spot.id)
        setShowReviewForm(false)
      }
    } catch {
      setPopupReviewsError("Failed to submit review")
    } finally {
      setIsSubmittingReview(false)
    }
  }

  if (!popup.isOpen || !popup.spot) {
    return null
  }

  return (
    <Card className={cn("w-full max-w-md shadow-lg", className)}>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
              <MapPin className="w-5 h-5 text-primary" />
            </div>
            <div className="flex-1 min-w-0">
              <CardTitle className="text-lg font-semibold text-foreground mb-1">
                {popup.spot.name}
              </CardTitle>
              <p className="text-sm text-muted-foreground mb-1">
                {popup.spot.category}
              </p>
              <p className="text-xs text-muted-foreground">
                {popup.spot.address}
              </p>
              {popup.averageRating && (
                <div className="flex items-center gap-1 mt-2">
                  <Star className="w-4 h-4 fill-yellow-400 text-yellow-400" />
                  <span className="text-sm font-medium">
                    {popup.averageRating.toFixed(1)}
                  </span>
                  <span className="text-xs text-muted-foreground">
                    ({popup.reviews.length} review{popup.reviews.length !== 1 ? "s" : ""})
                  </span>
                </div>
              )}
            </div>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={closePopup}
            className="flex-shrink-0 h-8 w-8 p-0"
          >
            <X className="w-4 h-4" />
          </Button>
        </div>
      </CardHeader>

      <CardContent className="pt-0">
        {/* Loading state */}
        {popup.isLoadingReviews && (
          <div className="flex items-center justify-center py-8">
            <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
            <span className="ml-2 text-sm text-muted-foreground">
              Loading reviews...
            </span>
          </div>
        )}

        {/* Error state */}
        {popup.reviewsError && (
          <div className="flex items-center gap-2 p-3 bg-destructive/10 text-destructive rounded-md mb-4">
            <AlertCircle className="w-4 h-4 flex-shrink-0" />
            <span className="text-sm">{popup.reviewsError}</span>
          </div>
        )}

        {/* Review form */}
        {showReviewForm && (
          <div className="mb-4">
            <ReviewForm
              spotId={popup.spot.id}
              onSubmit={handleSubmitReview}
              onCancel={() => setShowReviewForm(false)}
              className="border-none shadow-none"
            />
          </div>
        )}

        {/* Reviews content */}
        {!popup.isLoadingReviews && !popup.reviewsError && !showReviewForm && (
          <div className="space-y-4">
            {/* Review stats */}
            {popup.reviews.length > 0 && (
              <ReviewStats reviews={popup.reviews} />
            )}

            {/* Reviews list or empty state */}
            <div className="max-h-96 overflow-y-auto">
              {popup.reviews.length === 0 ? (
                <div className="text-center py-8">
                  <Star className="w-12 h-12 text-muted-foreground/30 mx-auto mb-4" />
                  <p className="text-sm text-muted-foreground mb-4">
                    No reviews yet. Be the first to share your experience!
                  </p>
                  {user ? (
                    <Button
                      onClick={() => setShowReviewForm(true)}
                      size="sm"
                      className="bg-primary hover:bg-primary/90"
                    >
                      <Plus className="w-4 h-4 mr-2" />
                      Write First Review
                    </Button>
                  ) : (
                    <div className="space-y-2">
                      <p className="text-xs text-muted-foreground">
                        Please log in to write a review
                      </p>
                      <Button
                        onClick={() => window.location.href = '/sign-in'}
                        variant="outline"
                        size="sm"
                      >
                        Log In
                      </Button>
                    </div>
                  )}
                </div>
              ) : (
                <div className="space-y-3">
                  {popup.reviews.map((review) => (
                    <div
                      key={review.id}
                      className="p-3 bg-muted/30 rounded-lg border"
                    >
                      <div className="flex items-start justify-between mb-2">
                        <div className="flex items-center gap-2">
                          {review.userAvatar && (
                            <img
                              src={review.userAvatar}
                              alt={review.userName}
                              className="w-6 h-6 rounded-full object-cover"
                            />
                          )}
                          <div>
                            <p className="text-sm font-medium text-foreground">
                              {review.userName}
                            </p>
                            <p className="text-xs text-muted-foreground">
                              {new Date(review.createdAt).toLocaleDateString()}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center gap-1">
                          <Star className="w-3 h-3 fill-yellow-400 text-yellow-400" />
                          <span className="text-sm font-medium">
                            {review.rating.toFixed(1)}
                          </span>
                        </div>
                      </div>
                      {review.comment && (
                        <p className="text-sm text-foreground mb-2">
                          {review.comment}
                        </p>
                      )}
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-1 text-xs text-muted-foreground">
                          <span>Solo-friendly:</span>
                          <Star className="w-3 h-3 fill-yellow-400 text-yellow-400" />
                          <span>{review.soloFriendlyRating.toFixed(1)}</span>
                        </div>
                        {review.tags && review.tags.length > 0 && (
                          <div className="flex gap-1">
                            {review.tags.slice(0, 2).map((tag, index) => (
                              <span
                                key={index}
                                className="px-2 py-1 text-xs bg-secondary text-secondary-foreground rounded"
                              >
                                {tag}
                              </span>
                            ))}
                            {review.tags.length > 2 && (
                              <span className="text-xs text-muted-foreground">
                                +{review.tags.length - 2}
                              </span>
                            )}
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                  
                  {/* Add review button */}
                  <div className="pt-2 border-t">
                    {user ? (
                      <Button
                        onClick={() => setShowReviewForm(true)}
                        variant="outline"
                        size="sm"
                        className="w-full"
                        disabled={isSubmittingReview}
                      >
                        <Plus className="w-4 h-4 mr-2" />
                        Add Review
                      </Button>
                    ) : (
                      <Button
                        onClick={() => window.location.href = '/sign-in'}
                        variant="outline"
                        size="sm"
                        className="w-full"
                      >
                        Log In to Add Review
                      </Button>
                    )}
                  </div>
                </div>
              )}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}