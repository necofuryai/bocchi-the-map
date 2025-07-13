import { setupWorker } from 'msw/browser'
import { authHandlers } from './handlers/auth-handlers'
import { reviewHandlers } from './handlers/review-handlers'

// Setup the browser worker with all handlers
export const worker = setupWorker(
  ...authHandlers,
  ...reviewHandlers
)

// Start the worker
if (typeof window !== 'undefined') {
  worker.start({
    onUnhandledRequest: 'warn',
    serviceWorker: {
      url: '/mockServiceWorker.js',
    },
  })
}