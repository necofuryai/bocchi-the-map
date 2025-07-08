/**
 * BDD E2E Test Example: Spot Search Feature
 * 
 * This file demonstrates the Outside-In TDD+BDD approach:
 * 1. Start with user behavior (BDD E2E test) - RED
 * 2. Implement components with TDD - GREEN
 * 3. Integration and refactoring - REFACTOR
 * 
 * Feature: Solo-friendly spot search
 * As a solo traveler
 * I want to search for spots suitable for solo activities
 * So that I can find comfortable places to visit alone
 */

import { test, expect } from '@playwright/test'

test.describe('Spot Search Feature', () => {
  test.describe('Given I am on the search page', () => {
    test.beforeEach(async ({ page }) => {
      // Start the application
      await page.goto('/search')
    })

    test('When I search for "quiet cafe", Then I should see relevant solo-friendly results', async ({ page }) => {
      // Given - Search page is loaded
      await expect(page.getByTestId('search-page')).toBeVisible()
      
      // When - I enter a search query
      const searchInput = page.getByTestId('search-input')
      await searchInput.fill('quiet cafe')
      await page.keyboard.press('Enter')
      
      // Then - Search results should be displayed
      await expect(page.getByTestId('search-results')).toBeVisible()
      
      // And - Results should contain relevant spots
      const spotItems = page.getByTestId('spot-item')
      await expect(spotItems).toHaveCount.greaterThan(0)
      
      // And - Each result should show spot information
      const firstSpot = spotItems.first()
      await expect(firstSpot.getByTestId('spot-name')).toBeVisible()
      await expect(firstSpot.getByTestId('spot-rating')).toBeVisible()
      await expect(firstSpot.getByTestId('solo-friendly-indicator')).toBeVisible()
    })

    test('When I apply solo-friendly filter, Then only solo-friendly spots should be shown', async ({ page }) => {
      // Given - Search page is loaded and I search for cafes
      await page.getByTestId('search-input').fill('cafe')
      await page.keyboard.press('Enter')
      await expect(page.getByTestId('search-results')).toBeVisible()
      
      // When - I apply the solo-friendly filter
      const filterButton = page.getByTestId('filter-toggle')
      await filterButton.click()
      
      const soloFriendlyFilter = page.getByTestId('solo-friendly-filter')
      await soloFriendlyFilter.check()
      
      const applyFiltersButton = page.getByTestId('apply-filters')
      await applyFiltersButton.click()
      
      // Then - Only solo-friendly spots should be displayed
      const spotItems = page.getByTestId('spot-item')
      const count = await spotItems.count()
      
      for (let i = 0; i < count; i++) {
        const spot = spotItems.nth(i)
        const soloFriendlyBadge = spot.getByTestId('solo-friendly-badge')
        await expect(soloFriendlyBadge).toBeVisible()
      }
      
      // And - Filter indicator should show active filter
      await expect(page.getByTestId('active-filters')).toContainText('Solo-friendly')
    })

    test('When I click on a search result, Then I should navigate to the spot details', async ({ page }) => {
      // Given - Search results are displayed
      await page.getByTestId('search-input').fill('coffee')
      await page.keyboard.press('Enter')
      await expect(page.getByTestId('search-results')).toBeVisible()
      
      // When - I click on the first search result
      const firstSpot = page.getByTestId('spot-item').first()
      await firstSpot.click()
      
      // Then - I should be navigated to the spot details page
      await expect(page).toHaveURL(/\/spots\/.*/)
      await expect(page.getByTestId('spot-details')).toBeVisible()
    })

    test('When search returns no results, Then I should see an appropriate message', async ({ page }) => {
      // Given - Search page is loaded
      
      // When - I search for something that doesn't exist
      await page.getByTestId('search-input').fill('nonexistent place xyz123')
      await page.keyboard.press('Enter')
      
      // Then - No results message should be displayed
      await expect(page.getByTestId('no-results-message')).toBeVisible()
      await expect(page.getByTestId('no-results-message')).toContainText('No spots found')
      
      // And - Suggestions should be provided
      await expect(page.getByTestId('search-suggestions')).toBeVisible()
    })

    test('When search fails, Then I should see an error message with retry option', async ({ page }) => {
      // Given - Mock API to return error
      await page.route('/api/spots/search*', route => route.abort())
      
      // When - I perform a search
      await page.getByTestId('search-input').fill('coffee')
      await page.keyboard.press('Enter')
      
      // Then - Error message should be displayed
      await expect(page.getByTestId('error-message')).toBeVisible()
      await expect(page.getByTestId('error-message')).toContainText('Failed to search')
      
      // And - Retry button should be available
      const retryButton = page.getByTestId('retry-search')
      await expect(retryButton).toBeVisible()
      
      // When - I click retry
      await page.unroute('/api/spots/search*') // Remove the mock
      await retryButton.click()
      
      // Then - Search should work normally
      await expect(page.getByTestId('search-results')).toBeVisible()
    })
  })

  test.describe('Given I am not authenticated', () => {
    test.beforeEach(async ({ page }) => {
      // Ensure user is not authenticated
      await page.goto('/search')
    })

    test('When I try to save a spot, Then I should be prompted to login', async ({ page }) => {
      // Given - Search results are displayed
      await page.getByTestId('search-input').fill('cafe')
      await page.keyboard.press('Enter')
      await expect(page.getByTestId('search-results')).toBeVisible()
      
      // When - I try to save a spot
      const firstSpot = page.getByTestId('spot-item').first()
      const saveButton = firstSpot.getByTestId('save-spot-button')
      await saveButton.click()
      
      // Then - Login prompt should appear
      await expect(page.getByTestId('login-prompt')).toBeVisible()
      await expect(page.getByTestId('login-prompt')).toContainText('Please log in to save spots')
      
      // And - Login button should be available
      await expect(page.getByTestId('login-button')).toBeVisible()
    })
  })

  test.describe('Given I am on a mobile device', () => {
    test.beforeEach(async ({ page }) => {
      // Set mobile viewport
      await page.setViewportSize({ width: 375, height: 667 })
      await page.goto('/search')
    })

    test('When I use the search interface, Then it should be mobile-optimized', async ({ page }) => {
      // Given - Mobile search page is loaded
      await expect(page.getByTestId('search-page')).toBeVisible()
      
      // Then - Search input should be properly sized
      const searchInput = page.getByTestId('search-input')
      await expect(searchInput).toBeVisible()
      
      // And - Filter button should be accessible
      const filterButton = page.getByTestId('filter-toggle')
      await expect(filterButton).toBeVisible()
      
      // When - I open filters
      await filterButton.click()
      
      // Then - Filter panel should slide in from bottom (mobile pattern)
      const filterPanel = page.getByTestId('filter-panel')
      await expect(filterPanel).toBeVisible()
      await expect(filterPanel).toHaveClass(/slide-up/)
      
      // And - I should be able to close the filter panel
      const closeFilters = page.getByTestId('close-filters')
      await closeFilters.click()
      await expect(filterPanel).not.toBeVisible()
    })

    test('When search results are displayed, Then they should be mobile-friendly', async ({ page }) => {
      // Given - I perform a search
      await page.getByTestId('search-input').fill('cafe')
      await page.keyboard.press('Enter')
      
      // Then - Results should be in a mobile-optimized layout
      const searchResults = page.getByTestId('search-results')
      await expect(searchResults).toBeVisible()
      
      // And - Each spot item should be properly sized for mobile
      const spotItems = page.getByTestId('spot-item')
      const firstSpot = spotItems.first()
      
      // Spot cards should stack vertically and fill width
      await expect(firstSpot).toHaveCSS('display', 'block')
      
      // And - Touch targets should be appropriately sized (44px minimum)
      const saveButton = firstSpot.getByTestId('save-spot-button')
      const boundingBox = await saveButton.boundingBox()
      expect(boundingBox?.height).toBeGreaterThanOrEqual(44)
    })
  })
})

/**
 * Test Data Setup
 * 
 * This test file assumes the following test data exists:
 * - Mock API responses for spot search
 * - Test spots with various solo-friendly ratings
 * - Error handling scenarios
 * 
 * The actual components tested:
 * - SearchPage component
 * - SearchInput component  
 * - SearchResults component
 * - SpotItem component
 * - FilterPanel component
 * 
 * Next Steps:
 * 1. Implement each component using TDD
 * 2. Create integration tests
 * 3. Ensure E2E tests pass
 */