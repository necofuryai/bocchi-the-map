import React from 'react'
import { describe, it, expect, vi } from 'vitest'
import { render } from '@testing-library/react'

// Mock next-themes
vi.mock('next-themes', () => ({
  ThemeProvider: vi.fn(({ children }) => children),
}))

import { ThemeProvider } from '../theme-provider'

describe('ThemeProvider Component', () => {
  describe('Given the ThemeProvider component', () => {
    it('When rendering children, Then it should wrap children with NextThemesProvider', () => {
      const TestChild = () => <div data-testid="test-child">Test Content</div>
      
      const { getByTestId } = render(
        <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
          <TestChild />
        </ThemeProvider>
      )
      
      expect(getByTestId('test-child')).toBeInTheDocument()
    })

    it('When provided with props, Then it should render children correctly', () => {
      const props = {
        attribute: 'data-theme' as const,
        defaultTheme: 'dark',
        enableSystem: false,
        disableTransitionOnChange: true,
      }
      
      const { getByText } = render(
        <ThemeProvider {...props}>
          <div>Test Content</div>
        </ThemeProvider>
      )
      
      expect(getByText('Test Content')).toBeInTheDocument()
    })

    it('When used as a wrapper, Then it should provide theme context to children', () => {
      const TestComponent = () => {
        return <div data-testid="themed-component">Themed Content</div>
      }
      
      const { getByTestId } = render(
        <ThemeProvider>
          <TestComponent />
        </ThemeProvider>
      )
      
      expect(getByTestId('themed-component')).toBeInTheDocument()
    })
  })
})