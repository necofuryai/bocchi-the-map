import { http, HttpResponse } from 'msw'

// Mock data for spots
const mockSpots = [
  {
    id: '1',
    name: 'Quiet Coffee House',
    type: 'cafe',
    address: '123 Peaceful St, Tokyo',
    latitude: 35.6762,
    longitude: 139.6503,
    soloFriendly: true,
    soloFriendlyRating: 4.8,
    averageRating: 4.5,
    reviewCount: 127,
    amenities: ['wifi', 'quiet', 'power_outlets', 'solo_seating'],
    description: 'A peaceful cafe perfect for solo work and relaxation',
    photos: ['/images/cafe1.jpg'],
    openingHours: {
      monday: '08:00-20:00',
      tuesday: '08:00-20:00',
      wednesday: '08:00-20:00',
      thursday: '08:00-20:00',
      friday: '08:00-20:00',
      saturday: '09:00-21:00',
      sunday: '09:00-21:00',
    },
  },
  {
    id: '2',
    name: 'Busy Downtown Cafe',
    type: 'cafe',
    address: '456 Crowded Ave, Tokyo',
    latitude: 35.6895,
    longitude: 139.6917,
    soloFriendly: false,
    soloFriendlyRating: 2.1,
    averageRating: 4.2,
    reviewCount: 89,
    amenities: ['wifi', 'busy', 'social'],
    description: 'A lively cafe with lots of social interaction',
    photos: ['/images/cafe2.jpg'],
    openingHours: {
      monday: '07:00-22:00',
      tuesday: '07:00-22:00',
      wednesday: '07:00-22:00',
      thursday: '07:00-22:00',
      friday: '07:00-23:00',
      saturday: '08:00-23:00',
      sunday: '08:00-21:00',
    },
  },
  {
    id: '3',
    name: 'Silent Study Library',
    type: 'library',
    address: '789 Knowledge Blvd, Tokyo',
    latitude: 35.7028,
    longitude: 139.7753,
    soloFriendly: true,
    soloFriendlyRating: 4.9,
    averageRating: 4.7,
    reviewCount: 203,
    amenities: ['wifi', 'very_quiet', 'power_outlets', 'study_space'],
    description: 'Perfect quiet space for focused solo work',
    photos: ['/images/library1.jpg'],
    openingHours: {
      monday: '09:00-21:00',
      tuesday: '09:00-21:00',
      wednesday: '09:00-21:00',
      thursday: '09:00-21:00',
      friday: '09:00-19:00',
      saturday: '10:00-18:00',
      sunday: 'closed',
    },
  },
]

export const spotHandlers = [
  // Search spots
  http.get('/api/spots/search', ({ request }) => {
    const url = new URL(request.url)
    const query = url.searchParams.get('q')?.toLowerCase() || ''
    const soloFriendlyOnly = url.searchParams.get('solo_friendly') === 'true'
    const type = url.searchParams.get('type')
    const limit = parseInt(url.searchParams.get('limit') || '10')
    const offset = parseInt(url.searchParams.get('offset') || '0')

    let filteredSpots = mockSpots

    // Filter by search query
    if (query) {
      filteredSpots = filteredSpots.filter(spot =>
        spot.name.toLowerCase().includes(query) ||
        spot.description.toLowerCase().includes(query) ||
        spot.amenities.some(amenity => amenity.includes(query.replace(' ', '_')))
      )
    }

    // Filter by solo-friendly
    if (soloFriendlyOnly) {
      filteredSpots = filteredSpots.filter(spot => spot.soloFriendly)
    }

    // Filter by type
    if (type) {
      filteredSpots = filteredSpots.filter(spot => spot.type === type)
    }

    // Pagination
    const paginatedSpots = filteredSpots.slice(offset, offset + limit)
    const hasMore = offset + limit < filteredSpots.length

    return HttpResponse.json({
      data: paginatedSpots,
      total: filteredSpots.length,
      hasMore,
      offset,
      limit,
    })
  }),

  // Get spot by ID
  http.get('/api/spots/:id', ({ params }) => {
    const spot = mockSpots.find(s => s.id === params.id)
    
    if (!spot) {
      return HttpResponse.json(
        { error: 'Spot not found' },
        { status: 404 }
      )
    }

    return HttpResponse.json({ data: spot })
  }),

  // Get nearby spots
  http.get('/api/spots/nearby', ({ request }) => {
    const url = new URL(request.url)
    const lat = parseFloat(url.searchParams.get('lat') || '0')
    const lng = parseFloat(url.searchParams.get('lng') || '0')
    const radius = parseFloat(url.searchParams.get('radius') || '1000') // meters
    const limit = parseInt(url.searchParams.get('limit') || '10')

    // Simple distance calculation (not accurate for production)
    const nearbySpots = mockSpots
      .map(spot => ({
        ...spot,
        distance: Math.sqrt(
          Math.pow(spot.latitude - lat, 2) + Math.pow(spot.longitude - lng, 2)
        ) * 111000 // Rough conversion to meters
      }))
      .filter(spot => spot.distance <= radius)
      .sort((a, b) => a.distance - b.distance)
      .slice(0, limit)

    return HttpResponse.json({
      data: nearbySpots,
      center: { latitude: lat, longitude: lng },
      radius,
    })
  }),

  // Create new spot (requires authentication)
  http.post('/api/spots', async ({ request }) => {
    const authHeader = request.headers.get('authorization')
    
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return HttpResponse.json(
        { error: 'Authentication required' },
        { status: 401 }
      )
    }

    const body = await request.json() as any
    
    // Validation
    if (!body.name || !body.address || !body.latitude || !body.longitude) {
      return HttpResponse.json(
        { error: 'Missing required fields: name, address, latitude, longitude' },
        { status: 400 }
      )
    }

    const newSpot = {
      id: String(mockSpots.length + 1),
      ...body,
      soloFriendly: body.soloFriendlyRating > 3.0,
      averageRating: body.soloFriendlyRating || 0,
      reviewCount: 0,
      amenities: body.amenities || [],
      photos: body.photos || [],
      openingHours: body.openingHours || {},
    }

    mockSpots.push(newSpot)

    return HttpResponse.json(
      { data: newSpot },
      { status: 201 }
    )
  }),

  // Update spot
  http.put('/api/spots/:id', async ({ params, request }) => {
    const authHeader = request.headers.get('authorization')
    
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return HttpResponse.json(
        { error: 'Authentication required' },
        { status: 401 }
      )
    }

    const spotIndex = mockSpots.findIndex(s => s.id === params.id)
    
    if (spotIndex === -1) {
      return HttpResponse.json(
        { error: 'Spot not found' },
        { status: 404 }
      )
    }

    const body = await request.json() as any
    const updatedSpot = { ...mockSpots[spotIndex], ...body }
    
    mockSpots[spotIndex] = updatedSpot

    return HttpResponse.json({ data: updatedSpot })
  }),

  // Delete spot
  http.delete('/api/spots/:id', ({ params, request }) => {
    const authHeader = request.headers.get('authorization')
    
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return HttpResponse.json(
        { error: 'Authentication required' },
        { status: 401 }
      )
    }

    const spotIndex = mockSpots.findIndex(s => s.id === params.id)
    
    if (spotIndex === -1) {
      return HttpResponse.json(
        { error: 'Spot not found' },
        { status: 404 }
      )
    }

    mockSpots.splice(spotIndex, 1)

    return HttpResponse.json({ success: true })
  }),
]