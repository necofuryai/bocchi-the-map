import { http, HttpResponse } from 'msw'

// Mock review data
const mockReviews = [
  {
    id: 'review-1',
    spotId: '1',
    userId: 'user-1',
    userName: 'Test User',
    userAvatar: '/images/avatar1.jpg',
    rating: 5,
    soloFriendlyRating: 5,
    comment: 'Perfect spot for solo work! Very quiet and peaceful.',
    tags: ['quiet', 'wifi', 'solo-friendly'],
    photos: ['/images/review1.jpg'],
    helpful: 12,
    notHelpful: 1,
    createdAt: '2024-06-01T10:00:00Z',
    updatedAt: '2024-06-01T10:00:00Z',
  },
  {
    id: 'review-2',
    spotId: '1',
    userId: 'user-2',
    userName: 'Solo Traveler',
    userAvatar: '/images/avatar2.jpg',
    rating: 4,
    soloFriendlyRating: 5,
    comment: 'Great coffee and atmosphere. Staff is very respectful of solo diners.',
    tags: ['coffee', 'respectful-staff', 'comfortable-seating'],
    photos: [],
    helpful: 8,
    notHelpful: 0,
    createdAt: '2024-05-15T14:30:00Z',
    updatedAt: '2024-05-15T14:30:00Z',
  },
  {
    id: 'review-3',
    spotId: '2',
    userId: 'user-1',
    userName: 'Test User',
    userAvatar: '/images/avatar1.jpg',
    rating: 3,
    soloFriendlyRating: 2,
    comment: 'Too crowded and noisy for solo activities. Better for groups.',
    tags: ['crowded', 'noisy', 'social'],
    photos: [],
    helpful: 5,
    notHelpful: 2,
    createdAt: '2024-05-20T16:45:00Z',
    updatedAt: '2024-05-20T16:45:00Z',
  },
]

// Helper function to get user from token
function getUserFromToken(authHeader: string | null): { userId: string } | null {
  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return null
  }
  
  const token = authHeader.substring(7)
  // In a real app, verify JWT token
  const match = token.match(/mock-token-(.+?)-/)
  return match ? { userId: match[1] } : null
}

