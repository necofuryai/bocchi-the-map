import { describe, it, expect, vi, afterEach } from 'vitest'
import { render } from '@testing-library/react'

// Mock next-themes
const mockThemeProvider = vi.fn(({ children }) => children)
vi.mock('next-themes', () => ({
  ThemeProvider: mockThemeProvider,
}))

import { ThemeProvider } from '../theme-provider'

describe('ThemeProvider Component', () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('Given the ThemeProvider component', () => {
    it('When rendering children, Then it should wrap children with NextThemesProvider', () => {
      const TestChild = () => <div data-testid="test-child">Test Content</div>
      
      const { getByTestId } = render(
        <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
          <TestChild />
        </ThemeProvider>
      )
      
      expect(getByTestId('test-child')).toBeInTheDocument()
      expect(mockThemeProvider).toHaveBeenCalledWith(
        expect.objectContaining({
          attribute: 'class',
          defaultTheme: 'system',
          enableSystem: true,
        })
      )
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
      expect(mockThemeProvider).toHaveBeenCalledWith(
        expect.objectContaining({
          attribute: 'data-theme',
          defaultTheme: 'dark',
          enableSystem: false,
          disableTransitionOnChange: true,
        })
      )
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
      expect(mockThemeProvider).toHaveBeenCalled()
    })
  })
})