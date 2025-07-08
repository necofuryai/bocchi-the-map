/**
 * SearchInput Component Implementation
 * 
 * This component was developed following TDD principles:
 * 1. Tests were written first (RED)
 * 2. Minimal implementation to pass tests (GREEN)
 * 3. Refactoring for clean code (REFACTOR)
 * 
 * The component supports the E2E BDD scenarios by providing
 * the search functionality required by users.
 */

'use client'

import { useState, useEffect, useRef, KeyboardEvent, ChangeEvent } from 'react'
import { MagnifyingGlassIcon, XMarkIcon, ArrowPathIcon } from '@heroicons/react/24/outline'

interface SearchInputProps {
  onSearch: (query: string) => void
  onClear: () => void
  placeholder?: string
  ariaLabel?: string
  defaultValue?: string
  loading?: boolean
  showSearchButton?: boolean
  recentSearches?: string[]
  maxLength?: number
  validateInput?: boolean
  className?: string
}

export function SearchInput({
  onSearch,
  onClear,
  placeholder = 'Search for spots...',
  ariaLabel = 'Search for spots',
  defaultValue = '',
  loading = false,
  showSearchButton = false,
  recentSearches = [],
  maxLength = 100,
  validateInput = false,
  className = '',
}: SearchInputProps) {
  const [value, setValue] = useState(defaultValue)
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)
  const suggestionsRef = useRef<HTMLDivElement>(null)

  // Set initial value
  useEffect(() => {
    setValue(defaultValue)
  }, [defaultValue])

  // Handle input change
  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value
    
    // Enforce max length
    if (maxLength && newValue.length > maxLength) {
      setError('Search query too long')
      return
    }
    
    // Validate input if required
    if (validateInput && newValue) {
      const invalidChars = /[<>"']/
      if (invalidChars.test(newValue)) {
        setError('Invalid characters detected')
        return
      }
    }
    
    // Clear error if input is valid
    setError(null)
    setValue(newValue)
  }

  // Handle Enter key press
  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleSearch()
    } else if (e.key === 'Escape') {
      setShowSuggestions(false)
    }
  }

  // Handle search execution
  const handleSearch = () => {
    const trimmedValue = value.trim()
    if (trimmedValue && !error && !loading) {
      onSearch(trimmedValue)
      setShowSuggestions(false)
    }
  }

  // Handle clear
  const handleClear = () => {
    setValue('')
    setError(null)
    setShowSuggestions(false)
    onClear()
    inputRef.current?.focus()
  }

  // Handle suggestion click
  const handleSuggestionClick = (suggestion: string) => {
    setValue(suggestion)
    setShowSuggestions(false)
    onSearch(suggestion)
  }

  // Handle input focus
  const handleFocus = () => {
    if (recentSearches.length > 0) {
      setShowSuggestions(true)
    }
  }

  // Handle click outside to close suggestions
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        suggestionsRef.current &&
        !suggestionsRef.current.contains(event.target as Node) &&
        inputRef.current &&
        !inputRef.current.contains(event.target as Node)
      ) {
        setShowSuggestions(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const hasValue = value.length > 0
  const isSearchDisabled = !hasValue || loading || !!error

  return (
    <div className={`relative w-full ${className}`}>
      {/* Search Input Container */}
      <div className="relative flex items-center">
        {/* Search Input */}
        <input
          ref={inputRef}
          data-testid="search-input"
          type="text"
          value={value}
          onChange={handleInputChange}
          onKeyDown={handleKeyDown}
          onFocus={handleFocus}
          placeholder={placeholder}
          aria-label={ariaLabel}
          role="searchbox"
          disabled={loading}
          className={`
            w-full px-4 py-3 pr-${hasValue ? '20' : '12'} 
            border border-gray-300 rounded-lg
            focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent
            disabled:bg-gray-100 disabled:cursor-not-allowed
            ${error ? 'border-red-500 focus:ring-red-500' : ''}
            ${showSearchButton ? 'pr-24' : ''}
          `}
        />

        {/* Loading Indicator */}
        {loading && (
          <div 
            data-testid="search-loading"
            className="absolute right-3 flex items-center"
          >
            <ArrowPathIcon className="w-5 h-5 animate-spin text-gray-400" />
          </div>
        )}

        {/* Clear Button */}
        {hasValue && !loading && (
          <button
            data-testid="clear-button"
            type="button"
            onClick={handleClear}
            className="absolute right-3 p-1 hover:bg-gray-100 rounded-full transition-colors"
            aria-label="Clear search"
          >
            <XMarkIcon className="w-4 h-4 text-gray-400" />
          </button>
        )}

        {/* Search Button */}
        {showSearchButton && (
          <button
            data-testid="search-button"
            type="button"
            onClick={handleSearch}
            disabled={isSearchDisabled}
            aria-label={loading ? 'Searching...' : 'Search'}
            className={`
              absolute right-2 px-3 py-1.5 
              bg-blue-600 text-white rounded-md
              hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed
              transition-colors flex items-center gap-2
            `}
          >
            {loading ? (
              <ArrowPathIcon className="w-4 h-4 animate-spin" />
            ) : (
              <MagnifyingGlassIcon className="w-4 h-4" />
            )}
            <span className="hidden sm:inline">Search</span>
          </button>
        )}
      </div>

      {/* Error Message */}
      {error && (
        <div 
          data-testid="search-error"
          className="mt-2 text-sm text-red-600"
          role="alert"
        >
          {error}
        </div>
      )}

      {/* Search Suggestions */}
      {showSuggestions && recentSearches.length > 0 && (
        <div
          ref={suggestionsRef}
          data-testid="search-suggestions"
          className="absolute z-10 w-full mt-1 bg-white border border-gray-200 rounded-lg shadow-lg max-h-48 overflow-y-auto"
        >
          <div className="px-3 py-2 text-xs font-medium text-gray-500 border-b">
            Recent Searches
          </div>
          {recentSearches.map((search, index) => (
            <button
              key={index}
              type="button"
              onClick={() => handleSuggestionClick(search)}
              className="w-full px-3 py-2 text-left hover:bg-gray-50 focus:bg-gray-50 focus:outline-none transition-colors"
            >
              <div className="flex items-center gap-2">
                <MagnifyingGlassIcon className="w-4 h-4 text-gray-400" />
                <span className="text-sm text-gray-700">{search}</span>
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  )
}

/**
 * Development Notes:
 * 
 * This component was built using TDD methodology:
 * 
 * 1. RED Phase:
 *    - Tests were written first defining expected behavior
 *    - All tests initially failed (RED)
 * 
 * 2. GREEN Phase:
 *    - Minimal implementation to make tests pass
 *    - Focus on functionality over perfect code
 * 
 * 3. REFACTOR Phase:
 *    - Clean up code while keeping tests green
 *    - Improve performance, readability, and maintainability
 * 
 * Key Features Implemented:
 * - Basic search input with validation
 * - Clear functionality
 * - Loading states
 * - Recent search suggestions
 * - Keyboard navigation (Enter, Escape)
 * - Accessibility support (ARIA labels, roles)
 * - Error handling and validation
 * - Mobile-responsive design
 * 
 * Next Steps:
 * - Implement SearchResults component
 * - Create SpotItem component
 * - Build FilterPanel component
 * - Integrate with search API
 */