export const reviewHandlers = [
  // Get reviews for a spot
  http.get('/api/spots/:spotId/reviews', ({ params, request }) => {
    const url = new URL(request.url)
    const limit = parseInt(url.searchParams.get('limit') || '10')
    const offset = parseInt(url.searchParams.get('offset') || '0')
    const sortBy = url.searchParams.get('sort_by') || 'created_at'
    const order = url.searchParams.get('order') || 'desc'

    let spotReviews = mockReviews.filter(review => review.spotId === params.spotId)

    // Sort reviews
    spotReviews.sort((a, b) => {
      let valueA: any, valueB: any
      
      switch (sortBy) {
        case 'rating':
          valueA = a.rating
          valueB = b.rating
          break
        case 'solo_friendly_rating':
          valueA = a.soloFriendlyRating
          valueB = b.soloFriendlyRating
          break
        case 'helpful':
          valueA = a.helpful
          valueB = b.helpful
          break
        case 'created_at':
        default:
          valueA = new Date(a.createdAt)
          valueB = new Date(b.createdAt)
      }

      if (order === 'asc') {
        return valueA > valueB ? 1 : -1
      } else {
        return valueA < valueB ? 1 : -1
      }
    })

    // Pagination
    const paginatedReviews = spotReviews.slice(offset, offset + limit)
    const hasMore = offset + limit < spotReviews.length

    // Calculate statistics
    const totalReviews = spotReviews.length
    const averageRating = totalReviews > 0 
      ? spotReviews.reduce((sum, r) => sum + r.rating, 0) / totalReviews 
      : 0
    const averageSoloFriendlyRating = totalReviews > 0
      ? spotReviews.reduce((sum, r) => sum + r.soloFriendlyRating, 0) / totalReviews
      : 0

    return HttpResponse.json({
      data: paginatedReviews,
      total: totalReviews,
      hasMore,
      offset,
      limit,
      statistics: {
        averageRating: Math.round(averageRating * 10) / 10,
        averageSoloFriendlyRating: Math.round(averageSoloFriendlyRating * 10) / 10,
        totalReviews,
      },
    })
  }),

  // Create a new review
  http.post('/api/spots/:spotId/reviews', async ({ params, request }) => {
    const user = getUserFromToken(request.headers.get('authorization'))
    
    if (!user) {
      return HttpResponse.json(
        { error: 'Authentication required' },
        { status: 401 }
      )
    }

    const body = await request.json() as any
    
    // Validation
    if (!body.rating || !body.soloFriendlyRating || !body.comment) {
      return HttpResponse.json(
        { error: 'Rating, solo-friendly rating, and comment are required' },
        { status: 400 }
      )
    }

    if (body.rating < 1 || body.rating > 5 || body.soloFriendlyRating < 1 || body.soloFriendlyRating > 5) {
      return HttpResponse.json(
        { error: 'Ratings must be between 1 and 5' },
        { status: 400 }
      )
    }

    // Check if user already reviewed this spot
    const existingReview = mockReviews.find(
      r => r.spotId === params.spotId && r.userId === user.userId
    )
    
    if (existingReview) {
      return HttpResponse.json(
        { error: 'You have already reviewed this spot' },
        { status: 409 }
      )
    }

    const newReview = {
      id: `review-${mockReviews.length + 1}`,
      spotId: params.spotId as string,
      userId: user.userId,
      userName: 'Test User', // In real app, get from user data
      userAvatar: '/images/avatar1.jpg',
      rating: body.rating,
      soloFriendlyRating: body.soloFriendlyRating,
      comment: body.comment,
      tags: body.tags || [],
      photos: body.photos || [],
      helpful: 0,
      notHelpful: 0,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }

    mockReviews.push(newReview)

    return HttpResponse.json(
      { data: newReview },
      { status: 201 }
    )
  }),

  // Get a specific review
  http.get('/api/reviews/:reviewId', ({ params }) => {
    const review = mockReviews.find(r => r.id === params.reviewId)
    
    if (!review) {
      return HttpResponse.json(
        { error: 'Review not found' },
        { status: 404 }
      )
    }

    return HttpResponse.json({ data: review })
  }),

  // Update a review
  http.put('/api/reviews/:reviewId', async ({ params, request }) => {
    const user = getUserFromToken(request.headers.get('authorization'))
    
    if (!user) {
      return HttpResponse.json(
        { error: 'Authentication required' },
        { status: 401 }
      )
    }

    const reviewIndex = mockReviews.findIndex(r => r.id === params.reviewId)
    
    if (reviewIndex === -1) {
      return HttpResponse.json(
        { error: 'Review not found' },
        { status: 404 }
      )
    }

    const review = mockReviews[reviewIndex]
    
    if (review.userId !== user.userId) {
      return HttpResponse.json(
        { error: 'You can only edit your own reviews' },
        { status: 403 }
      )
    }

    const body = await request.json() as any
    
    const updatedReview = {
      ...review,
      ...body,
      updatedAt: new Date().toISOString(),
    }

    mockReviews[reviewIndex] = updatedReview

    return HttpResponse.json({ data: updatedReview })
  }),

  // Delete a review
  http.delete('/api/reviews/:reviewId', ({ params, request }) => {
    const user = getUserFromToken(request.headers.get('authorization'))
    
    if (!user) {
      return HttpResponse.json(
        { error: 'Authentication required' },
        { status: 401 }
      )
    }

    const reviewIndex = mockReviews.findIndex(r => r.id === params.reviewId)
    
    if (reviewIndex === -1) {
      return HttpResponse.json(
        { error: 'Review not found' },
        { status: 404 }
      )
    }

    const review = mockReviews[reviewIndex]
    
    if (review.userId !== user.userId) {
      return HttpResponse.json(
        { error: 'You can only delete your own reviews' },
        { status: 403 }
      )
    }

    mockReviews.splice(reviewIndex, 1)

    return HttpResponse.json({ success: true })
  }),

  // Mark review as helpful/not helpful
  http.post('/api/reviews/:reviewId/helpful', async ({ params, request }) => {
    const user = getUserFromToken(request.headers.get('authorization'))
    
    if (!user) {
      return HttpResponse.json(
        { error: 'Authentication required' },
        { status: 401 }
      )
    }

    const body = await request.json() as any
    const isHelpful = body.helpful === true

    const reviewIndex = mockReviews.findIndex(r => r.id === params.reviewId)
    
    if (reviewIndex === -1) {
      return HttpResponse.json(
        { error: 'Review not found' },
        { status: 404 }
      )
    }

    const review = mockReviews[reviewIndex]
    
    // Update helpful count
    if (isHelpful) {
      review.helpful += 1
    } else {
      review.notHelpful += 1
    }

    mockReviews[reviewIndex] = review

    return HttpResponse.json({ data: review })
  }),

  // Get user's reviews
  http.get('/api/users/:userId/reviews', ({ params, request }) => {
    const user = getUserFromToken(request.headers.get('authorization'))
    
    // Allow users to see their own reviews, or make public reviews visible
    if (user && user.userId !== params.userId) {
      // In a real app, check if profile is public
    }

    const url = new URL(request.url)
    const limit = parseInt(url.searchParams.get('limit') || '10')
    const offset = parseInt(url.searchParams.get('offset') || '0')

    const userReviews = mockReviews
      .filter(review => review.userId === params.userId)
      .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())

    const paginatedReviews = userReviews.slice(offset, offset + limit)
    const hasMore = offset + limit < userReviews.length

    return HttpResponse.json({
      data: paginatedReviews,
      total: userReviews.length,
      hasMore,
      offset,
      limit,
    })
  }),
]