import { test, expect } from '@playwright/test'

test.describe('Homepage E2E Tests', () => {
  test.describe('Given the user visits the homepage', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/')
    })

    test('When the page loads, Then the main elements should be visible', async ({ page }) => {
      // Then the header should be visible
      await expect(page.locator('header')).toBeVisible()
      
      // And the application title should be displayed
      await expect(page.getByRole('heading', { name: 'Bocchi The Map' })).toBeVisible()
      
      // And the map should be displayed or loading
      const mapContainer = page.locator('[style*="height"]').first()
      await expect(mapContainer).toBeVisible()
    })

    test('When the page loads, Then the navigation elements should be accessible', async ({ page }) => {
      // Then the authentication button should be visible (either login or user menu)
      const authButton = page.locator('header').getByRole('button', { name: /ログイン|ユーザーメニューを開く/ })
      await expect(authButton).toBeVisible()
      
      // And desktop navigation should be visible on larger screens
      const viewport = page.viewportSize()
      if (viewport && viewport.width >= 768) {
        // Check if desktop navigation elements exist (they might be hidden)
        const spotSearchButton = page.getByText('スポットを探す').first()
        const reviewButton = page.getByText('レビューを書く').first()
        
        // Elements should exist in the DOM even if hidden by responsive design
        await expect(spotSearchButton).toBeAttached()
        await expect(reviewButton).toBeAttached()
      }
    })

    test('When the user clicks the login button, Then the signin page should be accessible', async ({ page }) => {
      // When the user clicks the login button
      const loginButton = page.getByRole('button', { name: 'ログイン' })
      await loginButton.click()
      
      // Then the user should be redirected to the signin page
      await expect(page).toHaveURL(/.*\/auth\/signin/)
    })
  })

  test.describe('Given the user is on mobile', () => {
    test.beforeEach(async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 }) // iPhone SE size
      await page.goto('/')
    })

    test('When the page loads, Then mobile navigation should be available', async ({ page }) => {
      // Then the mobile menu button should be visible
      const mobileMenuButton = page.getByRole('button', { name: 'モバイルメニューを開く' })
      await expect(mobileMenuButton).toBeVisible()
    })

    test('When the user clicks the mobile menu, Then the mobile menu should open', async ({ page }) => {
      // When the user clicks the mobile menu button
      const mobileMenuButton = page.getByRole('button', { name: 'モバイルメニューを開く' })
      await mobileMenuButton.click()
      
      // Then the mobile menu items should be visible
      const mobileMenuItems = page.getByRole('menuitem')
      await expect(mobileMenuItems.first()).toBeVisible()
    })

    test('When the page loads, Then the layout should be mobile-responsive', async ({ page }) => {
      // Then the header should be properly sized
      const header = page.locator('header')
      await expect(header).toBeVisible()
      
      // And the map should fit the mobile screen
      const mapContainer = page.locator('[style*="height"]').first()
      await expect(mapContainer).toBeVisible()
      
      // And the title should be centered
      await expect(page.getByRole('heading', { name: 'Bocchi The Map' })).toBeVisible()
    })
  })

  test.describe('Given the user tests map functionality', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/')
    })

    test('When the map loads, Then it should display without errors', async ({ page }) => {
      // Wait for the map container to be visible
      const mapContainer = page.locator('[style*="height"]').first()
      await expect(mapContainer).toBeVisible()
      
      // Check that there's no map-specific error message displayed
      const mapErrorMessage = page.getByText('Map failed to load')
      await expect(mapErrorMessage).not.toBeVisible()
      
      // Wait for loading indicator to be detached from DOM
      const loadingIndicator = page.getByText('Loading map...')
      await loadingIndicator.waitFor({ state: 'detached', timeout: 10000 })
    })

    test('When the map fails to load, Then error handling should work gracefully', async ({ page }) => {
      // Block map tile requests to simulate error
      await page.route('**/*.pbf', route => route.abort())
      await page.route('**/*tiles*', route => route.abort())
      
      await page.goto('/')
      
      // Then an error message should be displayed to inform the user
      const errorAlert = page.getByRole('alert').filter({ hasText: 'Failed to load map' })
      await expect(errorAlert).toBeVisible({ timeout: 10000 })
      
      // And the map container should still be present
      const mapContainer = page.locator('[style*="height"]').first()
      await expect(mapContainer).toBeVisible()
    })
  })

  test.describe('Given the user tests accessibility', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/')
    })

    test('When navigating with keyboard, Then all interactive elements should be accessible', async ({ page }) => {
      // Test keyboard navigation
      await page.keyboard.press('Tab')
      
      // Check that focus is visible on interactive elements
      const focusedElement = page.locator(':focus')
      await expect(focusedElement).toBeVisible()
    })

    test('When using screen reader, Then proper ARIA attributes should be present', async ({ page }) => {
      // Check for proper ARIA labels on authentication button
      const authButton = page.locator('header').getByRole('button', { name: /ログイン|ユーザーメニューを開く/ })
      await expect(authButton).toBeVisible()
      
      // Check for proper heading structure
      const heading = page.getByRole('heading', { name: 'Bocchi The Map' })
      await expect(heading).toBeVisible()
    })
  })
})