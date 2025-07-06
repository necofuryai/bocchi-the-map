import { describe, test, beforeEach, afterEach } from 'vitest'
import { cleanup } from '@testing-library/react'

/**
 * BDD-style test helpers for frontend testing
 * These functions provide Given-When-Then syntax for better test readability
 */

// Type definitions for BDD test functions
type BDDTestFn = () => void | Promise<void>
type BDDHookFn = () => void | Promise<void>

interface BDDContext {
  given: (description: string, fn: BDDTestFn) => void
  when: (description: string, fn: BDDTestFn) => void
  then: (description: string, fn: BDDTestFn) => void
  and: (description: string, fn: BDDTestFn) => void
  beforeEach: (fn: BDDHookFn) => void
  afterEach: (fn: BDDHookFn) => void
}

/**
 * Creates a BDD-style feature test suite
 * @param featureName - Name of the feature being tested
 * @param testSuite - Function containing the test scenarios
 */
export function Feature(featureName: string, testSuite: (context: BDDContext) => void) {
  describe(`Feature: ${featureName}`, () => {
    // Auto-cleanup after each test
    afterEach(() => {
      cleanup()
    })

    const context: BDDContext = {
      given: (description: string, fn: BDDTestFn) => {
        describe(`Given ${description}`, () => {
          test('scenario setup', fn)
        })
      },
      when: (description: string, fn: BDDTestFn) => {
        describe(`When ${description}`, () => {
          test('action execution', fn)
        })
      },
      then: (description: string, fn: BDDTestFn) => {
        test(`Then ${description}`, fn)
      },
      and: (description: string, fn: BDDTestFn) => {
        test(`And ${description}`, fn)
      },
      beforeEach: (fn: BDDHookFn) => beforeEach(fn),
      afterEach: (fn: BDDHookFn) => afterEach(fn),
    }

    testSuite(context)
  })
}

/**
 * Creates a BDD-style scenario test
 * @param scenarioName - Name of the scenario
 * @param testFn - Test function
 */
export function Scenario(scenarioName: string, testFn: BDDTestFn) {
  test(`Scenario: ${scenarioName}`, testFn)
}

/**
 * Creates a nested Given-When-Then structure for complex scenarios
 * @param description - Description of the given condition
 * @param testSuite - Function containing when/then blocks
 */
export function Given(description: string, testSuite: (context: BDDContext) => void) {
  describe(`Given ${description}`, () => {
    afterEach(() => {
      cleanup()
    })

    const context: BDDContext = {
      given: (desc: string, fn: BDDTestFn) => {
        describe(`Given ${desc}`, () => {
          test('nested setup', fn)
        })
      },
      when: (desc: string, fn: BDDTestFn) => {
        describe(`When ${desc}`, () => {
          test('action', fn)
        })
      },
      then: (desc: string, fn: BDDTestFn) => {
        test(`Then ${desc}`, fn)
      },
      and: (desc: string, fn: BDDTestFn) => {
        test(`And ${desc}`, fn)
      },
      beforeEach: (fn: BDDHookFn) => beforeEach(fn),
      afterEach: (fn: BDDHookFn) => afterEach(fn),
    }

    testSuite(context)
  })
}

/**
 * Creates a When block for describing actions
 * @param description - Description of the action
 * @param testSuite - Function containing then blocks
 */
export function When(description: string, testSuite: (context: BDDContext) => void) {
  describe(`When ${description}`, () => {
    const context: BDDContext = {
      given: (desc: string, fn: BDDTestFn) => {
        describe(`Given ${desc}`, () => {
          test('setup', fn)
        })
      },
      when: (desc: string, fn: BDDTestFn) => {
        describe(`When ${desc}`, () => {
          test('nested action', fn)
        })
      },
      then: (desc: string, fn: BDDTestFn) => {
        test(`Then ${desc}`, fn)
      },
      and: (desc: string, fn: BDDTestFn) => {
        test(`And ${desc}`, fn)
      },
      beforeEach: (fn: BDDHookFn) => beforeEach(fn),
      afterEach: (fn: BDDHookFn) => afterEach(fn),
    }

    testSuite(context)
  })
}

/**
 * Utility functions for common BDD patterns
 */
export const BDDUtils = {
  /**
   * Wait for an element to appear and become stable
   */
  waitForElement: async (getElement: () => HTMLElement | null, timeout = 5000): Promise<HTMLElement> => {
    return new Promise((resolve, reject) => {
      const startTime = Date.now()
      
      const checkElement = () => {
        const element = getElement()
        if (element) {
          resolve(element)
        } else if (Date.now() - startTime > timeout) {
          reject(new Error(`Element not found within ${timeout}ms`))
        } else {
          setTimeout(checkElement, 100)
        }
      }
      
      checkElement()
    })
  },

  /**
   * Simulate user interaction delay
   */
  userDelay: (ms = 100): Promise<void> => {
    return new Promise(resolve => setTimeout(resolve, ms))
  },

  /**
   * Mock console methods to suppress logs during tests
   */
  suppressConsole: () => {
    const originalConsole = { ...console }
    
    beforeEach(() => {
      vi.spyOn(console, 'log').mockImplementation(() => {})
      vi.spyOn(console, 'warn').mockImplementation(() => {})
      vi.spyOn(console, 'error').mockImplementation(() => {})
    })
    
    afterEach(() => {
      vi.restoreAllMocks()
    })
    
    return originalConsole
  },
}