import '@testing-library/jest-dom'
import { beforeAll, beforeEach, afterEach, afterAll, vi } from 'vitest'
import { type ReactNode } from 'react'
import { cleanup } from '@testing-library/react'
import { server } from '@/mocks/server'

// Setup MSW for API mocking
beforeAll(() => {
  // Start the server before all tests
  server.listen({ onUnhandledRequest: 'error' })
})

afterEach(() => {
  // Reset any request handlers that we may add during the tests
  server.resetHandlers()
  // Cleanup DOM after each test
  cleanup()
})

afterAll(() => {
  // Clean up after the tests are finished
  server.close()
})

// Global test setup
beforeEach(() => {
  // Reset any mocks or state before each test
  vi.clearAllMocks()
  vi.resetModules()
  
  // Reset any global state that might affect tests
  if (typeof window !== 'undefined') {
    // Clear localStorage
    window.localStorage.clear()
    // Clear sessionStorage
    window.sessionStorage.clear()
  }
})

// Mock Next.js router
vi.mock('next/navigation', () => ({
  useRouter: vi.fn(() => ({
    push: vi.fn(),
    back: vi.fn(),
    forward: vi.fn(),
    refresh: vi.fn(),
    replace: vi.fn(),
    prefetch: vi.fn(),
  })),
  usePathname: vi.fn(() => '/'),
  useSearchParams: vi.fn(() => new URLSearchParams()),
}))


// Mock MapLibre GL JS
vi.mock('maplibre-gl', () => ({
  Map: vi.fn(() => ({
    on: vi.fn(),
    off: vi.fn(),
    remove: vi.fn(),
    addControl: vi.fn(),
    removeControl: vi.fn(),
    addSource: vi.fn(),
    removeSource: vi.fn(),
    addLayer: vi.fn(),
    removeLayer: vi.fn(),
    setStyle: vi.fn(),
    getStyle: vi.fn(),
    flyTo: vi.fn(),
    panTo: vi.fn(),
    zoomTo: vi.fn(),
    getBounds: vi.fn(),
    getCenter: vi.fn(),
    getZoom: vi.fn(),
    loaded: vi.fn(() => true),
    resize: vi.fn(),
  })),
  NavigationControl: vi.fn(() => ({
    onAdd: vi.fn(),
    onRemove: vi.fn(),
  })),
  Marker: vi.fn(() => ({
    setLngLat: vi.fn().mockReturnThis(),
    addTo: vi.fn().mockReturnThis(),
    remove: vi.fn(),
    setPopup: vi.fn().mockReturnThis(),
  })),
  Popup: vi.fn(() => ({
    setLngLat: vi.fn().mockReturnThis(),
    setHTML: vi.fn().mockReturnThis(),
    addTo: vi.fn().mockReturnThis(),
    remove: vi.fn(),
  })),
}))

// Mock next-themes
vi.mock('next-themes', () => ({
  useTheme: vi.fn(() => ({
    theme: 'light',
    setTheme: vi.fn(),
  })),
  ThemeProvider: ({ children }: { children: ReactNode }) => children,
}))