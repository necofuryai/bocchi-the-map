/**
 * TDD Unit Test Example: SearchInput Component
 * 
 * This demonstrates the Inner Loop of TDD+BDD:
 * RED -> GREEN -> REFACTOR
 * 
 * Following the E2E test that drives this implementation,
 * we implement the SearchInput component using TDD.
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { SearchInput } from '../search-input'

describe('SearchInput Component', () => {
  const mockOnSearch = vi.fn()
  const mockOnClear = vi.fn()
  
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Given the SearchInput component is rendered', () => {
    describe('When the component loads', () => {
      it('Then it should display the search input with correct attributes', () => {
        // Given
        render(
          <SearchInput 
            onSearch={mockOnSearch} 
            onClear={mockOnClear}
            placeholder="Search for spots..."
          />
        )
        
        // Then
        const searchInput = screen.getByTestId('search-input')
        expect(searchInput).toBeInTheDocument()
        expect(searchInput).toHaveAttribute('type', 'text')
        expect(searchInput).toHaveAttribute('placeholder', 'Search for spots...')
        expect(searchInput).toHaveValue('')
      })

      it('Then it should have proper accessibility attributes', () => {
        // Given
        render(
          <SearchInput 
            onSearch={mockOnSearch} 
            onClear={mockOnClear}
            ariaLabel="Search for solo-friendly spots"
          />
        )
        
        // Then
        const searchInput = screen.getByTestId('search-input')
        expect(searchInput).toHaveAttribute('aria-label', 'Search for solo-friendly spots')
        expect(searchInput).toHaveAttribute('role', 'searchbox')
      })
    })

    describe('When user types in the search input', () => {
      it('Then the input value should update correctly', async () => {
        // Given
        const user = userEvent.setup()
        render(<SearchInput onSearch={mockOnSearch} onClear={mockOnClear} />)
        const searchInput = screen.getByTestId('search-input')
        
        // When
        await user.type(searchInput, 'quiet cafe')
        
        // Then
        expect(searchInput).toHaveValue('quiet cafe')
      })

      it('Then it should show a clear button when there is text', async () => {
        // Given
        const user = userEvent.setup()
        render(<SearchInput onSearch={mockOnSearch} onClear={mockOnClear} />)
        const searchInput = screen.getByTestId('search-input')
        
        // When
        await user.type(searchInput, 'coffee')
        
        // Then
        const clearButton = screen.getByTestId('clear-button')
        expect(clearButton).toBeInTheDocument()
        expect(clearButton).toBeVisible()
      })

      it('Then it should not show clear button when input is empty', () => {
        // Given
        render(<SearchInput onSearch={mockOnSearch} onClear={mockOnClear} />)
        
        // Then
        const clearButton = screen.queryByTestId('clear-button')
        expect(clearButton).not.toBeInTheDocument()
      })
    })

    describe('When user presses Enter key', () => {
      it('Then onSearch should be called with the input value', async () => {
        // Given
        const user = userEvent.setup()
        render(<SearchInput onSearch={mockOnSearch} onClear={mockOnClear} />)
        const searchInput = screen.getByTestId('search-input')
        
        // When
        await user.type(searchInput, 'solo-friendly restaurant')
        await user.keyboard('{Enter}')
        
        // Then
        expect(mockOnSearch).toHaveBeenCalledWith('solo-friendly restaurant')
        expect(mockOnSearch).toHaveBeenCalledTimes(1)
      })

      it('Then onSearch should not be called if input is empty', async () => {
        // Given
        const user = userEvent.setup()
        render(<SearchInput onSearch={mockOnSearch} onClear={mockOnClear} />)
        const searchInput = screen.getByTestId('search-input')
        
        // When
        await user.click(searchInput)
        await user.keyboard('{Enter}')
        
        // Then
        expect(mockOnSearch).not.toHaveBeenCalled()
      })

      it('Then onSearch should not be called if input contains only whitespace', async () => {
        // Given
        const user = userEvent.setup()
        render(<SearchInput onSearch={mockOnSearch} onClear={mockOnClear} />)
        const searchInput = screen.getByTestId('search-input')
        
        // When
        await user.type(searchInput, '   ')
        await user.keyboard('{Enter}')
        
        // Then
        expect(mockOnSearch).not.toHaveBeenCalled()
      })
    })

    describe('When user clicks the search button', () => {
      it('Then onSearch should be called with the input value', async () => {
        // Given
        const user = userEvent.setup()
        render(<SearchInput onSearch={mockOnSearch} onClear={mockOnClear} showSearchButton />)
        const searchInput = screen.getByTestId('search-input')
        const searchButton = screen.getByTestId('search-button')
        
        // When
        await user.type(searchInput, 'library')
        await user.click(searchButton)
        
        // Then
        expect(mockOnSearch).toHaveBeenCalledWith('library')
        expect(mockOnSearch).toHaveBeenCalledTimes(1)
      })

      it('Then search button should be disabled when input is empty', () => {
        // Given
        render(<SearchInput onSearch={mockOnSearch} onClear={mockOnClear} showSearchButton />)
        
        // Then
        const searchButton = screen.getByTestId('search-button')
        expect(searchButton).toBeDisabled()
      })

      it('Then search button should be enabled when input has value', async () => {
        // Given
        const user = userEvent.setup()
        render(<SearchInput onSearch={mockOnSearch} onClear={mockOnClear} showSearchButton />)
        const searchInput = screen.getByTestId('search-input')
        
        // When
        await user.type(searchInput, 'cafe')
        
        // Then
        const searchButton = screen.getByTestId('search-button')
        expect(searchButton).not.toBeDisabled()
      })
    })

    describe('When user clicks the clear button', () => {
      it('Then the input should be cleared and onClear should be called', async () => {
        // Given
        const user = userEvent.setup()
        render(<SearchInput onSearch={mockOnSearch} onClear={mockOnClear} />)
        const searchInput = screen.getByTestId('search-input')
        
        // When
        await user.type(searchInput, 'some text')
        const clearButton = screen.getByTestId('clear-button')
        await user.click(clearButton)
        
        // Then
        expect(searchInput).toHaveValue('')
        expect(mockOnClear).toHaveBeenCalledTimes(1)
        
        // And clear button should disappear
        expect(screen.queryByTestId('clear-button')).not.toBeInTheDocument()
      })
    })

    describe('When component receives a default value', () => {
      it('Then the input should display the default value', () => {
        // Given
        const defaultValue = 'quiet spaces'
        render(
          <SearchInput 
            onSearch={mockOnSearch} 
            onClear={mockOnClear} 
            defaultValue={defaultValue}
          />
        )
        
        // Then
        const searchInput = screen.getByTestId('search-input')
        expect(searchInput).toHaveValue(defaultValue)
        
        // And clear button should be visible
        expect(screen.getByTestId('clear-button')).toBeInTheDocument()
      })
    })

    describe('When component is in loading state', () => {
      it('Then it should show loading indicator and disable input', () => {
        // Given
        render(
          <SearchInput 
            onSearch={mockOnSearch} 
            onClear={mockOnClear}
            loading={true}
          />
        )
        
        // Then
        const searchInput = screen.getByTestId('search-input')
        const loadingIndicator = screen.getByTestId('search-loading')
        
        expect(searchInput).toBeDisabled()
        expect(loadingIndicator).toBeInTheDocument()
      })

      it('Then search button should show loading state', () => {
        // Given
        render(
          <SearchInput 
            onSearch={mockOnSearch} 
            onClear={mockOnClear}
            showSearchButton
            loading={true}
          />
        )
        
        // Then
        const searchButton = screen.getByTestId('search-button')
        expect(searchButton).toBeDisabled()
        expect(searchButton).toHaveAttribute('aria-label', 'Searching...')
      })
    })

    describe('When component receives recent searches', () => {
      it('Then it should show search suggestions dropdown', async () => {
        // Given
        const recentSearches = ['coffee shops', 'quiet libraries', 'co-working spaces']
        const user = userEvent.setup()
        
        render(
          <SearchInput 
            onSearch={mockOnSearch} 
            onClear={mockOnClear}
            recentSearches={recentSearches}
          />
        )
        
        // When
        const searchInput = screen.getByTestId('search-input')
        await user.click(searchInput)
        
        // Then
        const suggestions = screen.getByTestId('search-suggestions')
        expect(suggestions).toBeInTheDocument()
        
        recentSearches.forEach(search => {
          expect(screen.getByText(search)).toBeInTheDocument()
        })
      })

      it('Then clicking a suggestion should trigger search', async () => {
        // Given
        const recentSearches = ['coffee shops', 'quiet libraries']
        const user = userEvent.setup()
        
        render(
          <SearchInput 
            onSearch={mockOnSearch} 
            onClear={mockOnClear}
            recentSearches={recentSearches}
          />
        )
        
        // When
        const searchInput = screen.getByTestId('search-input')
        await user.click(searchInput)
        
        const suggestion = screen.getByText('coffee shops')
        await user.click(suggestion)
        
        // Then
        expect(mockOnSearch).toHaveBeenCalledWith('coffee shops')
        expect(searchInput).toHaveValue('coffee shops')
      })
    })
  })

  describe('Given the SearchInput has validation requirements', () => {
    describe('When input exceeds maximum length', () => {
      it('Then it should prevent further typing and show error', async () => {
        // Given
        const user = userEvent.setup()
        const maxLength = 50
        const longText = 'a'.repeat(maxLength + 10)
        
        render(
          <SearchInput 
            onSearch={mockOnSearch} 
            onClear={mockOnClear}
            maxLength={maxLength}
          />
        )
        
        // When
        const searchInput = screen.getByTestId('search-input')
        await user.type(searchInput, longText)
        
        // Then
        expect(searchInput).toHaveValue(longText.substring(0, maxLength))
        expect(screen.getByTestId('search-error')).toBeInTheDocument()
        expect(screen.getByTestId('search-error')).toHaveTextContent('Search query too long')
      })
    })

    describe('When input contains invalid characters', () => {
      it('Then it should show validation error', async () => {
        // Given
        const user = userEvent.setup()
        const invalidChars = ['<', '>', '"', "'"]
        
        render(
          <SearchInput 
            onSearch={mockOnSearch} 
            onClear={mockOnClear}
            validateInput={true}
          />
        )
        
        // When
        const searchInput = screen.getByTestId('search-input')
        await user.type(searchInput, 'coffee<script>')
        
        // Then
        expect(screen.getByTestId('search-error')).toBeInTheDocument()
        expect(screen.getByTestId('search-error')).toHaveTextContent('Invalid characters detected')
      })
    })
  })
})