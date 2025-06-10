// ドメインエンティティ
export interface Spot {
  id: string
  name: string
  nameI18n?: Record<string, string>
  latitude: number
  longitude: number
  category: string
  address: string
  addressI18n?: Record<string, string>
  countryCode: string
  averageRating: number
  reviewCount: number
  createdAt: string
  updatedAt: string
}

export interface Review {
  id: string
  spotId: string
  userId: string
  rating: number
  comment?: string
  ratingAspects?: Record<string, number>
  createdAt: string
  updatedAt: string
}

export interface User {
  id: string
  email: string
  displayName: string
  avatarUrl?: string
  preferences: UserPreferences
}

export interface UserPreferences {
  language: 'ja' | 'en'
  darkMode: boolean
  timezone: string
}

// 共通レスポンス型
export interface PaginationResponse {
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
}

// 地図関連の型
export interface MapPosition {
  latitude: number
  longitude: number
  zoom?: number
}

export interface POIProperties {
  name?: string
  kind?: string
  script?: string
  min_zoom?: number
}