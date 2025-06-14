import { describe, it, expect, vi } from 'vitest'
import { render } from '@testing-library/react'

// Mock next-auth/react
const mockSessionProvider = vi.fn(({ children }) => children)
vi.mock('next-auth/react', () => ({
  SessionProvider: mockSessionProvider,
}))

import { AuthProvider } from '../auth-provider'

describe('AuthProvider Component', () => {
  describe('Given the AuthProvider component', () => {
    it('When rendering children, Then it should wrap children with SessionProvider', () => {
      const TestChild = () => <div data-testid="test-child">Test Content</div>
      
      const { getByTestId } = render(
        <AuthProvider>
          <TestChild />
        </AuthProvider>
      )
      
      expect(getByTestId('test-child')).toBeInTheDocument()
      expect(mockSessionProvider).toHaveBeenCalled()
    })

    it('When used as a wrapper, Then it should provide authentication context to children', () => {
      const TestComponent = () => {
        return <div data-testid="auth-component">Authenticated Content</div>
      }
      
      const { getByTestId } = render(
        <AuthProvider>
          <TestComponent />
        </AuthProvider>
      )
      
      expect(getByTestId('auth-component')).toBeInTheDocument()
      expect(mockSessionProvider).toHaveBeenCalled()
    })

    it('When multiple children are provided, Then all children should be rendered', () => {
      const { getByTestId } = render(
        <AuthProvider>
          <div data-testid="child-1">Child 1</div>
          <div data-testid="child-2">Child 2</div>
        </AuthProvider>
      )
      
      expect(getByTestId('child-1')).toBeInTheDocument()
      expect(getByTestId('child-2')).toBeInTheDocument()
      expect(mockSessionProvider).toHaveBeenCalled()
    })
  })
})