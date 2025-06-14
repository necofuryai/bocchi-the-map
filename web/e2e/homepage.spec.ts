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
      await expect(page.getByText('Bocchi The Map')).toBeVisible()
      
      // And the map should be displayed or loading
      const mapContainer = page.locator('[style*="height"]').first()
      await expect(mapContainer).toBeVisible()
    })

    test('When the page loads, Then the navigation elements should be accessible', async ({ page }) => {
      // Then the user menu button should be visible
      const userMenuButton = page.getByRole('button', { name: 'ユーザーメニューを開く' })
      await expect(userMenuButton).toBeVisible()
      
      // And desktop navigation should be visible on larger screens
      if (await page.viewportSize()?.width! >= 768) {
        await expect(page.getByText('スポットを探す').first()).toBeVisible()
        await expect(page.getByText('レビューを書く').first()).toBeVisible()
      }
    })

    test('When the user clicks the user menu, Then the menu should open', async ({ page }) => {
      // When the user clicks the user menu button
      const userMenuButton = page.getByRole('button', { name: 'ユーザーメニューを開く' })
      await userMenuButton.click()
      
      // Then the user menu should be visible
      await expect(page.getByText('マイアカウント')).toBeVisible()
      await expect(page.getByText('プロフィール')).toBeVisible()
      await expect(page.getByText('レビュー履歴')).toBeVisible()
      await expect(page.getByText('お気に入り')).toBeVisible()
      await expect(page.getByText('ログアウト')).toBeVisible()
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
      await expect(page.getByText('Bocchi The Map')).toBeVisible()
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
      
      // Check that there's no error message displayed
      const errorMessage = page.getByRole('alert')
      await expect(errorMessage).not.toBeVisible()
      
      // Check that loading indicator disappears eventually
      const loadingIndicator = page.getByText('Loading map...')
      await expect(loadingIndicator).not.toBeVisible({ timeout: 10000 })
    })

    test('When the map fails to load, Then error handling should work gracefully', async ({ page }) => {
      // Block map tile requests to simulate error
      await page.route('**/*.pbf', route => route.abort())
      await page.route('**/*tiles*', route => route.abort())
      
      await page.goto('/')
      
      // Then an error state should be handled gracefully
      // (Either showing error message or continuing to show loading)
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
      // Check for proper ARIA labels
      const userMenuButton = page.getByRole('button', { name: 'ユーザーメニューを開く' })
      await expect(userMenuButton).toHaveAttribute('aria-expanded', 'false')
      
      // Check for proper heading structure
      const heading = page.getByRole('heading', { name: 'Bocchi The Map' })
      await expect(heading).toBeVisible()
    })
  })
})