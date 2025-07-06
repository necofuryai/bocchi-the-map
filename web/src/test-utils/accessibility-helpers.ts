import { render, screen } from '@testing-library/react'
import { axe, toHaveNoViolations } from 'jest-axe'
import { ReactElement } from 'react'
import { renderWithProviders } from './render-with-providers'

// Extend Jest matchers for accessibility testing
expect.extend(toHaveNoViolations)

/**
 * Test a component for accessibility violations using axe-core
 * @param component - React component to test
 * @param options - Axe configuration options
 */
export async function testAccessibility(
  component: ReactElement,
  options?: Parameters<typeof axe>[1]
) {
  const { container } = renderWithProviders(component)
  const results = await axe(container, options)
  expect(results).toHaveNoViolations()
  return results
}

/**
 * Test keyboard navigation for a component
 * @param component - React component to test
 * @param expectedFocusableElements - Number of expected focusable elements
 */
export async function testKeyboardNavigation(
  component: ReactElement,
  expectedFocusableElements?: number
) {
  renderWithProviders(component)
  
  // Get all focusable elements
  const focusableElements = screen.getAllByRole(/button|link|textbox|combobox|checkbox|radio|tab/)
  
  if (expectedFocusableElements !== undefined) {
    expect(focusableElements).toHaveLength(expectedFocusableElements)
  }
  
  // Test Tab navigation
  let currentIndex = 0
  for (const element of focusableElements) {
    element.focus()
    expect(document.activeElement).toBe(element)
    
    // Simulate Tab key
    if (currentIndex < focusableElements.length - 1) {
      await userEvent.keyboard('{Tab}')
      currentIndex++
    }
  }
  
  return focusableElements
}

/**
 * Test screen reader accessibility
 * @param component - React component to test
 */
export function testScreenReaderAccessibility(component: ReactElement) {
  const { container } = renderWithProviders(component)
  
  // Check for proper ARIA labels
  const elementsWithAriaLabel = container.querySelectorAll('[aria-label]')
  const elementsWithAriaLabelledBy = container.querySelectorAll('[aria-labelledby]')
  const elementsWithAriaDescribedBy = container.querySelectorAll('[aria-describedby]')
  
  // Check for proper heading structure
  const headings = container.querySelectorAll('h1, h2, h3, h4, h5, h6')
  
  // Check for alt text on images
  const images = container.querySelectorAll('img')
  images.forEach(img => {
    expect(img).toHaveAttribute('alt')
  })
  
  return {
    ariaLabels: elementsWithAriaLabel.length,
    ariaLabelledBy: elementsWithAriaLabelledBy.length,
    ariaDescribedBy: elementsWithAriaDescribedBy.length,
    headings: headings.length,
    images: images.length,
  }
}

/**
 * Test color contrast for accessibility
 * @param component - React component to test
 */
export async function testColorContrast(component: ReactElement) {
  const { container } = renderWithProviders(component)
  
  // Use axe to specifically test color contrast
  const results = await axe(container, {
    rules: {
      'color-contrast': { enabled: true },
      'color-contrast-enhanced': { enabled: true },
    },
  })
  
  expect(results).toHaveNoViolations()
  return results
}

/**
 * Comprehensive accessibility test suite
 * @param component - React component to test
 * @param options - Test configuration options
 */
export async function runAccessibilityTestSuite(
  component: ReactElement,
  options: {
    skipAxe?: boolean
    skipKeyboard?: boolean
    skipScreenReader?: boolean
    skipColorContrast?: boolean
    expectedFocusableElements?: number
  } = {}
) {
  const results: any = {}
  
  if (!options.skipAxe) {
    results.axe = await testAccessibility(component)
  }
  
  if (!options.skipKeyboard) {
    results.keyboard = await testKeyboardNavigation(
      component,
      options.expectedFocusableElements
    )
  }
  
  if (!options.skipScreenReader) {
    results.screenReader = testScreenReaderAccessibility(component)
  }
  
  if (!options.skipColorContrast) {
    results.colorContrast = await testColorContrast(component)
  }
  
  return results
}

/**
 * Custom accessibility matchers for testing
 */
export const accessibilityMatchers = {
  /**
   * Check if element has proper ARIA attributes
   */
  toHaveProperAriaAttributes(received: HTMLElement) {
    const hasAriaLabel = received.hasAttribute('aria-label')
    const hasAriaLabelledBy = received.hasAttribute('aria-labelledby')
    const hasAriaDescribedBy = received.hasAttribute('aria-describedby')
    
    const pass = hasAriaLabel || hasAriaLabelledBy || hasAriaDescribedBy
    
    return {
      message: () =>
        pass
          ? `Expected element not to have proper ARIA attributes`
          : `Expected element to have at least one of: aria-label, aria-labelledby, or aria-describedby`,
      pass,
    }
  },
  
  /**
   * Check if interactive element is focusable
   */
  toBeFocusable(received: HTMLElement) {
    const interactiveTags = ['button', 'a', 'input', 'select', 'textarea']
    const isInteractive = interactiveTags.includes(received.tagName.toLowerCase())
    const hasTabIndex = received.hasAttribute('tabindex')
    const tabIndex = received.getAttribute('tabindex')
    
    const pass = isInteractive || (hasTabIndex && tabIndex !== '-1')
    
    return {
      message: () =>
        pass
          ? `Expected element not to be focusable`
          : `Expected interactive element to be focusable (have tabindex >= 0 or be naturally focusable)`,
      pass,
    }
  },
}

// Extend expect with custom matchers
declare global {
  namespace Vi {
    interface AsymmetricMatchersContaining {
      toHaveProperAriaAttributes(): any
      toBeFocusable(): any
    }
  }
}

Object.entries(accessibilityMatchers).forEach(([name, matcher]) => {
  expect.extend({ [name]: matcher })
})