import { test, expect } from '@playwright/test'

test.describe('Authentication E2E Tests', () => {
  test.describe('Given an unauthenticated user', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/')
    })

    test('When the user is not signed in, Then sign-in option should be available', async ({ page }) => {
      // Then the login button should be visible
      const loginButton = page.getByRole('button', { name: 'ログイン' })
      await expect(loginButton).toBeVisible()
    })

    test('When the user clicks sign-in, Then authentication flow should be initiated', async ({ page }) => {
      // When the user clicks the login button
      const loginButton = page.getByRole('button', { name: 'ログイン' })
      await loginButton.click()
      
      // Then they should be redirected to the signin page
      await expect(page).toHaveURL(/.*\/auth\/signin/)
      
      // And OAuth provider options should be available
      await expect(page.getByText('Googleでログイン')).toBeVisible()
    })
  })

  test.describe('Given an authenticated user', () => {
    test.beforeEach(async ({ page }) => {
      // Mock authentication state by setting Auth.js v5 session cookie
      await page.context().addCookies([
        {
          name: 'authjs.session-token',
          value: 'mock-session-token-value',
          domain: 'localhost',
          path: '/',
          httpOnly: true,
          secure: false,
          sameSite: 'Lax'
        }
      ])
      
      // Also mock the session API response
      await page.route('**/api/auth/session', route => {
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            user: {
              id: 'mock-user-id',
              name: 'テストユーザー',
              email: 'test@example.com'
            },
            expires: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString()
          })
        })
      })
      
      await page.goto('/')
    })

    test('When the user is signed in, Then user profile should be visible', async ({ page }) => {
      const userMenuButton = page.getByRole('button', { name: 'ユーザーメニューを開く' })
      await userMenuButton.click()
      
      // Then user-specific options should be available
      await expect(page.getByText('プロフィール')).toBeVisible()
      await expect(page.getByText('レビュー履歴')).toBeVisible()
      await expect(page.getByText('お気に入り')).toBeVisible()
      await expect(page.getByText('ログアウト')).toBeVisible()
    })

    test('When the user clicks logout, Then the sign-out process should work', async ({ page }) => {
      const userMenuButton = page.getByRole('button', { name: 'ユーザーメニューを開く' })
      await userMenuButton.click()
      
      const logoutButton = page.getByText('ログアウト')
      await expect(logoutButton).toBeVisible()
      
      // Click logout button
      await logoutButton.click()
      
      // Wait for sign-out process to complete
      await page.waitForURL('/')
      
      // Verify user is signed out by checking UI state
      await expect(page.getByText('ログイン')).toBeVisible()
      await expect(page.getByRole('button', { name: 'ユーザーメニューを開く' })).not.toBeVisible()
      
      // Verify authentication cookies are cleared
      const cookies = await page.context().cookies()
      const authCookies = cookies.filter(cookie => 
        cookie.name === 'authjs.session-token' || 
        cookie.name === 'authjs.csrf-token'
      )
      expect(authCookies.length).toBe(0)
    })
  })

  test.describe('Given the authentication error handling', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto('/')
    })

    test('When authentication fails, Then error should be handled gracefully', async ({ page }) => {
      // Mock authentication error before navigation
      await page.route('**/api/auth/**', route => {
        route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Authentication failed' })
        })
      })
      
      await page.goto('/')
      
      // The application should continue to function even with auth errors
      const userMenuButton = page.getByRole('button', { name: 'ユーザーメニューを開く' })
      await expect(userMenuButton).toBeVisible()
    })

    test('When OAuth provider is unavailable, Then fallback should work', async ({ page }) => {
      // Mock OAuth provider being unavailable
      await page.route('**/api/auth/signin/**', route => {
        route.fulfill({
          status: 503,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Service unavailable' })
        })
      })
      
      await page.goto('/')
      
      // The application should still be usable
      await expect(page.getByRole('heading', { name: 'Bocchi The Map' })).toBeVisible()
    })
  })

  test.describe('Given session management', () => {
    test('When the session expires, Then the user should be handled appropriately', async ({ page }) => {
      // Start with a valid session cookie
      await page.context().addCookies([
        {
          name: 'authjs.session-token',
          value: 'valid-session-token',
          domain: 'localhost',
          path: '/',
          httpOnly: true,
          secure: false,
          sameSite: 'Lax'
        }
      ])
      
      // Mock valid session API response
      await page.route('**/api/auth/session', route => {
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            user: {
              id: 'user-id',
              name: 'テストユーザー',
              email: 'test@example.com'
            },
            expires: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString()
          })
        })
      })
      
      await page.goto('/')
      
      // Then expire the session by clearing cookie and mocking expired response
      await page.context().clearCookies()
      await page.route('**/api/auth/session', route => {
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({})
        })
      })
      await page.reload()
      
      // The application should handle the expired session gracefully
      const userMenuButton = page.getByRole('button', { name: 'ユーザーメニューを開く' })
      await expect(userMenuButton).toBeVisible()
    })

    test('When session is refreshed, Then user state should be maintained', async ({ page }) => {
      // Set Auth.js v5 session cookie
      await page.context().addCookies([
        {
          name: 'authjs.session-token',
          value: 'refreshed-session-token',
          domain: 'localhost',
          path: '/',
          httpOnly: true,
          secure: false,
          sameSite: 'Lax'
        }
      ])
      
      // Mock session API response
      await page.route('**/api/auth/session', route => {
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            user: {
              id: 'refreshed-user-id',
              name: 'リフレッシュユーザー',
              email: 'refresh@example.com'
            },
            expires: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString()
          })
        })
      })
      
      await page.goto('/')
      
      const userMenuButton = page.getByRole('button', { name: 'ユーザーメニューを開く' })
      await userMenuButton.click()
      
      // User menu should work consistently
      await expect(page.getByText('マイアカウント')).toBeVisible()
    })
  })
})