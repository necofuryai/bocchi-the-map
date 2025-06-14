import React from 'react'
import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MapErrorDisplay, MapLoadingDisplay } from '../map-status'
import type { MapError } from '../types'

describe('MapErrorDisplay Component', () => {
  describe('Given the MapErrorDisplay component', () => {
    it('When rendered with loading error, Then it should display the error message', () => {
      const error: MapError = {
        type: 'loading',
        message: 'Failed to load map tiles',
      }

      render(<MapErrorDisplay error={error} />)
      
      expect(screen.getByRole('alert')).toBeInTheDocument()
      expect(screen.getByText('Failed to load map tiles')).toBeInTheDocument()
      expect(screen.getByText('Failed to display map')).toBeInTheDocument()
    })

    it('When rendered with configuration error, Then it should display configuration message', () => {
      const error: MapError = {
        type: 'configuration',
        message: 'Invalid API key provided',
      }

      render(<MapErrorDisplay error={error} />)
      
      expect(screen.getByText('Invalid API key provided')).toBeInTheDocument()
      expect(screen.getByText('Please check your configuration')).toBeInTheDocument()
    })

    it('When rendered with initialization error, Then it should display the error message', () => {
      const error: MapError = {
        type: 'initialization',
        message: 'Map initialization failed',
      }

      render(<MapErrorDisplay error={error} />)
      
      expect(screen.getByText('Map initialization failed')).toBeInTheDocument()
      expect(screen.getByText('Failed to display map')).toBeInTheDocument()
    })

    it('When custom height is provided, Then it should apply the height style', () => {
      const error: MapError = {
        type: 'loading',
        message: 'Error message',
      }

      render(<MapErrorDisplay error={error} height="300px" />)
      
      const errorContainer = screen.getByRole('alert')
      expect(errorContainer).toHaveStyle({ height: '300px' })
    })

    it('When custom className is provided, Then it should apply the className', () => {
      const error: MapError = {
        type: 'loading',
        message: 'Error message',
      }

      render(<MapErrorDisplay error={error} className="custom-error-class" />)
      
      const errorContainer = screen.getByRole('alert')
      expect(errorContainer).toHaveClass('custom-error-class')
    })

    it('When rendered, Then it should have proper accessibility attributes', () => {
      const error: MapError = {
        type: 'loading',
        message: 'Accessibility test error',
      }

      render(<MapErrorDisplay error={error} />)
      
      const errorContainer = screen.getByRole('alert')
      expect(errorContainer).toBeInTheDocument()
      
      const errorIcon = screen.getByLabelText('Error')
      expect(errorIcon).toBeInTheDocument()
    })

    it('When rendered with default props, Then it should have default styling', () => {
      const error: MapError = {
        type: 'loading',
        message: 'Default styling test',
      }

      render(<MapErrorDisplay error={error} />)
      
      const errorContainer = screen.getByRole('alert')
      expect(errorContainer).toHaveClass(
        'w-full',
        'flex',
        'items-center',
        'justify-center',
        'bg-gray-100',
        'border',
        'border-gray-300',
        'rounded'
      )
    })
  })
})

describe('MapLoadingDisplay Component', () => {
  describe('Given the MapLoadingDisplay component', () => {
    it('When rendered, Then it should display loading message', () => {
      render(<MapLoadingDisplay />)
      
      expect(screen.getByText('Loading map...')).toBeInTheDocument()
    })

    it('When rendered, Then it should have proper accessibility attributes', () => {
      render(<MapLoadingDisplay />)
      
      const loadingContainer = screen.getByText('Loading map...').closest('div[aria-live="polite"]')
      expect(loadingContainer).toHaveAttribute('aria-live', 'polite')
    })

    it('When custom className is provided, Then it should apply the className', () => {
      render(<MapLoadingDisplay className="custom-loading-class" />)
      
      const loadingContainer = screen.getByText('Loading map...').closest('div[aria-live="polite"]')
      expect(loadingContainer).toHaveClass('custom-loading-class')
    })

    it('When rendered, Then it should have loading spinner', () => {
      render(<MapLoadingDisplay />)
      
      const loadingContainer = screen.getByText('Loading map...').closest('div[aria-live="polite"]')
      const spinner = loadingContainer?.querySelector('.animate-spin')
      expect(spinner).toBeInTheDocument()
      expect(spinner).toHaveClass('rounded-full', 'h-8', 'w-8', 'border-b-2', 'border-blue-600')
    })

    it('When rendered with default props, Then it should have default styling', () => {
      render(<MapLoadingDisplay />)
      
      const loadingContainer = screen.getByText('Loading map...').closest('div[aria-live="polite"]')
      expect(loadingContainer).toHaveClass(
        'absolute',
        'inset-0',
        'flex',
        'items-center',
        'justify-center',
        'bg-gray-50',
        'border',
        'border-gray-300',
        'rounded'
      )
    })

    it('When rendered, Then loading text should be properly styled', () => {
      render(<MapLoadingDisplay />)
      
      const loadingText = screen.getByText('Loading map...')
      expect(loadingText).toHaveClass('text-sm')
      expect(loadingText.parentElement).toHaveClass('text-center', 'text-gray-600')
    })
  })
})