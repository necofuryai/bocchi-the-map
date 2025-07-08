import { http, HttpResponse } from 'msw'

// Mock user data
const mockUsers = [
  {
    id: 'user-1',
    email: 'test@example.com',
    name: 'Test User',
    avatar: '/images/avatar1.jpg',
    preferences: {
      theme: 'light',
      notifications: true,
      language: 'ja',
    },
    createdAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 'user-2',
    email: 'solo@traveler.com',
    name: 'Solo Traveler',
    avatar: '/images/avatar2.jpg',
    preferences: {
      theme: 'dark',
      notifications: false,
      language: 'en',
    },
    createdAt: '2024-02-01T00:00:00Z',
  },
]

// Mock sessions/tokens
const mockSessions = new Map<string, { userId: string; expires: Date }>()

export const authHandlers = [
  // Login/Authentication
  http.post('/api/auth/login', async ({ request }) => {
    const body = await request.json() as any
    
    if (!body.email || !body.password) {
      return HttpResponse.json(
        { error: 'Email and password are required' },
        { status: 400 }
      )
    }

    const user = mockUsers.find(u => u.email === body.email)
    
    if (!user || body.password !== 'password123') {
      return HttpResponse.json(
        { error: 'Invalid credentials' },
        { status: 401 }
      )
    }

    // Generate mock token
    const token = `mock-token-${user.id}-${Date.now()}`
    const expires = new Date(Date.now() + 24 * 60 * 60 * 1000) // 24 hours
    
    mockSessions.set(token, { userId: user.id, expires })

    return HttpResponse.json({
      token,
      user: {
        id: user.id,
        email: user.email,
        name: user.name,
        avatar: user.avatar,
      },
      expires: expires.toISOString(),
    })
  }),

  // Logout
  http.post('/api/auth/logout', ({ request }) => {
    const authHeader = request.headers.get('authorization')
    
    if (authHeader && authHeader.startsWith('Bearer ')) {
      const token = authHeader.substring(7)
      mockSessions.delete(token)
    }

    return HttpResponse.json({ success: true })
  }),

  // Get current user
  http.get('/api/auth/me', ({ request }) => {
    const authHeader = request.headers.get('authorization')
    
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return HttpResponse.json(
        { error: 'Authentication required' },
        { status: 401 }
      )
    }

    const token = authHeader.substring(7)
    const session = mockSessions.get(token)
    
    if (!session || session.expires < new Date()) {
      return HttpResponse.json(
        { error: 'Invalid or expired token' },
        { status: 401 }
      )
    }

    const user = mockUsers.find(u => u.id === session.userId)
    
    if (!user) {
      return HttpResponse.json(
        { error: 'User not found' },
        { status: 404 }
      )
    }

    return HttpResponse.json({
      user: {
        id: user.id,
        email: user.email,
        name: user.name,
        avatar: user.avatar,
        preferences: user.preferences,
        createdAt: user.createdAt,
      },
    })
  }),

  // Update user profile
  http.put('/api/auth/profile', async ({ request }) => {
    const authHeader = request.headers.get('authorization')
    
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return HttpResponse.json(
        { error: 'Authentication required' },
        { status: 401 }
      )
    }

    const token = authHeader.substring(7)
    const session = mockSessions.get(token)
    
    if (!session || session.expires < new Date()) {
      return HttpResponse.json(
        { error: 'Invalid or expired token' },
        { status: 401 }
      )
    }

    const userIndex = mockUsers.findIndex(u => u.id === session.userId)
    
    if (userIndex === -1) {
      return HttpResponse.json(
        { error: 'User not found' },
        { status: 404 }
      )
    }

    const body = await request.json() as any
    const updatedUser = { ...mockUsers[userIndex], ...body }
    
    mockUsers[userIndex] = updatedUser

    return HttpResponse.json({
      user: {
        id: updatedUser.id,
        email: updatedUser.email,
        name: updatedUser.name,
        avatar: updatedUser.avatar,
        preferences: updatedUser.preferences,
        createdAt: updatedUser.createdAt,
      },
    })
  }),

  // Register new user
  http.post('/api/auth/register', async ({ request }) => {
    const body = await request.json() as any
    
    if (!body.email || !body.password || !body.name) {
      return HttpResponse.json(
        { error: 'Email, password, and name are required' },
        { status: 400 }
      )
    }

    // Check if user already exists
    const existingUser = mockUsers.find(u => u.email === body.email)
    if (existingUser) {
      return HttpResponse.json(
        { error: 'User with this email already exists' },
        { status: 409 }
      )
    }

    // Create new user
    const newUser = {
      id: `user-${mockUsers.length + 1}`,
      email: body.email,
      name: body.name,
      avatar: body.avatar || '/images/default-avatar.jpg',
      preferences: {
        theme: 'light',
        notifications: true,
        language: 'ja',
      },
      createdAt: new Date().toISOString(),
    }

    mockUsers.push(newUser)

    // Generate token
    const token = `mock-token-${newUser.id}-${Date.now()}`
    const expires = new Date(Date.now() + 24 * 60 * 60 * 1000)
    
    mockSessions.set(token, { userId: newUser.id, expires })

    return HttpResponse.json(
      {
        token,
        user: {
          id: newUser.id,
          email: newUser.email,
          name: newUser.name,
          avatar: newUser.avatar,
        },
        expires: expires.toISOString(),
      },
      { status: 201 }
    )
  }),

  // Password reset request
  http.post('/api/auth/password-reset-request', async ({ request }) => {
    const body = await request.json() as any
    
    if (!body.email) {
      return HttpResponse.json(
        { error: 'Email is required' },
        { status: 400 }
      )
    }

    const user = mockUsers.find(u => u.email === body.email)
    
    if (!user) {
      // For security, don't reveal whether email exists
      return HttpResponse.json({ success: true })
    }

    // In real implementation, send email with reset token
    return HttpResponse.json({ success: true })
  }),

  // Refresh token
  http.post('/api/auth/refresh', ({ request }) => {
    const authHeader = request.headers.get('authorization')
    
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return HttpResponse.json(
        { error: 'Authentication required' },
        { status: 401 }
      )
    }

    const token = authHeader.substring(7)
    const session = mockSessions.get(token)
    
    if (!session) {
      return HttpResponse.json(
        { error: 'Invalid token' },
        { status: 401 }
      )
    }

    // Generate new token
    const newToken = `mock-token-${session.userId}-${Date.now()}`
    const expires = new Date(Date.now() + 24 * 60 * 60 * 1000)
    
    mockSessions.delete(token)
    mockSessions.set(newToken, { userId: session.userId, expires })

    const user = mockUsers.find(u => u.id === session.userId)

    return HttpResponse.json({
      token: newToken,
      user: user ? {
        id: user.id,
        email: user.email,
        name: user.name,
        avatar: user.avatar,
      } : null,
      expires: expires.toISOString(),
    })
  }),
]