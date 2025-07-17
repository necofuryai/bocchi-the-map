// Domain entities - MVP simplified version
export interface Review {
  id: string
  spotId: string
  userId: string
  userName: string
  userAvatar?: string
  rating: number // 1-5 stars only
  soloFriendlyRating: number // 1-5 stars for solo-friendly rating
  comment?: string
  tags?: string[]
  photos?: string[]
  helpful: number // helpful votes count
  notHelpful: number // not helpful votes count
  createdAt: string
  updatedAt: string
}

// MVP simplified user (basic info only)
export interface DomainUser {
  id: string
  email: string
  name?: string
  createdAt: string
  updatedAt: string
}

// Spot entity - MVP simplified version
export interface Spot {
  id: string
  name: string
  category: string
  address: string
  countryCode: string
  createdAt: string
  updatedAt: string
}

// Common response types
export interface PaginationResponse<T = unknown> {
  items: T[]
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
}

// PMTiles and Map-related types
export interface PMTilesMetadata {
  vector_layers?: Array<{
    id: string;
    fields?: Record<string, string>;
    minzoom?: number;
    maxzoom?: number;
    description?: string;
  }>;
  name?: string;
  description?: string;
  version?: string;
  [key: string]: unknown;
}

export interface PMTilesHeader {
  tileType: number;
  minZoom: number;
  maxZoom: number;
  centerLon: number;
  centerLat: number;
  minLon: number;
  minLat: number;
  maxLon: number;
  maxLat: number;
  clustered?: boolean;
  [key: string]: unknown;
}

export interface PMTilesTileData {
  data: ArrayBuffer;
  [key: string]: unknown;
}

// POI data structure from PMTiles
export interface POIProperties {
  name?: string;
  class?: string;
  subclass?: string;
  kind?: string;
  kind_detail?: string;
  [key: string]: unknown;
}

// Map filter and expression types
export type MapFilter = [
  string,
  ...(string | number | boolean | MapFilter | Array<string | number>)[]
];

export type MapExpression = [
  string,
  ...(string | number | boolean | MapExpression | Array<string | number>)[]
];

// Note: Map-related types removed for MVP simplification
// Geographic search and POI features have been removed

// Clerk Auth related types
export interface ClerkUser {
  id: string
  emailAddresses: Array<{
    emailAddress: string
    id: string
  }>
  firstName?: string
  lastName?: string
  imageUrl?: string
  createdAt?: number
  updatedAt?: number
}

export interface AuthUser {
  id: string
  email: string
  name?: string
}

// Auth page component types
export interface AuthPageProps {
  redirectTo?: string
}

export interface LoginPageState {
  isLoading: boolean
  error?: string
}

