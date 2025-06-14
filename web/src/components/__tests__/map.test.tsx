import React from 'react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import MapComponent from '../map'
import type { MapError } from '../map/types'

// Mock the custom hooks
const mockUseMaplibre = vi.fn()
const mockUsePmtiles = vi.fn()

vi.mock('../hooks/useMaplibre', () => ({
  useMaplibre: mockUseMaplibre,
}))

vi.mock('../hooks/usePmtiles', () => ({
  usePmtiles: mockUsePmtiles,
}))

// Mock map status components
vi.mock('../map/map-status', () => ({
  MapErrorDisplay: ({ error, className, height }: { error: MapError; className?: string; height?: string }) => (
    <div data-testid="map-error" className={className} style={{ height }} role="alert">
      Error: {error.message}
    </div>
  ),
  MapLoadingDisplay: ({ className }: { className?: string }) => (
    <div data-testid="map-loading" className={className} aria-live="polite">
      Loading map...
    </div>
  ),
}))

describe('MapComponent', () => {
  const mockContainerRef = { current: document.createElement('div') }

  beforeEach(() => {
    mockUseMaplibre.mockReturnValue({
      containerRef: mockContainerRef,
      mapState: 'loaded',
      error: null,
    })
    mockUsePmtiles.mockReturnValue(undefined)
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('Given the MapComponent is rendered', () => {
    it('When the component loads successfully, Then it should display the map container', () => {
      render(<MapComponent />)
      
      const mapContainer = screen.getByTestId('map-container')
      expect(mapContainer).toBeInTheDocument()
      expect(mapContainer).toHaveStyle({ height: '480px' })
    })

    it('When custom height is provided, Then it should use the custom height', () => {
      render(<MapComponent height="600px" />)
      
      const mapContainer = screen.getByRole('generic')
      expect(mapContainer).toHaveStyle({ height: '600px' })
    })

    it('When custom className is provided, Then it should apply the className', () => {
      render(<MapComponent className="custom-map-class" />)
      
      const mapContainer = screen.getByRole('generic')
      expect(mapContainer).toHaveClass('custom-map-class')
    })

    it('When component mounts, Then it should initialize PMTiles protocol', async () => {
      render(<MapComponent />)
      
      await waitFor(() => {
        expect(mockUsePmtiles).toHaveBeenCalled()
      })
    })

    it('When component mounts, Then it should initialize MapLibre with correct props', () => {
      const mockOnClick = vi.fn()
      const mockOnLoad = vi.fn()
      const mockOnError = vi.fn()
      const mockPoiFilter: ['==', string, string] = ['==', 'type', 'restaurant']

      render(
        <MapComponent
          onClick={mockOnClick}
          onLoad={mockOnLoad}
          onError={mockOnError}
          poiFilter={mockPoiFilter}
        />
      )
      
      expect(mockUseMaplibre).toHaveBeenCalledWith({
        onClick: mockOnClick,
        onLoad: mockOnLoad,
        onError: mockOnError,
        poiFilter: mockPoiFilter,
      })
    })
  })

  describe('Given the MapComponent is in loading state', () => {
    beforeEach(() => {
      mockUseMaplibre.mockReturnValue({
        containerRef: mockContainerRef,
        mapState: 'loading',
        error: null,
      })
    })

    it('When the map is loading, Then it should display the loading indicator', () => {
      render(<MapComponent />)
      
      expect(screen.getByTestId('map-loading')).toBeInTheDocument()
      expect(screen.getByText('Loading map...')).toBeInTheDocument()
    })

    it('When the map is loading, Then it should display both container and loading indicator', () => {
      render(<MapComponent />)
      
      expect(screen.getByRole('generic')).toBeInTheDocument()
      expect(screen.getByTestId('map-loading')).toBeInTheDocument()
    })
  })

  describe('Given the MapComponent is in error state', () => {
    const mockError: MapError = {
      type: 'loading',
      message: 'Failed to load map tiles',
    }

    beforeEach(() => {
      mockUseMaplibre.mockReturnValue({
        containerRef: mockContainerRef,
        mapState: 'error',
        error: mockError,
      })
    })

    it('When the map has an error, Then it should display the error message', () => {
      render(<MapComponent />)
      
      expect(screen.getByTestId('map-error')).toBeInTheDocument()
      expect(screen.getByText('Error: Failed to load map tiles')).toBeInTheDocument()
    })

    it('When the map has an error, Then it should not display the map container', () => {
      render(<MapComponent />)
      
      expect(screen.queryByRole('generic')).not.toBeInTheDocument()
      expect(screen.getByTestId('map-error')).toBeInTheDocument()
    })

    it('When the map has an error with custom height, Then error display should use custom height', () => {
      render(<MapComponent height="300px" />)
      
      const errorDisplay = screen.getByTestId('map-error')
      expect(errorDisplay).toHaveStyle({ height: '300px' })
    })

    it('When the map has an error with custom className, Then error display should use custom className', () => {
      render(<MapComponent className="error-map-class" />)
      
      const errorDisplay = screen.getByTestId('map-error')
      expect(errorDisplay).toHaveClass('error-map-class')
    })
  })

  describe('Given the MapComponent handles different error types', () => {
    it('When there is a configuration error, Then it should display the error', () => {
      const configError: MapError = {
        type: 'configuration',
        message: 'Invalid API key',
      }

      mockUseMaplibre.mockReturnValue({
        containerRef: mockContainerRef,
        mapState: 'error',
        error: configError,
      })

      render(<MapComponent />)
      
      expect(screen.getByText('Error: Invalid API key')).toBeInTheDocument()
    })

    it('When there is an initialization error, Then it should display the error', () => {
      const initError: MapError = {
        type: 'initialization',
        message: 'Failed to initialize map',
      }

      mockUseMaplibre.mockReturnValue({
        containerRef: mockContainerRef,
        mapState: 'error',
        error: initError,
      })

      render(<MapComponent />)
      
      expect(screen.getByText('Error: Failed to initialize map')).toBeInTheDocument()
    })
  })

  describe('Given the MapComponent accessibility features', () => {
    it('When the map is in error state, Then error display should have proper ARIA attributes', () => {
      const mockError: MapError = {
        type: 'loading',
        message: 'Network error',
      }

      mockUseMaplibre.mockReturnValue({
        containerRef: mockContainerRef,
        mapState: 'error',
        error: mockError,
      })

      render(<MapComponent />)
      
      expect(screen.getByRole('alert')).toBeInTheDocument()
    })

    it('When the map is loading, Then loading display should have proper ARIA attributes', () => {
      mockUseMaplibre.mockReturnValue({
        containerRef: mockContainerRef,
        mapState: 'loading',
        error: null,
      })

      render(<MapComponent />)
      
      const loadingDisplay = screen.getByTestId('map-loading')
      expect(loadingDisplay).toHaveAttribute('aria-live', 'polite')
    })
  })
})