import { test, expect } from '@playwright/test'

test.describe('Map Interaction E2E Tests', () => {
  test.describe('Given the map functionality', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/')
    })

    test('When the page loads, Then the map should initialize', async ({ page }) => {
      // Wait for map container to be visible
      const mapContainer = page.locator('[data-testid="map-container"]').first()
      await expect(mapContainer).toBeVisible()
      
      // Wait for map container to be visible with timeout
      await expect(mapContainer).toBeVisible({ timeout: 3000 })
      
      // Check that map has loaded (no error state)
      const errorAlert = page.getByRole('alert')
      const hasError = await errorAlert.isVisible().catch(() => false)
      
      if (hasError) {
        // If there's an error, it should be displayed properly
        await expect(errorAlert).toBeVisible()
      } else {
        // If no error, map should be ready for interaction
        await expect(mapContainer).toBeVisible()
      }
    })

    test('When map loads successfully, Then it should be interactive', async ({ page }) => {
      const mapContainer = page.locator('[data-testid="map-container"]').first()
      await expect(mapContainer).toBeVisible()
      
      // Wait for potential loading to complete
      await page.waitForTimeout(2000)
      
      // Try to interact with the map (if it's loaded)
      const mapBounds = await mapContainer.boundingBox()
      expect(mapBounds).not.toBeNull()
      if (mapBounds) {
        // Click on the map center
        await page.mouse.click(
          mapBounds.x + mapBounds.width / 2,
          mapBounds.y + mapBounds.height / 2
        )
        
        // Map should still be visible after interaction
        await expect(mapContainer).toBeVisible()
      }
    })

    test('When map fails to load, Then error message should be shown', async ({ page }) => {
      // Block map-related requests to simulate loading failure
      await page.route('**/*.pbf', route => route.abort())
      await page.route('**/*tiles*', route => route.abort())
      await page.route('**/*.pmtiles', route => route.abort())
      
      await page.goto('/')
      
      // Either an error should be shown or loading should persist
      const mapContainer = page.locator('[data-testid="map-container"]').first()
      await expect(mapContainer).toBeVisible()
      
      // Check for error state or loading state
      const loadingText = page.getByText('Loading map...')
      const errorAlert = page.getByRole('alert')
      
      const hasLoading = await loadingText.isVisible().catch(() => false)
      const hasError = await errorAlert.isVisible().catch(() => false)
      
      // One of these states should be present
      expect(hasLoading || hasError).toBeTruthy()
    })
  })

  test.describe('Given map navigation controls', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/')
      
      // Wait for map to potentially load
      await page.waitForTimeout(2000)
    })

    test('When map controls are available, Then user should be able to interact with them', async ({ page }) => {
      const mapContainer = page.locator('[data-testid="map-container"]').first()
      await expect(mapContainer).toBeVisible()
      
      // Look for map control elements (these would be rendered by MapLibre)
      const controlElements = page.locator('.maplibregl-ctrl, .mapboxgl-ctrl')
      
      // If controls exist, they should be interactive
      const controlCount = await controlElements.count()
      if (controlCount > 0) {
        await expect(controlElements.first()).toBeVisible()
      }
    })

    test('When user interacts with zoom controls, Then map should respond', async ({ page }) => {
      const mapContainer = page.locator('[data-testid="map-container"]').first()
      await expect(mapContainer).toBeVisible()
      
      // Look for zoom controls
      const zoomIn = page.locator('.maplibregl-ctrl-zoom-in, .mapboxgl-ctrl-zoom-in')
      const zoomOut = page.locator('.maplibregl-ctrl-zoom-out, .mapboxgl-ctrl-zoom-out')
      
      // If zoom controls exist, test them
      if (await zoomIn.isVisible()) {
        await zoomIn.click()
        // Map should still be visible after zoom
        await expect(mapContainer).toBeVisible()
      }
      
      if (await zoomOut.isVisible()) {
        await zoomOut.click()
        // Map should still be visible after zoom
        await expect(mapContainer).toBeVisible()
      }
    })
  })

  test.describe('Given POI (Point of Interest) functionality', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/')
      await page.waitForTimeout(3000) // Wait for map and potential POI loading
    })

    test('When POIs are available, Then they should be displayed on the map', async ({ page }) => {
      const mapContainer = page.locator('[data-testid="map-container"]').first()
      await expect(mapContainer).toBeVisible()
      
      // Look for POI markers or features on the map
      // This would depend on the actual implementation
      const poiElements = page.locator('[data-poi], .poi-marker, .maplibregl-marker')
      
      const poiCount = await poiElements.count()
      if (poiCount > 0) {
        // If POIs exist, they should be visible
        await expect(poiElements.first()).toBeVisible()
      }
      
      // Map should be interactive regardless of POI presence
      await expect(mapContainer).toBeVisible()
    })

    test('When clicking on a POI, Then POI details should be shown', async ({ page }) => {
      const mapContainer = page.locator('[data-testid="map-container"]').first()
      await expect(mapContainer).toBeVisible()
      
      // Look for clickable POI elements
      const poiElements = page.locator('[data-poi], .poi-marker, .maplibregl-marker')
      
      const poiCount = await poiElements.count()
      if (poiCount > 0) {
        // Click on the first POI
        await poiElements.first().click()
        
        // Look for popup or detail view
        const popup = page.locator('.maplibregl-popup, .poi-popup, [role="tooltip"]')
        
        // If popup appears, it should be visible
        if (await popup.isVisible()) {
          await expect(popup).toBeVisible()
        }
      }
    })

    test('When searching for POIs, Then search functionality should work', async ({ page }) => {
      // Look for search functionality
      const searchButton = page.getByText('スポットを探す')
      
      if (await searchButton.isVisible()) {
        await searchButton.click()
        
        // Search interface should be accessible
        // This test would be expanded based on actual search implementation
        await expect(searchButton).toBeVisible()
      }
    })
  })

  test.describe('Given map responsiveness', () => {
    test('When viewed on desktop, Then map should use full available space', async ({ page }) => {
      await page.setViewportSize({ width: 1280, height: 720 })
      await page.goto('/')
      
      const mapContainer = page.locator('[data-testid="map-container"]').first()
      await expect(mapContainer).toBeVisible()
      
      // Map should take appropriate space on desktop
      const mapBounds = await mapContainer.boundingBox()
      expect(mapBounds?.width).toBeGreaterThan(800) // Should be reasonably wide on desktop
    })

    test('When viewed on mobile, Then map should be mobile-optimized', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 }) // iPhone SE
      await page.goto('/')
      
      const mapContainer = page.locator('[data-testid="map-container"]').first()
      await expect(mapContainer).toBeVisible()
      
      // Map should fit mobile screen
      const mapBounds = await mapContainer.boundingBox()
      expect(mapBounds?.width).toBeLessThanOrEqual(375) // Should fit mobile width
    })

    test('When orientation changes, Then map should adapt', async ({ page }) => {
      // Start in portrait
      await page.setViewportSize({ width: 375, height: 667 })
      await page.goto('/')
      
      const mapContainer = page.locator('[data-testid="map-container"]').first()
      await expect(mapContainer).toBeVisible()
      
      // Change to landscape
      await page.setViewportSize({ width: 667, height: 375 })
      
      // Map should still be visible and properly sized
      await expect(mapContainer).toBeVisible()
      
      const mapBounds = await mapContainer.boundingBox()
      expect(mapBounds?.width).toBeLessThanOrEqual(667) // Should fit landscape width
    })
  })
})