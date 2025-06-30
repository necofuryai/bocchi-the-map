import { NextResponse } from 'next/server';
import { auth0 } from '@/lib/auth0';

export async function GET() {
  try {
    // Get the access token from Auth0
    const result = await auth0.getAccessToken();
    
    if (!result || !result.token) {
      return NextResponse.json(
        { 
          error: 'No access token available',
          message: 'User is not authenticated or session has expired'
        },
        { status: 401 }
      );
    }

    return NextResponse.json({ 
      accessToken: result.token,
      expiresAt: result.expiresAt || null
    });
  } catch (error) {
    console.error('Error getting access token:', error);
    
    // More specific error handling
    if (error instanceof Error) {
      if (error.message.includes('not authenticated')) {
        return NextResponse.json(
          { 
            error: 'Authentication required',
            message: 'Please log in to get an access token'
          },
          { status: 401 }
        );
      }
    }
    
    return NextResponse.json(
      { 
        error: 'Failed to get access token',
        message: 'An internal error occurred while retrieving the access token'
      },
      { status: 500 }
    );
  }
}